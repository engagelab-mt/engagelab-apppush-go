package engagelab

import (
	"context"
	"net/url"
	"strconv"
	"strings"
)

type StatusService struct {
	client *Client
}

// Users returns user statistics (new, active, online) for a time range.
// GET /v4/status/users?time_unit={}&start={}&duration={}
func (s *StatusService) Users(ctx context.Context, timeUnit string, start string, duration int) (*UserStatusGetResult, error) {
	query := url.Values{}
	query.Set("time_unit", timeUnit)
	query.Set("start", start)
	query.Set("duration", strconv.Itoa(duration))
	var result UserStatusGetResult
	err := s.client.doGet(ctx, "/v4/status/users", query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// MessageDetail returns delivery statistics for given message IDs.
// GET /v4/status/detail?message_ids={ids}
func (s *StatusService) MessageDetail(ctx context.Context, messageIDs []string) (map[string]MessageStatusGetResult, error) {
	query := url.Values{}
	query.Set("message_ids", strings.Join(messageIDs, ","))
	result := make(map[string]MessageStatusGetResult)
	err := s.client.doGet(ctx, "/v4/status/detail", query, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// MessageLifecycle returns the lifecycle status for a message on specific devices.
// GET /v4/status/message?message_id={}&registration_ids={}
func (s *StatusService) MessageLifecycle(ctx context.Context, messageID string, registrationIDs []string) (map[string]MessageLifecycleGetResult, error) {
	query := url.Values{}
	query.Set("message_id", messageID)
	query.Set("registration_ids", strings.Join(registrationIDs, ","))
	result := make(map[string]MessageLifecycleGetResult)
	err := s.client.doGet(ctx, "/v4/status/message", query, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// BatchMessageDetail returns delivery statistics for multiple messages (batch query).
// GET /v4/status/batch/message?message_ids={}
func (s *StatusService) BatchMessageDetail(ctx context.Context, messageIDs []string) (map[string]MessageStatusGetResult, error) {
	query := url.Values{}
	query.Set("message_ids", strings.Join(messageIDs, ","))
	result := make(map[string]MessageStatusGetResult)
	err := s.client.doGet(ctx, "/v4/status/batch/message", query, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// PlanDetail returns message statistics for a push plan.
// GET /v4/status/plan/detail?plan_id={}&message_ids={}
func (s *StatusService) PlanDetail(ctx context.Context, planID string, messageIDs []string) (map[string]MessageStatusGetResult, error) {
	query := url.Values{}
	query.Set("plan_id", planID)
	if len(messageIDs) > 0 {
		query.Set("message_ids", strings.Join(messageIDs, ","))
	}
	result := make(map[string]MessageStatusGetResult)
	err := s.client.doGet(ctx, "/v4/status/plan/detail", query, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
