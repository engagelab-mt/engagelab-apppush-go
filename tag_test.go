package engagelab

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTagService_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v4/tags" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		encodeJSON(t, w, TagsGetResult{Tags: []string{"vip", "test", "beta"}})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Tag.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Tags) != 3 {
		t.Errorf("tags count = %d, want 3", len(result.Tags))
	}
}

func TestTagService_Set(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v4/tags/vip" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var param TagSetParam
		if err := json.Unmarshal(body, &param); err != nil {
			t.Fatal(err)
		}
		if len(param.RegistrationIDs.Add) != 2 {
			t.Errorf("add count = %d, want 2", len(param.RegistrationIDs.Add))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	err := c.Tag.Set(context.Background(), "vip", &TagSetParam{
		RegistrationIDs: &TagRegistrationIDs{
			Add: []string{"reg1", "reg2"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTagService_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/v4/tags/old_tag" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("platform") != "android,ios" {
			t.Errorf("platform = %q", r.URL.Query().Get("platform"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	err := c.Tag.Delete(context.Background(), "old_tag", []string{"android", "ios"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTagService_Delete_NoPlatform(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("expected empty query, got %q", r.URL.RawQuery)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	err := c.Tag.Delete(context.Background(), "tag1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTagService_GetCount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/tags_count" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if len(r.URL.Query()["tags"]) != 2 || r.URL.Query().Get("platform") != "android" {
			t.Errorf("query = %q", r.URL.RawQuery)
		}
		encodeJSON(t, w, TagsCountGetResult{
			TagsCount: map[string]int64{"vip": 100, "beta": 50},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Tag.GetCount(context.Background(), []string{"vip", "beta"}, "android")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TagsCount["vip"] != 100 {
		t.Errorf("vip count = %d, want 100", result.TagsCount["vip"])
	}
}

func TestTagService_GetDeviceStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/tags/vip/registration_ids/reg123" {
			t.Errorf("path = %s", r.URL.Path)
		}
		encodeJSON(t, w, TagStatusGetResult{Result: true})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Tag.GetDeviceStatus(context.Background(), "vip", "reg123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Result {
		t.Errorf("result = %#v", result)
	}
}

func TestTagService_GetQuota(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/tags/quota-info" {
			t.Errorf("path = %s", r.URL.Path)
		}
		encodeJSON(t, w, TagQuotaGetResult{
			Data: &TagQuotaData{
				TotalTagQuota:   1000,
				UseTagQuota:     100,
				TotalAliasQuota: 500,
				UseAliasQuota:   50,
			},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Tag.GetQuota(context.Background(), []string{"vip"}, "android")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Data.TotalTagQuota != 1000 {
		t.Errorf("TotalTagQuota = %d, want 1000", result.Data.TotalTagQuota)
	}
	if result.Data.UseTagQuota != 100 {
		t.Errorf("UseTagQuota = %d, want 100", result.Data.UseTagQuota)
	}
}
