package engagelab

import "context"

type AppService struct {
	client *Client
}

type AppVIPStatusResult struct {
	VIPStatus  int   `json:"vip_status"`
	VIPEndTime int64 `json:"vip_end_time"`
}

// GetVIPStatus returns the application's VIP status.
// GET /v4/app/vip/status
func (s *AppService) GetVIPStatus(ctx context.Context) (*AppVIPStatusResult, error) {
	var result AppVIPStatusResult
	if err := s.client.doGet(ctx, "/v4/app/vip/status", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
