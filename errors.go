package postghost

import "fmt"

type ErrorDetail struct {
	Code    string `json:"code"`
	Slug    string `json:"slug"`
	Message string `json:"message"`
}

type APIError struct {
	StatusCode int
	Detail     ErrorDetail
	RawBody    string
}

func (e *APIError) Error() string {
	if e == nil {
		return "api error"
	}
	if e.Detail.Message != "" {
		return fmt.Sprintf("postghost api error (%d): %s", e.StatusCode, e.Detail.Message)
	}
	return fmt.Sprintf("postghost api error (%d)", e.StatusCode)
}
