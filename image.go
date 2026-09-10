package engagelab

import (
	"context"
	"errors"
)

type ImageService struct {
	client *Client
}

type OppoImageParam struct {
	BigPictureURL   string `json:"big_picture_url,omitempty"`
	SmallPictureURL string `json:"small_picture_url,omitempty"`
}

type ImageUploadResult struct {
	BigPictureID   string `json:"big_picture_id,omitempty"`
	SmallPictureID string `json:"small_picture_id,omitempty"`
}

// UploadOppo registers image URLs for OPPO notification pushes.
// POST /v4/image/oppo
func (s *ImageService) UploadOppo(ctx context.Context, param *OppoImageParam) (*ImageUploadResult, error) {
	if param == nil {
		return nil, errors.New("oppo image param is required")
	}
	var result ImageUploadResult
	if err := s.client.doPost(ctx, "/v4/image/oppo", param, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
