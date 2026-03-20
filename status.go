package engagelab

import (
	"context"
	"net/url"
	"strconv"
	"strings"
)

// UserStatusGetResult is the response for user statistics.
type UserStatusGetResult struct {
	TimeUnit string           `json:"TimeUnit,omitempty"`
	Start    string           `json:"start,omitempty"`
	Duration int              `json:"duration,omitempty"`
	Items    []UserStatusItem `json:"items,omitempty"`
}

type UserStatusItem struct {
	Time    string              `json:"time,omitempty"`
	Android *UserStatusPlatform `json:"android,omitempty"`
	IOS     *UserStatusPlatform `json:"ios,omitempty"`
}

type UserStatusPlatform struct {
	New    int64 `json:"new"`
	Active int64 `json:"active"`
	Online int64 `json:"online"`
}

// MessageStatusGetResult is the response for message delivery statistics.
type MessageStatusGetResult struct {
	Targets     int64             `json:"targets"`
	Sent        int64             `json:"sent"`
	Delivered   int64             `json:"delivered"`
	Impressions int64             `json:"impressions"`
	Clicks      int64             `json:"clicks"`
	Sub         *MessageStatusSub `json:"sub,omitempty"`
}

type MessageStatusSub struct {
	Notification *MessageStatusDetail `json:"notification,omitempty"`
	Message      *MessageStatusDetail `json:"message,omitempty"`
}

type MessageStatusDetail struct {
	Target      int64                 `json:"target"`
	Sent        int64                 `json:"sent"`
	Delivered   int64                 `json:"delivered"`
	Impressions int64                 `json:"impressions"`
	Click       int64                 `json:"click"`
	SubAndroid  *MessageStatusAndroid `json:"sub_android,omitempty"`
	SubIOS      *MessageStatusIOS     `json:"sub_ios,omitempty"`
}

type MessageStatusAndroid struct {
	EngageLabAndroid *MessageStatusChannel `json:"engageLab_android,omitempty"`
	Xiaomi           *MessageStatusChannel `json:"xiaomi,omitempty"`
	Huawei           *MessageStatusChannel `json:"huawei,omitempty"`
	Honor            *MessageStatusChannel `json:"honor,omitempty"`
	Meizu            *MessageStatusChannel `json:"meizu,omitempty"`
	Oppo             *MessageStatusChannel `json:"oppo,omitempty"`
	Vivo             *MessageStatusChannel `json:"vivo,omitempty"`
	FCM              *MessageStatusChannel `json:"fcm,omitempty"`
}

type MessageStatusIOS struct {
	EngageLabIOS *MessageStatusChannel `json:"engageLab_ios,omitempty"`
	APNS         *MessageStatusChannel `json:"apns,omitempty"`
}

type MessageStatusChannel struct {
	Targets     int64 `json:"targets"`
	Sent        int64 `json:"sent"`
	Delivered   int64 `json:"delivered"`
	Impressions int64 `json:"impressions"`
	Clicks      int64 `json:"clicks"`
}

// MessageLifecycleGetResult is the lifecycle status for a message on a device.
type MessageLifecycleGetResult struct {
	Status       string `json:"status,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// --- StatusService ---

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
