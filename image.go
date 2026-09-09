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

// UploadOppo registers one image URL for OPPO notification pushes.
// Exactly one of big_picture_url and small_picture_url must be provided.
// POST /v4/image/oppo
func (s *ImageService) UploadOppo(ctx context.Context, param *OppoImageParam) (*ImageUploadResult, error) {
	if param == nil || (param.BigPictureURL == "") == (param.SmallPictureURL == "") {
		return nil, errors.New("exactly one of big_picture_url and small_picture_url is required")
	}
	var result ImageUploadResult
	if err := s.client.doPost(ctx, "/v4/image/oppo", param, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
