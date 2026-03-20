package engagelab

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

type TagService struct {
	client *Client
}

// List returns all tags for the application.
// GET /v4/tags
func (s *TagService) List(ctx context.Context) (*TagsGetResult, error) {
	var result TagsGetResult
	err := s.client.doGet(ctx, "/v4/tags", nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Set adds or removes registration IDs for a tag.
// POST /v4/tags/{tag}
func (s *TagService) Set(ctx context.Context, tag string, param *TagSetParam) error {
	return s.client.doPost(ctx, fmt.Sprintf("/v4/tags/%s", tag), param, nil)
}

// Delete removes a tag. platforms is a list like ["android", "ios", "hmos"].
// DELETE /v4/tags/{tag}?platform={platforms}
func (s *TagService) Delete(ctx context.Context, tag string, platforms []string) error {
	query := url.Values{}
	if len(platforms) > 0 {
		query.Set("platform", strings.Join(platforms, ","))
	}
	return s.client.doDeleteWithQuery(ctx, fmt.Sprintf("/v4/tags/%s", tag), query, nil)
}

// GetCount returns the count of devices for given tags and platforms.
// GET /v4/tags_count?tags={tags}&platform={platforms}
func (s *TagService) GetCount(ctx context.Context, tags []string, platforms []string) (*TagsCountGetResult, error) {
	query := url.Values{}
	query.Set("tags", strings.Join(tags, ","))
	if len(platforms) > 0 {
		query.Set("platform", strings.Join(platforms, ","))
	}
	var result TagsCountGetResult
	err := s.client.doGet(ctx, "/v4/tags_count", query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDeviceStatus checks if a registration ID has a specific tag.
// GET /v4/tags/{tag}/registration_ids/{registration_id}
func (s *TagService) GetDeviceStatus(ctx context.Context, tag, registrationID string) (*TagsGetResult, error) {
	var result TagsGetResult
	err := s.client.doGet(ctx, fmt.Sprintf("/v4/tags/%s/registration_ids/%s", tag, registrationID), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetQuota returns tag/alias quota information.
// GET /v4/tags/quota-info?tags={tags}&platform={platforms}
func (s *TagService) GetQuota(ctx context.Context, tags []string, platforms []string) (*TagQuotaGetResult, error) {
	query := url.Values{}
	if len(tags) > 0 {
		query.Set("tags", strings.Join(tags, ","))
	}
	if len(platforms) > 0 {
		query.Set("platform", strings.Join(platforms, ","))
	}
	var result TagQuotaGetResult
	err := s.client.doGet(ctx, "/v4/tags/quota-info", query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
