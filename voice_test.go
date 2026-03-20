package engagelab

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVoiceService_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v4/voices" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var param VoiceParam
		json.Unmarshal(body, &param)
		if param.Language != "zh-CN" {
			t.Errorf("Language = %q", param.Language)
		}
		json.NewEncoder(w).Encode(VoiceResult{
			Language: "zh-CN",
			Content:  "你的验证码是{code}",
			TTSType:  "standard",
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Voice.Create(context.Background(), &VoiceParam{
		Language: "zh-CN",
		Content:  "你的验证码是{code}",
		TTSType:  "standard",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Language != "zh-CN" {
		t.Errorf("Language = %q", result.Language)
	}
}

func TestVoiceService_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v4/voices" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		json.NewEncoder(w).Encode(VoiceListResult{
			Voices: []VoiceResult{
				{Language: "zh-CN", Content: "test", TTSType: "standard"},
				{Language: "en-US", Content: "test en", TTSType: "standard"},
			},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Voice.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Voices) != 2 {
		t.Errorf("Voices len = %d, want 2", len(result.Voices))
	}
}

func TestVoiceService_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/voices/zh-CN" {
			t.Errorf("path = %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(VoiceResult{Language: "zh-CN", Content: "hello"})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Voice.Get(context.Background(), "zh-CN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Language != "zh-CN" {
		t.Errorf("Language = %q", result.Language)
	}
}

func TestVoiceService_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/v4/voices/en-US" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	err := c.Voice.Delete(context.Background(), "en-US")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
