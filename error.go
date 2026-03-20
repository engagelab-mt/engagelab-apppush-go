package engagelab

import (
	"encoding/json"
	"fmt"
)

type ApiError struct {
	StatusCode int
	ErrorBody  ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("engagelab api error: status=%d code=%d message=%s",
		e.StatusCode, e.ErrorBody.Code, e.ErrorBody.Message)
}

func parseApiError(statusCode int, body []byte) *ApiError {
	apiErr := &ApiError{StatusCode: statusCode}
	var wrapper struct {
		Error ErrorDetail `json:"error"`
	}
	if err := json.Unmarshal(body, &wrapper); err == nil {
		apiErr.ErrorBody = wrapper.Error
	}
	return apiErr
}
