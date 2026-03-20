package engagelab

import (
	"context"
	"fmt"
)

// PushParam represents a push request payload.
type PushParam struct {
	From       string                 `json:"from,omitempty"`
	To         interface{}            `json:"to,omitempty"` // "all" or *PushTo
	RequestID  string                 `json:"request_id,omitempty"`
	CustomArgs map[string]interface{} `json:"custom_args,omitempty"`
	Body       *PushBody              `json:"body,omitempty"`
}

type PushBody struct {
	Platform     interface{}          `json:"platform,omitempty"` // "all" or []string{"android","ios","hmos"}
	Notification *NotificationMessage `json:"notification,omitempty"`
	Message      *CustomMessage       `json:"message,omitempty"`
	LiveActivity *LiveActivityMessage `json:"live_activity,omitempty"`
	Options      *Options             `json:"options,omitempty"`
}

// PushTo specifies target audience for push.
type PushTo struct {
	RegistrationID []string `json:"registration_id,omitempty"`
	Tag            []string `json:"tag,omitempty"`
	TagAnd         []string `json:"tag_and,omitempty"`
	TagNot         []string `json:"tag_not,omitempty"`
	Alias          []string `json:"alias,omitempty"`
	LiveActivityID string   `json:"live_activity_id,omitempty"`
	Seg            *Seg     `json:"seg,omitempty"`
}

type Seg struct {
	ID string `json:"id,omitempty"`
}

// NotificationMessage is the notification payload.
type NotificationMessage struct {
	Alert   string               `json:"alert,omitempty"`
	Android *AndroidNotification `json:"android,omitempty"`
	IOS     *IOSNotification     `json:"ios,omitempty"`
	Hmos    *HmosNotification    `json:"hmos,omitempty"`
}

type AndroidNotification struct {
	Alert             string                 `json:"alert,omitempty"`
	Title             string                 `json:"title,omitempty"`
	BuilderID         *int                   `json:"builder_id,omitempty"`
	ChannelID         string                 `json:"channel_id,omitempty"`
	Priority          *int                   `json:"priority,omitempty"`
	Category          string                 `json:"category,omitempty"`
	Style             *int                   `json:"style,omitempty"`
	BigText           string                 `json:"big_text,omitempty"`
	Inbox             map[string]interface{} `json:"inbox,omitempty"`
	BigPicPath        string                 `json:"big_pic_path,omitempty"`
	Extras            map[string]interface{} `json:"extras,omitempty"`
	Intent            *AndroidIntent         `json:"intent,omitempty"`
	LargeIcon         string                 `json:"large_icon,omitempty"`
	SmallIcon         string                 `json:"small_icon,omitempty"`
	Sound             string                 `json:"sound,omitempty"`
	BadgeAddNum       *int                   `json:"badge_add_num,omitempty"`
	BadgeClass        string                 `json:"badge_class,omitempty"`
	DisplayForeground string                 `json:"display_foreground,omitempty"`
	GroupID           string                 `json:"group_id,omitempty"`
}

type AndroidIntent struct {
	URL string `json:"url,omitempty"`
}

type IOSNotification struct {
	Alert             interface{}            `json:"alert,omitempty"`
	Sound             interface{}            `json:"sound,omitempty"`
	Badge             interface{}            `json:"badge,omitempty"`
	ContentAvailable  *bool                  `json:"content-available,omitempty"`
	MutableContent    *bool                  `json:"mutable-content,omitempty"`
	Category          string                 `json:"category,omitempty"`
	Extras            map[string]interface{} `json:"extras,omitempty"`
	ThreadID          string                 `json:"thread-id,omitempty"`
	InterruptionLevel string                 `json:"interruption-level,omitempty"`
}

type HmosNotification struct {
	Alert             string                 `json:"alert,omitempty"`
	Title             string                 `json:"title,omitempty"`
	Category          string                 `json:"category,omitempty"`
	LargeIcon         string                 `json:"large_icon,omitempty"`
	Intent            *HmosIntent            `json:"intent,omitempty"`
	BadgeAddNum       *int                   `json:"badge_add_num,omitempty"`
	BadgeSetNum       *int                   `json:"badge_set_num,omitempty"`
	TestMessage       *bool                  `json:"test_message,omitempty"`
	ReceiptID         string                 `json:"receipt_id,omitempty"`
	Extras            map[string]interface{} `json:"extras,omitempty"`
	Style             *int                   `json:"style,omitempty"`
	InboxContent      []string               `json:"inbox_content,omitempty"`
	PushType          *int                   `json:"push_type,omitempty"`
	ExtraData         string                 `json:"extra_data,omitempty"`
	DisplayForeground string                 `json:"display_foreground,omitempty"`
}

type HmosIntent struct {
	URL string `json:"url,omitempty"`
}

// CustomMessage is the custom/passthrough message payload.
type CustomMessage struct {
	Title       string                 `json:"title,omitempty"`
	MsgContent  string                 `json:"msg_content,omitempty"`
	ContentType string                 `json:"content_type,omitempty"`
	Extras      map[string]interface{} `json:"extras,omitempty"`
}

// LiveActivityMessage is for iOS Live Activity.
type LiveActivityMessage struct {
	IOS *LiveActivityIOS `json:"ios,omitempty"`
}

type LiveActivityIOS struct {
	Event          string                 `json:"event,omitempty"`
	AttributesType string                 `json:"attributes-type,omitempty"`
	ContentState   map[string]interface{} `json:"content-state,omitempty"`
	Alert          *LiveActivityAlert     `json:"alert,omitempty"`
	DismissalDate  *int64                 `json:"dismissal-date,omitempty"`
	StaleDate      *int64                 `json:"stale-date,omitempty"`
	Attributes     map[string]interface{} `json:"attributes,omitempty"`
	RelevanceScore *int                   `json:"relevance-score,omitempty"`
	APNSPriority   *int                   `json:"apns-priority,omitempty"`
}

type LiveActivityAlert struct {
	Sound string `json:"sound,omitempty"`
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

// Options controls push behavior.
type Options struct {
	TimeToLive        *int                   `json:"time_to_live,omitempty"`
	OverrideMsgID     *int64                 `json:"override_msg_id,omitempty"`
	APNSProduction    *bool                  `json:"apns_production,omitempty"`
	APNSCollapseID    string                 `json:"apns_collapse_id,omitempty"`
	BigPushDuration   *int                   `json:"big_push_duration,omitempty"`
	MultiLanguage     map[string]interface{} `json:"multi_language,omitempty"`
	ThirdPartyChannel map[string]interface{} `json:"third_party_channel,omitempty"`
	Classification    *int                   `json:"classification,omitempty"`
	VoiceValue        string                 `json:"voice_value,omitempty"`
	EnhancMessage     *bool                  `json:"enhanc_message,omitempty"`
	PlanID            string                 `json:"plan_id,omitempty"`
	CID               string                 `json:"cid,omitempty"`
}

// PushResult is the response for a push request.
type PushResult struct {
	RequestID string `json:"request_id,omitempty"`
	MsgID     string `json:"msg_id"`
}

// PushWithdrawResult is the response for a push withdraw request.
type PushWithdrawResult struct {
	RequestID string `json:"request_id,omitempty"`
	MsgID     string `json:"msg_id,omitempty"`
}

// BatchPushParam for batch push by regid or alias.
type BatchPushParam struct {
	Requests []BatchPushRequest `json:"requests"`
}

type BatchPushRequest struct {
	Target       string                 `json:"target,omitempty"`
	Platform     interface{}            `json:"platform,omitempty"` // string or []string
	Notification *NotificationMessage   `json:"notification,omitempty"`
	Message      *CustomMessage         `json:"message,omitempty"`
	Options      *Options               `json:"options,omitempty"`
	CustomArgs   map[string]interface{} `json:"custom_args,omitempty"`
}

type BatchPushResult struct {
	Results map[string]BatchPushSingleResult `json:"results,omitempty"`
}

type BatchPushSingleResult struct {
	Target  string `json:"target,omitempty"`
	Success bool   `json:"success"`
	MsgID   int64  `json:"msg_id"`
}

// --- PushService ---

type PushService struct {
	client *Client
}

// Send creates and sends a push notification.
// POST /v4/push
func (s *PushService) Send(ctx context.Context, param *PushParam) (*PushResult, error) {
	var result PushResult
	err := s.client.doPost(ctx, "/v4/push", param, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SendRaw sends a push with a raw JSON body (map or any serializable object).
// POST /v4/push
func (s *PushService) SendRaw(ctx context.Context, param interface{}) (*PushResult, error) {
	var result PushResult
	err := s.client.doPost(ctx, "/v4/push", param, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Validate validates a push payload without actually sending it.
// POST /v4/push/validate
func (s *PushService) Validate(ctx context.Context, param *PushParam) (*PushResult, error) {
	var result PushResult
	err := s.client.doPost(ctx, "/v4/push/validate", param, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Withdraw withdraws a push message by message ID.
// Only messages within 24 hours can be withdrawn. No duplicate withdrawal.
// DELETE /v4/push/withdraw/{msgId}
func (s *PushService) Withdraw(ctx context.Context, msgID string) (*PushWithdrawResult, error) {
	var result PushWithdrawResult
	err := s.client.doDelete(ctx, fmt.Sprintf("/v4/push/withdraw/%s", msgID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BatchByRegID sends batch push by registration IDs.
// POST /v4/batch/push/regid
func (s *PushService) BatchByRegID(ctx context.Context, param *BatchPushParam) (*BatchPushResult, error) {
	var result BatchPushResult
	err := s.client.doPost(ctx, "/v4/batch/push/regid", param, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BatchByAlias sends batch push by aliases.
// POST /v4/batch/push/alias
func (s *PushService) BatchByAlias(ctx context.Context, param *BatchPushParam) (*BatchPushResult, error) {
	var result BatchPushResult
	err := s.client.doPost(ctx, "/v4/batch/push/alias", param, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
