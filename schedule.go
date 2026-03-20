package engagelab

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// SchedulePushParam is the request for creating a scheduled push.
type SchedulePushParam struct {
	Name    string           `json:"name,omitempty"`
	Enabled *bool            `json:"enabled,omitempty"`
	Trigger *ScheduleTrigger `json:"trigger,omitempty"`
	Push    *PushParam       `json:"push,omitempty"`
}

type ScheduleTrigger struct {
	Single     *TriggerSingle     `json:"single,omitempty"`
	Periodical *TriggerPeriodical `json:"periodical,omitempty"`
}

type TriggerSingle struct {
	Time     string `json:"time,omitempty"`     // "yyyy-MM-dd HH:mm:ss"
	ZoneType *int   `json:"zone_type,omitempty"`
}

type TriggerPeriodical struct {
	Start     string   `json:"start,omitempty"`     // "yyyy-MM-dd HH:mm:ss"
	End       string   `json:"end,omitempty"`       // "yyyy-MM-dd HH:mm:ss"
	Time      string   `json:"time,omitempty"`      // "HH:mm:ss"
	Frequency *int     `json:"frequency,omitempty"`
	TimeUnit  string   `json:"time_unit,omitempty"` // "day", "WEEK", "MONTH"
	Point     []string `json:"point,omitempty"`
	ZoneType  *int     `json:"zone_type,omitempty"`
}

// SchedulePushResult is the response after creating a schedule.
type SchedulePushResult struct {
	ScheduleID string `json:"schedule_id,omitempty"`
	Name       string `json:"name,omitempty"`
}

// SchedulePushGetResult is the response for getting a schedule by ID.
type SchedulePushGetResult struct {
	ScheduleID string           `json:"schedule_id,omitempty"`
	Name       string           `json:"name,omitempty"`
	Enabled    *bool            `json:"enabled,omitempty"`
	Trigger    *ScheduleTrigger `json:"trigger,omitempty"`
	Push       *PushParam       `json:"push,omitempty"`
}

// SchedulePushListResult is the paginated schedule list response.
type SchedulePushListResult struct {
	TotalCount  int                      `json:"total_count"`
	TotalPages  int                      `json:"total_pages"`
	CurrentPage int                      `json:"page"`
	Schedules   []SchedulePushListDetail `json:"schedules,omitempty"`
}

type SchedulePushListDetail struct {
	ScheduleID string      `json:"schedule_id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Enabled    *bool       `json:"enabled,omitempty"`
	Trigger    interface{} `json:"trigger,omitempty"`
	Push       interface{} `json:"push,omitempty"`
}

// SchedulePushDetailGetResult contains msg IDs for a schedule.
type SchedulePushDetailGetResult struct {
	Count  int64         `json:"count"`
	MsgIDs []interface{} `json:"MsgIds,omitempty"`
}

// --- ScheduleService ---

type ScheduleService struct {
	client *Client
}

// Create creates a new scheduled push task.
// POST /v4/schedules
func (s *ScheduleService) Create(ctx context.Context, param *SchedulePushParam) (*SchedulePushResult, error) {
	var result SchedulePushResult
	err := s.client.doPost(ctx, "/v4/schedules", param, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Update modifies an existing scheduled push task.
// PUT /v4/schedules/{schedule_id}
func (s *ScheduleService) Update(ctx context.Context, scheduleID string, param *SchedulePushParam) (*SchedulePushGetResult, error) {
	var result SchedulePushGetResult
	err := s.client.doPut(ctx, fmt.Sprintf("/v4/schedules/%s", scheduleID), param, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes a scheduled push task.
// DELETE /v4/schedules/{schedule_id}
func (s *ScheduleService) Delete(ctx context.Context, scheduleID string) error {
	return s.client.doDelete(ctx, fmt.Sprintf("/v4/schedules/%s", scheduleID), nil)
}

// Get retrieves a scheduled push task by ID.
// GET /v4/schedules/{schedule_id}
func (s *ScheduleService) Get(ctx context.Context, scheduleID string) ([]SchedulePushGetResult, error) {
	var result []SchedulePushGetResult
	err := s.client.doGet(ctx, fmt.Sprintf("/v4/schedules/%s", scheduleID), nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// List returns a paginated list of scheduled push tasks.
// GET /v4/schedules?page={page}
func (s *ScheduleService) List(ctx context.Context, page int) (*SchedulePushListResult, error) {
	query := url.Values{}
	query.Set("page", strconv.Itoa(page))
	var result SchedulePushListResult
	err := s.client.doGet(ctx, "/v4/schedules", query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMsgIDs retrieves message IDs for a scheduled push.
// GET /v4/schedules/{schedule_id}/msg-ids
func (s *ScheduleService) GetMsgIDs(ctx context.Context, scheduleID string) (*SchedulePushDetailGetResult, error) {
	var result SchedulePushDetailGetResult
	err := s.client.doGet(ctx, fmt.Sprintf("/v4/schedules/%s/msg-ids", scheduleID), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
