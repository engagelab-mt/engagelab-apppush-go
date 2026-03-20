package engagelab

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestServer(handler http.HandlerFunc) (*httptest.Server, *Client) {
	ts := httptest.NewServer(handler)
	c := NewClient("test-key", "test-secret", WithBaseURL(ts.URL))
	return ts, c
}

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("myKey", "mySecret")

	if c.baseURL != string(Singapore) {
		t.Errorf("default baseURL = %q, want %q", c.baseURL, Singapore)
	}

	wantAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("myKey:mySecret"))
	if c.authHeader != wantAuth {
		t.Errorf("authHeader = %q, want %q", c.authHeader, wantAuth)
	}

	if c.Push == nil || c.Device == nil || c.Tag == nil || c.Alias == nil ||
		c.Schedule == nil || c.Status == nil || c.Plan == nil || c.Voice == nil || c.Image == nil {
		t.Error("one or more services are nil")
	}
}

func TestWithDataCenter(t *testing.T) {
	tests := []struct {
		dc   DataCenter
		want string
	}{
		{Singapore, "https://pushapi-sgp.engagelab.com"},
		{HongKong, "https://pushapi-hk.engagelab.com"},
		{Virginia, "https://pushapi-usva.engagelab.com"},
		{Frankfurt, "https://pushapi-defra.engagelab.com"},
	}
	for _, tt := range tests {
		c := NewClient("k", "s", WithDataCenter(tt.dc))
		if c.baseURL != tt.want {
			t.Errorf("WithDataCenter(%v): baseURL = %q, want %q", tt.dc, c.baseURL, tt.want)
		}
	}
}

func TestWithBaseURL(t *testing.T) {
	c := NewClient("k", "s", WithBaseURL("https://custom.example.com"))
	if c.baseURL != "https://custom.example.com" {
		t.Errorf("WithBaseURL: got %q", c.baseURL)
	}
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{Timeout: 99 * time.Second}
	c := NewClient("k", "s", WithHTTPClient(custom))
	if c.httpClient != custom {
		t.Error("WithHTTPClient did not set custom http client")
	}
}

func TestWithTimeout(t *testing.T) {
	c := NewClient("k", "s", WithTimeout(5*time.Second))
	if c.httpClient.Timeout != 5*time.Second {
		t.Errorf("WithTimeout: got %v, want 5s", c.httpClient.Timeout)
	}
}

func TestDoRequest_SetsHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got == "" {
			t.Error("missing Authorization header")
		}
		if got := r.Header.Get("Content-Type"); got != "application/json; charset=utf-8" {
			t.Errorf("Content-Type = %q", got)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Errorf("Accept = %q", got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	_ = c.doPost(context.Background(), "/test", map[string]string{"a": "b"}, nil)
}

func TestDoRequest_PostBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var data map[string]string
		json.Unmarshal(body, &data)
		if data["key"] != "value" {
			t.Errorf("request body key = %q, want %q", data["key"], "value")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	_ = c.doPost(context.Background(), "/test", map[string]string{"key": "value"}, nil)
}

func TestDoRequest_ParsesResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"msg_id": "12345"})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	var result PushResult
	err := c.doPost(context.Background(), "/test", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.MsgID != "12345" {
		t.Errorf("MsgID = %q, want %q", result.MsgID, "12345")
	}
}

func TestDoRequest_ApiError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    1003,
				"message": "auth failed",
			},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	err := c.doPost(context.Background(), "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*ApiError)
	if !ok {
		t.Fatalf("expected *ApiError, got %T", err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("StatusCode = %d, want 401", apiErr.StatusCode)
	}
	if apiErr.ErrorBody.Code != 1003 {
		t.Errorf("ErrorBody.Code = %d, want 1003", apiErr.ErrorBody.Code)
	}
	if apiErr.ErrorBody.Message != "auth failed" {
		t.Errorf("ErrorBody.Message = %q", apiErr.ErrorBody.Message)
	}
}

func TestApiError_ErrorString(t *testing.T) {
	e := &ApiError{
		StatusCode: 401,
		ErrorBody:  ErrorDetail{Code: 1003, Message: "auth failed"},
	}
	want := "engagelab api error: status=401 code=1003 message=auth failed"
	if e.Error() != want {
		t.Errorf("Error() = %q, want %q", e.Error(), want)
	}
}

func TestDoGet_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Query().Get("page") != "1" {
			t.Errorf("query page = %q, want %q", r.URL.Query().Get("page"), "1")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	_ = c.doGet(context.Background(), "/test", map[string][]string{"page": {"1"}}, nil)
}

func TestDoDelete_Method(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	_ = c.doDelete(context.Background(), "/test", nil)
}

func TestDoPut_Method(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	_ = c.doPut(context.Background(), "/test", nil, nil)
}

func TestDoRequest_ContextCancelled(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := c.doPost(ctx, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestDoRequest_EmptyResponseBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	err := c.doPost(context.Background(), "/test", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
