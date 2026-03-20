package engagelab

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPushService_Send(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v4/push" {
			t.Errorf("path = %s, want /v4/push", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var param PushParam
		json.Unmarshal(body, &param)
		if param.From != "test-app" {
			t.Errorf("From = %q, want %q", param.From, "test-app")
		}

		json.NewEncoder(w).Encode(PushResult{MsgID: "msg_001", RequestID: "req_001"})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Push.Send(context.Background(), &PushParam{
		From: "test-app",
		To:   "all",
		Body: &PushBody{
			Platform: "all",
			Notification: &NotificationMessage{
				Alert: "Hello!",
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.MsgID != "msg_001" {
		t.Errorf("MsgID = %q, want %q", result.MsgID, "msg_001")
	}
	if result.RequestID != "req_001" {
		t.Errorf("RequestID = %q, want %q", result.RequestID, "req_001")
	}
}

func TestPushService_Send_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 2002, "message": "invalid param"},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	_, err := c.Push.Send(context.Background(), &PushParam{})
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*ApiError)
	if !ok {
		t.Fatalf("expected *ApiError, got %T", err)
	}
	if apiErr.ErrorBody.Code != 2002 {
		t.Errorf("error code = %d, want 2002", apiErr.ErrorBody.Code)
	}
}

func TestPushService_SendRaw(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/push" {
			t.Errorf("path = %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var raw map[string]interface{}
		json.Unmarshal(body, &raw)
		if raw["custom_key"] != "custom_value" {
			t.Errorf("raw body missing custom_key")
		}
		json.NewEncoder(w).Encode(PushResult{MsgID: "raw_001"})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Push.SendRaw(context.Background(), map[string]string{"custom_key": "custom_value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.MsgID != "raw_001" {
		t.Errorf("MsgID = %q", result.MsgID)
	}
}

func TestPushService_Validate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/push/validate" {
			t.Errorf("path = %s, want /v4/push/validate", r.URL.Path)
		}
		json.NewEncoder(w).Encode(PushResult{MsgID: "validate_001"})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Push.Validate(context.Background(), &PushParam{
		To: "all",
		Body: &PushBody{
			Platform:     "all",
			Notification: &NotificationMessage{Alert: "test"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.MsgID != "validate_001" {
		t.Errorf("MsgID = %q", result.MsgID)
	}
}

func TestPushService_Withdraw(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/v4/push/withdraw/msg123" {
			t.Errorf("path = %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(PushWithdrawResult{MsgID: "msg123"})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Push.Withdraw(context.Background(), "msg123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.MsgID != "msg123" {
		t.Errorf("MsgID = %q", result.MsgID)
	}
}

func TestPushService_BatchByRegID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/batch/push/regid" {
			t.Errorf("path = %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(BatchPushResult{
			Results: map[string]BatchPushSingleResult{
				"regid1": {Target: "regid1", Success: true, MsgID: 100},
			},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Push.BatchByRegID(context.Background(), &BatchPushParam{
		Requests: []BatchPushRequest{
			{
				Target:       "regid1",
				Platform:     "android",
				Notification: &NotificationMessage{Alert: "batch test"},
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r, ok := result.Results["regid1"]; !ok || !r.Success {
		t.Error("batch result missing or not successful for regid1")
	}
}

func TestPushService_BatchByAlias(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/batch/push/alias" {
			t.Errorf("path = %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(BatchPushResult{
			Results: map[string]BatchPushSingleResult{
				"alias1": {Target: "alias1", Success: true, MsgID: 200},
			},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Push.BatchByAlias(context.Background(), &BatchPushParam{
		Requests: []BatchPushRequest{
			{
				Target:       "alias1",
				Platform:     "ios",
				Notification: &NotificationMessage{Alert: "alias batch"},
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r, ok := result.Results["alias1"]; !ok || r.MsgID != 200 {
		t.Error("batch result incorrect for alias1")
	}
}

func TestPushParam_JSONSerialization(t *testing.T) {
	builderID := 1
	priority := 2
	param := &PushParam{
		From: "sender",
		To:   "all",
		Body: &PushBody{
			Platform: "all",
			Notification: &NotificationMessage{
				Alert: "Hello World",
				Android: &AndroidNotification{
					Alert:     "Android Alert",
					Title:     "Title",
					BuilderID: &builderID,
					Priority:  &priority,
					Extras:    map[string]interface{}{"key": "val"},
				},
				IOS: &IOSNotification{
					Alert: "iOS Alert",
					Sound: "default",
					Badge: 1,
				},
			},
		},
	}

	data, err := json.Marshal(param)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded map[string]interface{}
	json.Unmarshal(data, &decoded)

	if decoded["from"] != "sender" {
		t.Errorf("from = %v", decoded["from"])
	}
	if decoded["to"] != "all" {
		t.Errorf("to = %v", decoded["to"])
	}

	body := decoded["body"].(map[string]interface{})
	notification := body["notification"].(map[string]interface{})
	if notification["alert"] != "Hello World" {
		t.Errorf("alert = %v", notification["alert"])
	}
}
