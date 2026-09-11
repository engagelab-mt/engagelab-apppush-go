package engagelab

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAliasService_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v4/aliases/user001" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("platform") != "android,ios" {
			t.Errorf("platform = %q", r.URL.Query().Get("platform"))
		}
		if err := json.NewEncoder(w).Encode(AliasStatusGetResult{
			RegistrationIDs: []string{"reg1", "reg2", "reg3"},
		}); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Alias.Get(context.Background(), "user001", []string{"android", "ios"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.RegistrationIDs) != 3 {
		t.Errorf("registration_ids count = %d, want 3", len(result.RegistrationIDs))
	}
}

func TestAliasService_Get_NoPlatform(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("expected no query, got %q", r.URL.RawQuery)
		}
		if err := json.NewEncoder(w).Encode(AliasStatusGetResult{
			RegistrationIDs: []string{"reg1"},
		}); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Alias.Get(context.Background(), "user001", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.RegistrationIDs) != 1 {
		t.Errorf("registration_ids count = %d, want 1", len(result.RegistrationIDs))
	}
}

func TestAliasService_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/v4/aliases/user001" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("platform") != "android" {
			t.Errorf("platform = %q", r.URL.Query().Get("platform"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	err := c.Alias.Delete(context.Background(), "user001", []string{"android"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAliasService_Delete_NoPlatform(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("expected no query, got %q", r.URL.RawQuery)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	err := c.Alias.Delete(context.Background(), "user001", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
