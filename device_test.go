package engagelab

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeviceService_GetStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v4/devices/status" {
			t.Errorf("path = %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var param DeviceStatusGetParam
		json.Unmarshal(body, &param)
		if len(param.RegistrationIDs) != 2 {
			t.Errorf("registration_ids len = %d, want 2", len(param.RegistrationIDs))
		}

		online := true
		offline := false
		json.NewEncoder(w).Encode([]DeviceStatusGetResult{
			{RegistrationID: "reg1", Online: &online, LastOnlineTime: "2025-01-01 12:00:00"},
			{RegistrationID: "reg2", Online: &offline, LastOnlineTime: "2025-01-01 11:00:00"},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	results, err := c.Device.GetStatus(context.Background(), &DeviceStatusGetParam{
		RegistrationIDs: []string{"reg1", "reg2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("results len = %d, want 2", len(results))
	}
	if results[0].RegistrationID != "reg1" || !*results[0].Online {
		t.Error("first device should be reg1 online")
	}
	if results[1].RegistrationID != "reg2" || *results[1].Online {
		t.Error("second device should be reg2 offline")
	}
}

func TestDeviceService_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v4/devices/reg123" {
			t.Errorf("path = %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(DeviceGetResult{
			Tags:  []string{"tag1", "tag2"},
			Alias: "user_alias",
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Device.Get(context.Background(), "reg123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Tags) != 2 {
		t.Errorf("tags len = %d, want 2", len(result.Tags))
	}
	if result.Alias != "user_alias" {
		t.Errorf("alias = %q, want %q", result.Alias, "user_alias")
	}
}

func TestDeviceService_Set(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v4/devices/reg123" {
			t.Errorf("path = %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var param DeviceSetParam
		json.Unmarshal(body, &param)
		if param.Alias != "new_alias" {
			t.Errorf("alias = %q, want %q", param.Alias, "new_alias")
		}
		if len(param.Tags.Add) != 1 || param.Tags.Add[0] != "vip" {
			t.Error("tags add should contain 'vip'")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	err := c.Device.Set(context.Background(), "reg123", &DeviceSetParam{
		Alias: "new_alias",
		Tags: &DeviceSetTags{
			Add: []string{"vip"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeviceService_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/v4/devices/reg123" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	err := c.Device.Delete(context.Background(), "reg123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeviceService_GetStatus_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 1004, "message": "forbidden"},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	_, err := c.Device.GetStatus(context.Background(), &DeviceStatusGetParam{
		RegistrationIDs: []string{"reg1"},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*ApiError)
	if !ok {
		t.Fatalf("expected *ApiError, got %T", err)
	}
	if apiErr.StatusCode != 403 {
		t.Errorf("status = %d, want 403", apiErr.StatusCode)
	}
}
