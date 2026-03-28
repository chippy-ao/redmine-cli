package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client は Redmine API へのHTTPクライアントです。
type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

// New は新しい Client を生成します。baseURL の末尾スラッシュは除去されます。
func New(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

// Get は指定パスに対してGETリクエストを送信し、結果を result にデコードします。
func (c *Client) Get(path string, params map[string]string, result any) error {
	u := c.baseURL + path
	if len(params) > 0 {
		v := url.Values{}
		for key, val := range params {
			v.Set(key, val)
		}
		u += "?" + v.Encode()
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("リクエスト作成エラー: %w", err)
	}
	return c.doRequest(req, result)
}

// GetRawQuery は生のクエリ文字列をそのまま使用してGETリクエストを送信します。
func (c *Client) GetRawQuery(path string, rawQuery string, result any) error {
	u := c.baseURL + path
	if rawQuery != "" {
		u += "?" + rawQuery
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("リクエスト作成エラー: %w", err)
	}
	return c.doRequest(req, result)
}

// Post は指定パスに対してPOSTリクエストを送信し、結果を result にデコードします。
func (c *Client) Post(path string, body any, result any) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("リクエストボディのJSON変換エラー: %w", err)
	}

	u := c.baseURL + path
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, u, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("リクエスト作成エラー: %w", err)
	}
	return c.doRequest(req, result)
}

// Delete は指定パスに対してDELETEリクエストを送信します。
func (c *Client) Delete(path string) error {
	u := c.baseURL + path
	req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, u, nil)
	if err != nil {
		return fmt.Errorf("リクエスト作成エラー: %w", err)
	}
	return c.doRequest(req, nil)
}

func (c *Client) doRequest(req *http.Request, result any) error {
	req.Header.Set("X-Redmine-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("リクエスト送信エラー: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		// 200/201: JSON デコード
	case http.StatusNoContent:
		// 204: デコード不要
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("認証エラー: APIキーが無効です。")
	case http.StatusForbidden:
		return fmt.Errorf("権限エラー: アクセス権がありません。")
	case http.StatusNotFound:
		return fmt.Errorf("リソースが見つかりません（404）")
	case http.StatusUnprocessableEntity:
		return c.parseValidationError(resp.Body)
	default:
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Redmine APIエラー（%d）: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("レスポンスデコードエラー: %w", err)
		}
	}

	return nil
}

func (c *Client) parseValidationError(body io.Reader) error {
	var redmineErr struct {
		Errors []string `json:"errors"`
	}
	data, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("バリデーションエラー（レスポンス読み取り失敗）")
	}
	if err := json.Unmarshal(data, &redmineErr); err != nil || len(redmineErr.Errors) == 0 {
		return fmt.Errorf("バリデーションエラー: %s", string(data))
	}
	return fmt.Errorf("バリデーションエラー: %s", strings.Join(redmineErr.Errors, ", "))
}
