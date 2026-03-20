package engagelab

import (
	"context"
	"fmt"
)

type DeviceService struct {
	client *Client
}

// GetStatus queries online status for a list of registration IDs.
// POST /v4/devices/status
func (s *DeviceService) GetStatus(ctx context.Context, param *DeviceStatusGetParam) ([]DeviceStatusGetResult, error) {
	var result []DeviceStatusGetResult
	err := s.client.doPost(ctx, "/v4/devices/status", param, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Get retrieves device info (tags + alias) by registration ID.
// GET /v4/devices/{registration_id}
func (s *DeviceService) Get(ctx context.Context, registrationID string) (*DeviceGetResult, error) {
	var result DeviceGetResult
	err := s.client.doGet(ctx, fmt.Sprintf("/v4/devices/%s", registrationID), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Set updates tags and alias for a device.
// POST /v4/devices/{registration_id}
func (s *DeviceService) Set(ctx context.Context, registrationID string, param *DeviceSetParam) error {
	return s.client.doPost(ctx, fmt.Sprintf("/v4/devices/%s", registrationID), param, nil)
}

// Delete removes a device and all its associated data (tags, alias, etc.).
// This operation is asynchronous and irreversible.
// DELETE /v4/devices/{registration_id}
func (s *DeviceService) Delete(ctx context.Context, registrationID string) error {
	return s.client.doDelete(ctx, fmt.Sprintf("/v4/devices/%s", registrationID), nil)
}
