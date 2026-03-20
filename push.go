package engagelab

import (
	"context"
	"fmt"
)

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
