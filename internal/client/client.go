package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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
		http:    &http.Client{},
	}
}

// Get は指定パスに対してGETリクエストを送信し、結果を result にデコードします。
func (c *Client) Get(path string, params map[string]string, result interface{}) error {
	u := c.baseURL + path
	if len(params) > 0 {
		v := url.Values{}
		for key, val := range params {
			v.Set(key, val)
		}
		u += "?" + v.Encode()
	}
	return c.doGet(u, result)
}

// GetRawQuery は生のクエリ文字列をそのまま使用してGETリクエストを送信します。
func (c *Client) GetRawQuery(path string, rawQuery string, result interface{}) error {
	u := c.baseURL + path
	if rawQuery != "" {
		u += "?" + rawQuery
	}
	return c.doGet(u, result)
}

func (c *Client) doGet(rawURL string, result interface{}) error {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return fmt.Errorf("リクエスト作成エラー: %w", err)
	}

	req.Header.Set("X-Redmine-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("リクエスト送信エラー: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return fmt.Errorf("認証エラー: APIキーが無効です。")
		case http.StatusForbidden:
			return fmt.Errorf("権限エラー: アクセス権がありません。")
		case http.StatusNotFound:
			return fmt.Errorf("リソースが見つかりません（404）")
		default:
			return fmt.Errorf("Redmine APIエラー（%d）: %s", resp.StatusCode, string(body))
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("レスポンスデコードエラー: %w", err)
	}

	return nil
}
