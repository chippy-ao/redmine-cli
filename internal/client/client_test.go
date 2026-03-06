package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestGet_SendsCorrectURLAndAPIKey(t *testing.T) {
	var receivedPath string
	var receivedQuery string
	var receivedAPIKey string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedQuery = r.URL.RawQuery
		receivedAPIKey = r.Header.Get("X-Redmine-API-Key")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(testResponse{ID: 1, Name: "test"})
	}))
	defer srv.Close()

	c := New(srv.URL, "test-api-key")
	var result testResponse
	err := c.Get("/issues.json", map[string]string{"status_id": "open", "limit": "10"}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedPath != "/issues.json" {
		t.Errorf("expected path /issues.json, got %s", receivedPath)
	}
	if receivedAPIKey != "test-api-key" {
		t.Errorf("expected API key test-api-key, got %s", receivedAPIKey)
	}
	if !strings.Contains(receivedQuery, "status_id=open") {
		t.Errorf("expected query to contain status_id=open, got %s", receivedQuery)
	}
	if !strings.Contains(receivedQuery, "limit=10") {
		t.Errorf("expected query to contain limit=10, got %s", receivedQuery)
	}
}

func TestGetRawQuery_PreservesRawQueryString(t *testing.T) {
	var receivedRawQuery string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedRawQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(testResponse{ID: 2, Name: "raw"})
	}))
	defer srv.Close()

	c := New(srv.URL, "test-api-key")
	rawQuery := "f[]=status_id&op[status_id]==&v[status_id][]=1"
	var result testResponse
	err := c.GetRawQuery("/issues.json", rawQuery, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedRawQuery != rawQuery {
		t.Errorf("expected raw query %q, got %q", rawQuery, receivedRawQuery)
	}
}

func TestGet_401ReturnsAuthError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	}))
	defer srv.Close()

	c := New(srv.URL, "bad-key")
	var result testResponse
	err := c.Get("/issues.json", nil, &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expected := "認証エラー: APIキーが無効です。"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestGet_403ReturnsForbiddenError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
	}))
	defer srv.Close()

	c := New(srv.URL, "key")
	var result testResponse
	err := c.Get("/issues.json", nil, &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expected := "権限エラー: アクセス権がありません。"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestGet_404ReturnsNotFoundError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer srv.Close()

	c := New(srv.URL, "key")
	var result testResponse
	err := c.Get("/issues/999.json", nil, &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expected := "リソースが見つかりません（404）"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestGet_SuccessfulResponseDecoded(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(testResponse{ID: 42, Name: "decoded"})
	}))
	defer srv.Close()

	c := New(srv.URL, "key")
	var result testResponse
	err := c.Get("/test.json", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 42 {
		t.Errorf("expected ID 42, got %d", result.ID)
	}
	if result.Name != "decoded" {
		t.Errorf("expected Name decoded, got %s", result.Name)
	}
}

func TestNew_TrimsTrailingSlash(t *testing.T) {
	c := New("https://redmine.example.com/", "key")
	if c.baseURL != "https://redmine.example.com" {
		t.Errorf("expected trailing slash trimmed, got %s", c.baseURL)
	}

	c2 := New("https://redmine.example.com///", "key")
	if c2.baseURL != "https://redmine.example.com" {
		t.Errorf("expected all trailing slashes trimmed, got %s", c2.baseURL)
	}
}

func TestGet_OtherErrorReturnsStatusAndBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer srv.Close()

	c := New(srv.URL, "key")
	var result testResponse
	err := c.Get("/test.json", nil, &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to contain status code 500, got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "internal server error") {
		t.Errorf("expected error to contain body, got %q", err.Error())
	}
}
