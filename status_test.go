package engagelab

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusService_Users(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v4/status/users" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("time_unit") != "DAY" {
			t.Errorf("time_unit = %q", r.URL.Query().Get("time_unit"))
		}
		if r.URL.Query().Get("start") != "2025-01-01" {
			t.Errorf("start = %q", r.URL.Query().Get("start"))
		}
		if r.URL.Query().Get("duration") != "7" {
			t.Errorf("duration = %q", r.URL.Query().Get("duration"))
		}
		json.NewEncoder(w).Encode(UserStatusGetResult{
			TimeUnit: "DAY",
			Start:    "2025-01-01",
			Duration: 7,
			Items: []UserStatusItem{
				{
					Time:    "2025-01-01",
					Android: &UserStatusPlatform{New: 100, Active: 500, Online: 200},
					IOS:     &UserStatusPlatform{New: 50, Active: 300, Online: 100},
				},
			},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Status.Users(context.Background(), "DAY", "2025-01-01", 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Duration != 7 {
		t.Errorf("Duration = %d, want 7", result.Duration)
	}
	if len(result.Items) != 1 {
		t.Fatalf("Items len = %d, want 1", len(result.Items))
	}
	if result.Items[0].Android.New != 100 {
		t.Errorf("Android.New = %d, want 100", result.Items[0].Android.New)
	}
}

func TestStatusService_MessageDetail(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/status/detail" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("message_ids") != "msg1,msg2" {
			t.Errorf("message_ids = %q", r.URL.Query().Get("message_ids"))
		}
		json.NewEncoder(w).Encode(map[string]MessageStatusGetResult{
			"msg1": {Targets: 1000, Sent: 990, Delivered: 900},
			"msg2": {Targets: 500, Sent: 495, Delivered: 480},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Status.MessageDetail(context.Background(), []string{"msg1", "msg2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("result len = %d, want 2", len(result))
	}
	if result["msg1"].Targets != 1000 {
		t.Errorf("msg1 targets = %d, want 1000", result["msg1"].Targets)
	}
}

func TestStatusService_MessageLifecycle(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/status/message" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("message_id") != "msg1" {
			t.Errorf("message_id = %q", r.URL.Query().Get("message_id"))
		}
		json.NewEncoder(w).Encode(map[string]MessageLifecycleGetResult{
			"reg1": {Status: "delivered"},
			"reg2": {Status: "failed", ErrorMessage: "device offline"},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Status.MessageLifecycle(context.Background(), "msg1", []string{"reg1", "reg2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["reg1"].Status != "delivered" {
		t.Errorf("reg1 status = %q", result["reg1"].Status)
	}
	if result["reg2"].ErrorMessage != "device offline" {
		t.Errorf("reg2 error = %q", result["reg2"].ErrorMessage)
	}
}

func TestStatusService_BatchMessageDetail(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/status/batch/message" {
			t.Errorf("path = %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]MessageStatusGetResult{
			"msg1": {Targets: 100, Sent: 99},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Status.BatchMessageDetail(context.Background(), []string{"msg1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["msg1"].Sent != 99 {
		t.Errorf("msg1 sent = %d, want 99", result["msg1"].Sent)
	}
}

func TestStatusService_PlanDetail(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/status/plan/detail" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("plan_id") != "plan1" {
			t.Errorf("plan_id = %q", r.URL.Query().Get("plan_id"))
		}
		json.NewEncoder(w).Encode(map[string]MessageStatusGetResult{
			"msg1": {Targets: 200, Delivered: 180},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Status.PlanDetail(context.Background(), "plan1", []string{"msg1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["msg1"].Delivered != 180 {
		t.Errorf("msg1 delivered = %d, want 180", result["msg1"].Delivered)
	}
}

func TestStatusService_PlanDetail_NoMsgIDs(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("message_ids") != "" {
			t.Errorf("message_ids should be empty, got %q", r.URL.Query().Get("message_ids"))
		}
		json.NewEncoder(w).Encode(map[string]MessageStatusGetResult{})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	_, err := c.Status.PlanDetail(context.Background(), "plan1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
