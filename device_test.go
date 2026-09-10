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
		if err := json.Unmarshal(body, &param); err != nil {
			t.Fatal(err)
		}
		if len(param.RegistrationIDs) != 2 {
			t.Errorf("registration_ids len = %d, want 2", len(param.RegistrationIDs))
		}

		online := true
		offline := false
		if err := json.NewEncoder(w).Encode([]DeviceStatusGetResult{
			{RegistrationID: "reg1", Online: &online, LastOnlineTime: "2025-01-01 12:00:00"},
			{RegistrationID: "reg2", Online: &offline, LastOnlineTime: "2025-01-01 11:00:00"},
		}); err != nil {
			t.Fatal(err)
		}
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
		if err := json.NewEncoder(w).Encode(DeviceGetResult{
			Tags:  []string{"tag1", "tag2"},
			Alias: "user_alias",
		}); err != nil {
			t.Fatal(err)
		}
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
		var param struct {
			Tags  DeviceSetTags `json:"tags"`
			Alias string        `json:"alias"`
		}
		if err := json.Unmarshal(body, &param); err != nil {
			t.Fatal(err)
		}
		if param.Alias != "new_alias" {
			t.Errorf("alias = %q, want %q", param.Alias, "new_alias")
		}
		if len(param.Tags.Add) != 1 || param.Tags.Add[0] != "vip" {
			t.Errorf("unexpected tags: %#v", param.Tags)
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

func TestDeviceService_Set_ClearTags(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var param map[string]interface{}
		if err := json.Unmarshal(body, &param); err != nil {
			t.Fatal(err)
		}
		if tags, ok := param["tags"].(string); !ok || tags != "" {
			t.Errorf("unexpected tags: %#v", param["tags"])
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	if err := c.Device.Set(context.Background(), "reg123", &DeviceSetParam{ClearTags: true}); err != nil {
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
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 1004, "message": "forbidden"},
		}); err != nil {
			t.Fatal(err)
		}
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

func TestDeviceService_RegisterToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v4/devices/token/registration_id" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var param DeviceTokenRegisterParam
		if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
			t.Fatal(err)
		}
		if param.Platform != "android" || len(param.Tokens) != 1 {
			t.Errorf("unexpected param: %#v", param)
		}
		if err := json.NewEncoder(w).Encode(DeviceTokenRegisterResult{Results: []DeviceTokenResult{{
			Token: "t1", RegistrationID: "r1", IsNew: true, Code: 0,
		}, {
			Token: "", IsNew: false, Code: 21003, Message: "invalid fcm token format",
		}}}); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()
	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Device.RegisterToken(context.Background(), &DeviceTokenRegisterParam{
		Platform: "android", Tokens: []string{"t1"},
	})
	if err != nil || len(result.Results) != 2 || result.Results[0].RegistrationID != "r1" || result.Results[1].Code != 21003 {
		t.Fatalf("result=%#v err=%v", result, err)
	}
}
