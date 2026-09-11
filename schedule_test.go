package engagelab

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScheduleService_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v4/schedules" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var param SchedulePushParam
		if err := json.Unmarshal(body, &param); err != nil {
			t.Fatal(err)
		}
		if param.Name != "daily_push" {
			t.Errorf("name = %q", param.Name)
		}
		encodeJSON(t, w, SchedulePushResult{
			ScheduleID: "sched_001",
			Name:       "daily_push",
		})
	}))
	defer ts.Close()

	enabled := true
	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Schedule.Create(context.Background(), &SchedulePushParam{
		Name:    "daily_push",
		Enabled: &enabled,
		Trigger: &ScheduleTrigger{
			Single: &TriggerSingle{Time: "2025-12-01 10:00:00"},
		},
		Push: &PushParam{
			To:   "all",
			Body: &PushBody{Platform: "all", Notification: &NotificationMessage{Alert: "scheduled!"}},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ScheduleID != "sched_001" {
		t.Errorf("ScheduleID = %q", result.ScheduleID)
	}
}

func TestScheduleService_Update(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/v4/schedules/sched_001" {
			t.Errorf("path = %s", r.URL.Path)
		}
		enabled := true
		encodeJSON(t, w, SchedulePushGetResult{
			ScheduleID: "sched_001",
			Name:       "updated_push",
			Enabled:    &enabled,
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Schedule.Update(context.Background(), "sched_001", &SchedulePushParam{
		Name: "updated_push",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "updated_push" {
		t.Errorf("Name = %q", result.Name)
	}
}

func TestScheduleService_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/v4/schedules/sched_001" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	err := c.Schedule.Delete(context.Background(), "sched_001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScheduleService_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v4/schedules/sched_001" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		enabled := true
		encodeJSON(t, w, []SchedulePushGetResult{
			{ScheduleID: "sched_001", Name: "test", Enabled: &enabled},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	results, err := c.Schedule.Get(context.Background(), "sched_001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("results len = %d, want 1", len(results))
	}
	if results[0].ScheduleID != "sched_001" {
		t.Errorf("ScheduleID = %q", results[0].ScheduleID)
	}
}

func TestScheduleService_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/schedules" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "1" {
			t.Errorf("page = %q", r.URL.Query().Get("page"))
		}
		encodeJSON(t, w, SchedulePushListResult{
			TotalCount:  10,
			TotalPages:  2,
			CurrentPage: 1,
			Schedules: []SchedulePushListDetail{
				{ScheduleID: "s1", Name: "first"},
			},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Schedule.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalCount != 10 {
		t.Errorf("TotalCount = %d, want 10", result.TotalCount)
	}
	if len(result.Schedules) != 1 {
		t.Errorf("Schedules len = %d, want 1", len(result.Schedules))
	}
}

func TestScheduleService_GetMsgIDs(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/schedules/sched_001/msg-ids" {
			t.Errorf("path = %s", r.URL.Path)
		}
		encodeJSON(t, w, SchedulePushDetailGetResult{
			Count:  3,
			MsgIDs: []interface{}{"m1", "m2", "m3"},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Schedule.GetMsgIDs(context.Background(), "sched_001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Count != 3 {
		t.Errorf("Count = %d, want 3", result.Count)
	}
	if len(result.MsgIDs) != 3 {
		t.Errorf("MsgIDs len = %d, want 3", len(result.MsgIDs))
	}
}
