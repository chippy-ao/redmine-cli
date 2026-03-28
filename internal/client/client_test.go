package client

import (
	"encoding/json"
	"io"
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

func TestPost_SendsCorrectMethodAndBody(t *testing.T) {
	var receivedMethod string
	var receivedBody []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"issue": map[string]any{"id": 123}})
	}))
	defer srv.Close()

	c := New(srv.URL, "test-api-key")
	body := map[string]any{"issue": map[string]any{"project_id": "test", "subject": "テスト"}}
	var result map[string]any
	err := c.Post("/issues.json", body, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedMethod != http.MethodPost {
		t.Errorf("expected POST method, got %s", receivedMethod)
	}

	var sentBody map[string]any
	if err := json.Unmarshal(receivedBody, &sentBody); err != nil {
		t.Fatalf("failed to parse sent body: %v", err)
	}
	issue, ok := sentBody["issue"].(map[string]any)
	if !ok {
		t.Fatal("expected issue key in body")
	}
	if issue["subject"] != "テスト" {
		t.Errorf("expected subject テスト, got %v", issue["subject"])
	}
}

func TestPost_APIKeyHeader(t *testing.T) {
	var receivedAPIKey string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAPIKey = r.Header.Get("X-Redmine-API-Key")
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{})
	}))
	defer srv.Close()

	c := New(srv.URL, "my-secret-key")
	var result any
	err := c.Post("/issues.json", map[string]any{}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedAPIKey != "my-secret-key" {
		t.Errorf("expected API key my-secret-key, got %s", receivedAPIKey)
	}
}

func TestPost_422ReturnsValidationErrors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{"errors": []string{"Subject cannot be blank", "Project is not valid"}})
	}))
	defer srv.Close()

	c := New(srv.URL, "key")
	var result any
	err := c.Post("/issues.json", map[string]any{}, &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "Subject cannot be blank") {
		t.Errorf("expected validation error message, got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "Project is not valid") {
		t.Errorf("expected second validation error, got %q", err.Error())
	}
}

func TestPost_422FallbackOnInvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	c := New(srv.URL, "key")
	var result any
	err := c.Post("/issues.json", map[string]any{}, &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "not json") {
		t.Errorf("expected raw body in error, got %q", err.Error())
	}
}

func TestDelete_SendsDeleteRequest(t *testing.T) {
	var receivedMethod string
	var receivedPath string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := New(srv.URL, "key")
	err := c.Delete("/relations/1.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedMethod != http.MethodDelete {
		t.Errorf("expected DELETE method, got %s", receivedMethod)
	}
	if receivedPath != "/relations/1.json" {
		t.Errorf("expected path /relations/1.json, got %s", receivedPath)
	}
}

func TestDelete_404ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer srv.Close()

	c := New(srv.URL, "key")
	err := c.Delete("/relations/999.json")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("expected 404 error, got %q", err.Error())
	}
}
