package engagelab

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

type AliasService struct {
	client *Client
}

// Get returns the registration IDs associated with an alias.
// GET /v4/aliases/{alias}?platform={platforms}
func (s *AliasService) Get(ctx context.Context, alias string, platforms []string) (*AliasStatusGetResult, error) {
	query := url.Values{}
	if len(platforms) > 0 {
		query.Set("platform", strings.Join(platforms, ","))
	}
	var result AliasStatusGetResult
	err := s.client.doGet(ctx, fmt.Sprintf("/v4/aliases/%s", alias), query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes an alias binding for given platforms.
// DELETE /v4/aliases/{alias}?platform={platforms}
func (s *AliasService) Delete(ctx context.Context, alias string, platforms []string) error {
	query := url.Values{}
	if len(platforms) > 0 {
		query.Set("platform", strings.Join(platforms, ","))
	}
	return s.client.doDeleteWithQuery(ctx, fmt.Sprintf("/v4/aliases/%s", alias), query, nil)
}
