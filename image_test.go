package engagelab

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImageService_UploadOppo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v4/image/oppo" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json; charset=utf-8" {
			t.Errorf("unexpected Content-Type: %s", r.Header.Get("Content-Type"))
		}
		var param OppoImageParam
		if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
			t.Fatal(err)
		}
		if param.BigPictureURL != "https://example.com/big.png" || param.SmallPictureURL != "" {
			t.Errorf("unexpected param: %#v", param)
		}
		if err := json.NewEncoder(w).Encode(ImageUploadResult{BigPictureID: "big_001"}); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Image.UploadOppo(context.Background(), &OppoImageParam{BigPictureURL: "https://example.com/big.png"})
	if err != nil {
		t.Fatal(err)
	}
	if result.BigPictureID != "big_001" {
		t.Errorf("unexpected result: %#v", result)
	}
}
