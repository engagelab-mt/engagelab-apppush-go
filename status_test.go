package engagelab

import (
	"context"
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
		encodeJSON(t, w, UserStatusGetResult{
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
		encodeJSON(t, w, map[string]interface{}{
			"msg1": map[string]interface{}{
				"targets": 1000, "sent": 990, "delivered": 900,
				"sub": map[string]interface{}{
					"notification":  map[string]interface{}{"target": 11, "click": 12},
					"message":       map[string]interface{}{"targets": 21, "clicks": 22},
					"live_activity": map[string]interface{}{"targets": 31, "clicks": 32},
					"voip": map[string]interface{}{
						"targets": 41, "clicks": 42, "delivered": 2,
						"sub_hmos": map[string]interface{}{
							"harmonyos": map[string]interface{}{"delivered": 12},
						},
					},
					"inapp_message": map[string]interface{}{"targets": 51, "clicks": 52},
				},
			},
			"msg2": map[string]interface{}{"targets": 500, "sent": 495, "delivered": 480},
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
	if result["msg1"].Sub.Notification.Target != 11 || result["msg1"].Sub.Notification.Click != 12 ||
		result["msg1"].Sub.Message.Target != 21 || result["msg1"].Sub.Message.Click != 22 ||
		result["msg1"].Sub.LiveActivity.Target != 31 || result["msg1"].Sub.LiveActivity.Click != 32 ||
		result["msg1"].Sub.VoIP.Target != 41 || result["msg1"].Sub.VoIP.Click != 42 ||
		result["msg1"].Sub.InAppMessage.Target != 51 || result["msg1"].Sub.InAppMessage.Click != 52 ||
		result["msg1"].Sub.VoIP.Delivered != 2 || result["msg1"].Sub.VoIP.SubHMOS.HarmonyOS.Delivered != 12 {
		t.Errorf("nested status fields were not decoded: %#v", result["msg1"].Sub)
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
		encodeJSON(t, w, map[string]MessageLifecycleGetResult{
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
		encodeJSON(t, w, []MessageLifecycleGetResult{{MessageID: "msg1", Status: "sent"}})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Status.BatchMessageDetail(context.Background(), []string{"msg1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].MessageID != "msg1" {
		t.Errorf("unexpected result: %#v", result)
	}
}

func TestStatusService_PlanDetail(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/status/plan/detail" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("plan_ids") != "plan1" {
			t.Errorf("plan_ids = %q", r.URL.Query().Get("plan_ids"))
		}
		encodeJSON(t, w, map[string]MessageStatusGetResult{
			"msg1": {Targets: 200, Delivered: 180},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Status.PlanDetail(context.Background(), []string{"plan1"}, "2026-01-01", "2026-01-02")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["msg1"].Delivered != 180 {
		t.Errorf("msg1 delivered = %d, want 180", result["msg1"].Delivered)
	}
}

func TestStatusService_PlanDetail_Query(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("start_date") != "2026-01-01" || r.URL.Query().Get("end_date") != "2026-01-31" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		encodeJSON(t, w, map[string]MessageStatusGetResult{})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	_, err := c.Status.PlanDetail(context.Background(), []string{"plan1"}, "2026-01-01", "2026-01-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
