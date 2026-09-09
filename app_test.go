package engagelab

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppService_GetVIPStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v4/app/vip/status" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		json.NewEncoder(w).Encode(AppVIPStatusResult{VIPStatus: 1, VIPEndTime: 1775059200})
	}))
	defer ts.Close()
	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.App.GetVIPStatus(context.Background())
	if err != nil || result.VIPStatus != 1 || result.VIPEndTime != 1775059200 {
		t.Fatalf("result=%#v err=%v", result, err)
	}
}
