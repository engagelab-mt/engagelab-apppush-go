package engagelab

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type VoiceService struct {
	client *Client
}

type VoiceResult struct {
	Language string `json:"language,omitempty"`
	FileURL  string `json:"file_url,omitempty"`
}

// Create uploads a voice file for one language.
// POST /v4/voices
func (s *VoiceService) Create(ctx context.Context, language, filePath string) (*VoiceResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open voice file: %w", err)
	}
	defer func() { _ = file.Close() }()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("language", language); err != nil {
		return nil, fmt.Errorf("write language field: %w", err)
	}
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("create voice file field: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("copy voice file: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.client.baseURL+"/v4/voices", body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", s.client.authHeader)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseApiError(resp.StatusCode, respBody)
	}
	var result VoiceResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &result, nil
}

func (s *VoiceService) List(ctx context.Context) ([]VoiceResult, error) {
	var result []VoiceResult
	if err := s.client.doGet(ctx, "/v4/voices", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *VoiceService) Get(ctx context.Context, language string) (*VoiceResult, error) {
	var result VoiceResult
	if err := s.client.doGet(ctx, fmt.Sprintf("/v4/voices/%s", language), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *VoiceService) Delete(ctx context.Context, language string) error {
	return s.client.doDelete(ctx, fmt.Sprintf("/v4/voices/%s", language), nil)
}
