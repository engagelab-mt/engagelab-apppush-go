package engagelab

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type DataCenter string

const (
	Singapore DataCenter = "https://pushapi-sgp.engagelab.com"
	HongKong  DataCenter = "https://pushapi-hk.engagelab.com"
	Virginia  DataCenter = "https://pushapi-usva.engagelab.com"
	Frankfurt DataCenter = "https://pushapi-defra.engagelab.com"
	Japan     DataCenter = "https://pushapi-jpn.engagelab.com"
	Brazil    DataCenter = "https://pushapi-bra.engagelab.com"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	authHeader string

	Push     *PushService
	Device   *DeviceService
	Tag      *TagService
	Alias    *AliasService
	Schedule *ScheduleService
	Status   *StatusService
	Plan     *PlanService
	Voice    *VoiceService
	Image    *ImageService
	App      *AppService
}

type Option func(*Client)

func WithDataCenter(dc DataCenter) Option {
	return func(c *Client) {
		c.baseURL = string(dc)
	}
}

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// NewClient creates a new EngageLab AppPush client.
// appKey and masterSecret are required. Default data center is Singapore.
func NewClient(appKey, masterSecret string, opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    string(Singapore),
		authHeader: basicAuth(appKey, masterSecret),
	}

	for _, opt := range opts {
		opt(c)
	}

	c.Push = &PushService{client: c}
	c.Device = &DeviceService{client: c}
	c.Tag = &TagService{client: c}
	c.Alias = &AliasService{client: c}
	c.Schedule = &ScheduleService{client: c}
	c.Status = &StatusService{client: c}
	c.Plan = &PlanService{client: c}
	c.Voice = &VoiceService{client: c}
	c.Image = &ImageService{client: c}
	c.App = &AppService{client: c}

	return c
}

func basicAuth(username, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	fullURL := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseApiError(resp.StatusCode, respBody)
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}

func (c *Client) doGet(ctx context.Context, path string, query url.Values, result interface{}) error {
	if len(query) > 0 {
		path = path + "?" + query.Encode()
	}
	return c.doRequest(ctx, http.MethodGet, path, nil, result)
}

func (c *Client) doPost(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.doRequest(ctx, http.MethodPost, path, body, result)
}

func (c *Client) doPut(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.doRequest(ctx, http.MethodPut, path, body, result)
}

func (c *Client) doDelete(ctx context.Context, path string, result interface{}) error {
	return c.doRequest(ctx, http.MethodDelete, path, nil, result)
}

func (c *Client) doDeleteWithQuery(ctx context.Context, path string, query url.Values, result interface{}) error {
	if len(query) > 0 {
		path = path + "?" + query.Encode()
	}
	return c.doRequest(ctx, http.MethodDelete, path, nil, result)
}
