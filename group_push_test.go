package engagelab

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewGroupPushClient_Defaults(t *testing.T) {
	gc := NewGroupPushClient("groupKey", "groupSecret")
	wantAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("group-groupKey:groupSecret"))
	if gc.authHeader != wantAuth {
		t.Errorf("authHeader = %q, want %q", gc.authHeader, wantAuth)
	}
	if gc.baseURL != string(Singapore) {
		t.Errorf("baseURL = %q, want %q", gc.baseURL, Singapore)
	}
}

func TestNewGroupPushClient_WithOptions(t *testing.T) {
	gc := NewGroupPushClient("gk", "gs", WithDataCenter(Frankfurt))
	if gc.baseURL != string(Frankfurt) {
		t.Errorf("baseURL = %q, want %q", gc.baseURL, Frankfurt)
	}
}

func TestGroupPushClient_Send(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v4/grouppush" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("missing Authorization header")
		}

		body, _ := io.ReadAll(r.Body)
		var param GroupPushParam
		json.Unmarshal(body, &param)
		if param.From != "group-app" {
			t.Errorf("From = %q, want %q", param.From, "group-app")
		}

		json.NewEncoder(w).Encode(GroupPushResult{
			GroupMsgID: "gmsg_001",
		})
	}))
	defer ts.Close()

	gc := NewGroupPushClient("gk", "gs", WithBaseURL(ts.URL))
	result, err := gc.Send(context.Background(), &GroupPushParam{
		From: "group-app",
		To:   "all",
		Body: &PushBody{
			Platform:     "all",
			Notification: &NotificationMessage{Alert: "group push!"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.GroupMsgID != "gmsg_001" {
		t.Errorf("GroupMsgID = %q", result.GroupMsgID)
	}
}

func TestGroupPushClient_Send_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 1003, "message": "auth failed"},
		})
	}))
	defer ts.Close()

	gc := NewGroupPushClient("bad", "cred", WithBaseURL(ts.URL))
	_, err := gc.Send(context.Background(), &GroupPushParam{})
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*ApiError)
	if !ok {
		t.Fatalf("expected *ApiError, got %T", err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("StatusCode = %d, want 401", apiErr.StatusCode)
	}
}
