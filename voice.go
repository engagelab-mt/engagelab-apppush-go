package engagelab

import (
	"context"
	"fmt"
)

type VoiceService struct {
	client *Client
}

// VoiceParam is the request for creating a voice/TTS template.
type VoiceParam struct {
	Language string `json:"language,omitempty"`
	Content  string `json:"content,omitempty"`
	TTSType  string `json:"tts_type,omitempty"`
}

// VoiceResult is the response for voice operations.
type VoiceResult struct {
	Language string `json:"language,omitempty"`
	Content  string `json:"content,omitempty"`
	TTSType  string `json:"tts_type,omitempty"`
}

// VoiceListResult is the list of voice templates.
type VoiceListResult struct {
	Voices []VoiceResult `json:"voices,omitempty"`
}

// Create creates a new voice/TTS template.
// POST /v4/voices
func (s *VoiceService) Create(ctx context.Context, param *VoiceParam) (*VoiceResult, error) {
	var result VoiceResult
	err := s.client.doPost(ctx, "/v4/voices", param, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// List returns all voice/TTS templates.
// GET /v4/voices
func (s *VoiceService) List(ctx context.Context) (*VoiceListResult, error) {
	var result VoiceListResult
	err := s.client.doGet(ctx, "/v4/voices", nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a voice/TTS template by language.
// GET /v4/voices/{language}
func (s *VoiceService) Get(ctx context.Context, language string) (*VoiceResult, error) {
	var result VoiceResult
	err := s.client.doGet(ctx, fmt.Sprintf("/v4/voices/%s", language), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes a voice/TTS template by language.
// DELETE /v4/voices/{language}
func (s *VoiceService) Delete(ctx context.Context, language string) error {
	return s.client.doDelete(ctx, fmt.Sprintf("/v4/voices/%s", language), nil)
}
