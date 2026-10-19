package engagelab

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

// GroupPushParam for group push (inherits PushParam structure).
type GroupPushParam = PushParam

// GroupPushResult is the response for a group push request.
type GroupPushResult struct {
	GroupMsgID string                          `json:"group_msgid,omitempty"`
	Successes  map[string]PushResult           `json:"-"`
	Errors     map[string]GroupPushErrorDetail `json:"-"`
}

func (r *GroupPushResult) UnmarshalJSON(data []byte) error {
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	r.Successes = make(map[string]PushResult)
	r.Errors = make(map[string]GroupPushErrorDetail)
	for key, raw := range payload {
		if key == "group_msgid" {
			if err := json.Unmarshal(raw, &r.GroupMsgID); err != nil {
				return err
			}
			continue
		}
		var apiError struct {
			Error *GroupPushErrorDetail `json:"error"`
		}
		if err := json.Unmarshal(raw, &apiError); err == nil && apiError.Error != nil {
			r.Errors[key] = *apiError.Error
			continue
		}
		var success PushResult
		if err := json.Unmarshal(raw, &success); err != nil {
			return err
		}
		r.Successes[key] = success
	}
	return nil
}

type GroupPushErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GroupPushClient is a separate client for Group Push API.
// It uses group-{groupKey}:{groupMasterSecret} for Basic Auth.
type GroupPushClient struct {
	httpClient *http.Client
	baseURL    string
	authHeader string
}

// NewGroupPushClient creates a client for EngageLab Group Push API.
// Authentication uses group-{groupKey}:{groupMasterSecret}.
func NewGroupPushClient(groupKey, groupMasterSecret string, opts ...Option) *GroupPushClient {
	authStr := "group-" + groupKey + ":" + groupMasterSecret
	gc := &GroupPushClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    string(Singapore),
		authHeader: "Basic " + base64.StdEncoding.EncodeToString([]byte(authStr)),
	}

	// Apply options via a temporary Client wrapper
	temp := &Client{
		httpClient: gc.httpClient,
		baseURL:    gc.baseURL,
	}
	for _, opt := range opts {
		opt(temp)
	}
	gc.httpClient = temp.httpClient
	gc.baseURL = temp.baseURL

	return gc
}

// Send sends a group push to multiple apps.
// POST /v4/grouppush
func (gc *GroupPushClient) Send(ctx context.Context, param *GroupPushParam) (*GroupPushResult, error) {
	c := &Client{
		httpClient: gc.httpClient,
		baseURL:    gc.baseURL,
		authHeader: gc.authHeader,
	}
	var result GroupPushResult
	err := c.doPost(ctx, "/v4/grouppush", param, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
