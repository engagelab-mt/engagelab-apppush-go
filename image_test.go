package engagelab

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestImageService_UploadOppo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v4/image/oppo" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("Content-Type = %q, want multipart/form-data", ct)
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("FormFile error: %v", err)
		}
		defer file.Close()
		if header.Filename != "test.png" {
			t.Errorf("filename = %q, want test.png", header.Filename)
		}

		json.NewEncoder(w).Encode(ImageUploadResult{MediaID: "media_001"})
	}))
	defer ts.Close()

	tmpFile, err := os.CreateTemp("", "test*.png")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Write([]byte("fake png content"))
	tmpFile.Close()

	// Rename to test.png for predictable filename
	tmpPath := tmpFile.Name()
	testPath := strings.TrimSuffix(tmpPath, tmpPath[strings.LastIndex(tmpPath, "/"):]) + "/test.png"
	os.Rename(tmpPath, testPath)
	defer os.Remove(testPath)

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Image.UploadOppo(context.Background(), testPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.MediaID != "media_001" {
		t.Errorf("MediaID = %q", result.MediaID)
	}
}

func TestImageService_UploadOppo_FileNotFound(t *testing.T) {
	c := NewClient("k", "s")
	_, err := c.Image.UploadOppo(context.Background(), "/nonexistent/file.png")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestImageService_UploadOppoFromReader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v4/image/oppo" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("FormFile error: %v", err)
		}
		defer file.Close()
		if header.Filename != "banner.jpg" {
			t.Errorf("filename = %q, want banner.jpg", header.Filename)
		}

		json.NewEncoder(w).Encode(ImageUploadResult{MediaID: "media_002"})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	reader := bytes.NewReader([]byte("fake jpg content"))
	result, err := c.Image.UploadOppoFromReader(context.Background(), "banner.jpg", reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.MediaID != "media_002" {
		t.Errorf("MediaID = %q", result.MediaID)
	}
}

func TestImageService_UploadOppoFromReader_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 4001, "message": "invalid image"},
		})
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	_, err := c.Image.UploadOppoFromReader(context.Background(), "bad.jpg", bytes.NewReader([]byte{}))
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*ApiError)
	if !ok {
		t.Fatalf("expected *ApiError, got %T", err)
	}
	if apiErr.ErrorBody.Code != 4001 {
		t.Errorf("error code = %d, want 4001", apiErr.ErrorBody.Code)
	}
}
