package engagelab

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type PlanService struct {
	client *Client
}

// CreateOrUpdate creates a new push plan or updates an existing one.
// POST /v4/push_plan
func (s *PlanService) CreateOrUpdate(ctx context.Context, param *PushPlanParam) (*PushPlanResult, error) {
	var result PushPlanResult
	err := s.client.doPost(ctx, "/v4/push_plan", param, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// List returns a paginated list of push plans.
// GET /v4/push_plan/list?page_index=&page_size=&send_source=&search_description=
func (s *PlanService) List(ctx context.Context, pageIndex, pageSize int, sendSource *int, searchDescription string) (*PushPlanListResult, error) {
	query := url.Values{}
	query.Set("page_index", strconv.Itoa(pageIndex))
	query.Set("page_size", strconv.Itoa(pageSize))
	if sendSource != nil {
		query.Set("send_source", strconv.Itoa(*sendSource))
	}
	if searchDescription != "" {
		query.Set("search_description", searchDescription)
	}
	var result PushPlanListResult
	err := s.client.doGet(ctx, "/v4/push_plan/list", query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// QueryMsg queries message IDs by push plan IDs within a date range.
// GET /v4/status/plan/msg/?plan_ids=&start_date=&end_date=
func (s *PlanService) QueryMsg(ctx context.Context, planIDs string, startDate, endDate string) (*PushPlanMsgQueryResult, error) {
	query := url.Values{}
	query.Set("plan_ids", planIDs)
	if startDate != "" {
		query.Set("start_date", startDate)
	}
	if endDate != "" {
		query.Set("end_date", endDate)
	}
	var result PushPlanMsgQueryResult
	err := s.client.doGet(ctx, "/v4/status/plan/msg/", query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes a push plan by ID.
// DELETE /v4/push_plan/{planId}
func (s *PlanService) Delete(ctx context.Context, planID string) (*PushPlanDeleteResult, error) {
	var result PushPlanDeleteResult
	err := s.client.doDelete(ctx, fmt.Sprintf("/v4/push_plan/%s", planID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BatchDelete removes multiple push plans by IDs (comma-separated).
// DELETE /v4/push_plan/batch/{planIds}
func (s *PlanService) BatchDelete(ctx context.Context, planIDs string) error {
	return s.client.doDelete(ctx, fmt.Sprintf("/v4/push_plan/batch/%s", planIDs), nil)
}
