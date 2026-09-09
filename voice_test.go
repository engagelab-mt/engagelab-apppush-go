package engagelab

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestVoiceService_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v4/voices" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatal(err)
		}
		if r.FormValue("language") != "zh-CN" {
			t.Errorf("language = %q", r.FormValue("language"))
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		data, _ := io.ReadAll(file)
		if string(data) != "voice" {
			t.Errorf("file = %q", data)
		}
		json.NewEncoder(w).Encode(VoiceResult{FileURL: "https://example.com/voice.mp3"})
	}))
	defer ts.Close()

	voiceFile, err := os.CreateTemp("", "voice-*.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(voiceFile.Name())
	voiceFile.WriteString("voice")
	voiceFile.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Voice.Create(context.Background(), "zh-CN", voiceFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if result.FileURL == "" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestVoiceService_ListGetDelete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v4/voices":
			json.NewEncoder(w).Encode([]VoiceResult{{Language: "en", FileURL: "u"}})
		case r.Method == http.MethodGet && r.URL.Path == "/v4/voices/en":
			json.NewEncoder(w).Encode(VoiceResult{Language: "en", FileURL: "u"})
		case r.Method == http.MethodDelete && r.URL.Path == "/v4/voices/en":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
	}))
	defer ts.Close()
	c := NewClient("k", "s", WithBaseURL(ts.URL))
	list, err := c.Voice.List(context.Background())
	if err != nil || len(list) != 1 {
		t.Fatalf("list=%#v err=%v", list, err)
	}
	result, err := c.Voice.Get(context.Background(), "en")
	if err != nil || result.FileURL != "u" {
		t.Fatalf("get=%#v err=%v", result, err)
	}
	if err := c.Voice.Delete(context.Background(), "en"); err != nil {
		t.Fatal(err)
	}
}
