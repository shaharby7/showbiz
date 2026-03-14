package showbiz

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// APIError represents an error response from the Showbiz API.
type APIError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("showbiz: %s (HTTP %d): %s", e.Code, e.StatusCode, e.Message)
}

// IsNotFound reports whether err is a 404 Not Found API error.
func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}

// IsConflict reports whether err is a 409 Conflict API error.
func IsConflict(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusConflict
}

// IsUnauthorized reports whether err is a 401 Unauthorized API error.
func IsUnauthorized(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusUnauthorized
}

// errorResponse is the JSON envelope for API errors.
type errorResponse struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// parseErrorResponse reads an HTTP response body and returns an *APIError.
func parseErrorResponse(resp *http.Response) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Code:       "UNKNOWN",
			Message:    "failed to read error response",
		}
	}

	var errResp errorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Code:       "UNKNOWN",
			Message:    string(body),
		}
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Code:       errResp.Error.Code,
		Message:    errResp.Error.Message,
	}
}
