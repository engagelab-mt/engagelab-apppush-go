package engagelab

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlanService_CreateOrUpdate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v4/push_plan" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var param PushPlanParam
		if err := json.Unmarshal(body, &param); err != nil {
			t.Fatal(err)
		}
		if param.PlanDescription != "test plan" {
			t.Errorf("PlanDescription = %q", param.PlanDescription)
		}
		if err := json.NewEncoder(w).Encode(PushPlanResult{PlanID: "plan_001"}); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Plan.CreateOrUpdate(context.Background(), &PushPlanParam{
		PlanDescription: "test plan",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.PlanID != "plan_001" {
		t.Errorf("PlanID = %q", result.PlanID)
	}
}

func TestPlanService_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/push_plan/list" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("page_index") != "0" {
			t.Errorf("page_index = %q", r.URL.Query().Get("page_index"))
		}
		if r.URL.Query().Get("page_size") != "10" {
			t.Errorf("page_size = %q", r.URL.Query().Get("page_size"))
		}
		if err := json.NewEncoder(w).Encode(PushPlanListResult{
			Total: 5,
			PushPlanInfo: []PushPlanInfo{
				{PlanID: "p1", PlanDescription: "plan one", Count: 10},
				{PlanID: "p2", PlanDescription: "plan two", Count: 20},
			},
		}); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Plan.List(context.Background(), 0, 10, nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 5 {
		t.Errorf("Total = %d, want 5", result.Total)
	}
	if len(result.PushPlanInfo) != 2 {
		t.Errorf("PushPlanInfo len = %d, want 2", len(result.PushPlanInfo))
	}
}

func TestPlanService_List_WithFilters(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("send_source") != "1" {
			t.Errorf("send_source = %q", r.URL.Query().Get("send_source"))
		}
		if r.URL.Query().Get("search_description") != "promo" {
			t.Errorf("search_description = %q", r.URL.Query().Get("search_description"))
		}
		if err := json.NewEncoder(w).Encode(PushPlanListResult{Total: 1}); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	src := 1
	_, err := c.Plan.List(context.Background(), 0, 10, &src, "promo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPlanService_QueryMsg(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/status/plan/msg/" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("plan_ids") != "p1,p2" {
			t.Errorf("plan_ids = %q", r.URL.Query().Get("plan_ids"))
		}
		result := PushPlanMsgQueryResult{
			"p1": PlanMsgInfo{MsgIDs: []string{"m1", "m2"}},
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Plan.QueryMsg(context.Background(), "p1,p2", "2025-01-01", "2025-01-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info, ok := (*result)["p1"]; !ok || len(info.MsgIDs) != 2 {
		t.Errorf("unexpected result for p1: %+v", result)
	}
}

func TestPlanService_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/v4/push_plan/plan_001" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode(PushPlanDeleteResult{PlanID: "plan_001"}); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	result, err := c.Plan.Delete(context.Background(), "plan_001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.PlanID != "plan_001" {
		t.Errorf("PlanID = %q", result.PlanID)
	}
}

func TestPlanService_BatchDelete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/v4/push_plan/batch/p1,p2,p3" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("k", "s", WithBaseURL(ts.URL))
	err := c.Plan.BatchDelete(context.Background(), "p1,p2,p3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
