package postghost

import (
	"fmt"
	"strings"
)

type PulsePayload map[string]any

type PulseAcceptedResponse struct {
	Status  string `json:"status"`
	EventID string `json:"eventId"`
}

type ApiMonitoringPayload struct {
	ServiceName  string         `json:"serviceName"`
	Method       string         `json:"method"`
	URL          string         `json:"url"`
	StatusCode   int            `json:"statusCode"`
	Request      map[string]any `json:"request,omitempty"`
	Response     map[string]any `json:"response,omitempty"`
	Meta         map[string]any `json:"meta,omitempty"`
	LatencyMS    *int           `json:"latencyMs,omitempty"`
	ErrorMessage *string        `json:"errorMessage,omitempty"`
}

type ApiMonitoringAcceptedResponse struct {
	Status   string   `json:"status"`
	EventID  string   `json:"eventId"`
	EventIDs []string `json:"eventIds,omitempty"`
}

type LogsPayload struct {
	Level      string         `json:"level,omitempty"`
	Message    string         `json:"message,omitempty"`
	Timestamp  string         `json:"timestamp,omitempty"`
	Attributes map[string]any `json:"attributes,omitempty"`
	JobID      string         `json:"jobId,omitempty"`
	Source     string         `json:"source,omitempty"`
	RawLog     string         `json:"rawLog,omitempty"`
}

type LogsAcceptedResponse struct {
	Status   string   `json:"status"`
	EventID  string   `json:"eventId"`
	EventIDs []string `json:"eventIds"`
}

func (p ApiMonitoringPayload) validate() error {
	if strings.TrimSpace(p.ServiceName) == "" {
		return fmt.Errorf("serviceName is required")
	}
	if strings.TrimSpace(p.Method) == "" {
		return fmt.Errorf("method is required")
	}
	if strings.TrimSpace(p.URL) == "" {
		return fmt.Errorf("url is required")
	}
	if p.StatusCode == 0 {
		return fmt.Errorf("statusCode is required")
	}
	return nil
}

func (p LogsPayload) validate() error {
	if strings.TrimSpace(p.RawLog) != "" {
		return nil
	}
	if strings.TrimSpace(p.Level) == "" {
		return fmt.Errorf("level is required when rawLog is empty")
	}
	if strings.TrimSpace(p.Message) == "" {
		return fmt.Errorf("message is required when rawLog is empty")
	}
	return nil
}
