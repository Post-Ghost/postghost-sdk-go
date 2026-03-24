package postghost

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultBaseURL      = "https://api.postghost.dev/api/v1"
	defaultLocalBaseURL = "http://localhost:8080/api/v1"
	defaultTimeout      = 30 * time.Second
	defaultSDKName      = "go"
	defaultSDKVersion   = "1.0.0"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	sdkName    string
	sdkVersion string
}

func NewClient(apiKey string, useLocal bool) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("api key is required")
	}

	baseURL := defaultBaseURL
	if useLocal {
		baseURL = defaultLocalBaseURL
	}

	httpClient := &http.Client{Timeout: defaultTimeout}

	return &Client{
		apiKey:     apiKey,
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: httpClient,
		sdkName:    defaultSDKName,
		sdkVersion: defaultSDKVersion,
	}, nil
}

func (c *Client) IngestPulse(ctx context.Context, externalID string, payload *PulsePayload) (*PulseAcceptedResponse, error) {
	if strings.TrimSpace(externalID) == "" {
		return nil, fmt.Errorf("externalID is required")
	}
	path := fmt.Sprintf("/ingestion/pulse/%s", externalID)
	return doJSON[PulseAcceptedResponse](ctx, c, http.MethodPost, path, payload)
}

func (c *Client) IngestAPIMonitoring(ctx context.Context, payload ApiMonitoringPayload) (*ApiMonitoringAcceptedResponse, error) {
	if err := payload.validate(); err != nil {
		return nil, err
	}
	return doJSON[ApiMonitoringAcceptedResponse](ctx, c, http.MethodPost, "/ingestion/api-monitoring", payload)
}

func (c *Client) IngestAPIMonitoringBatch(ctx context.Context, payloads []ApiMonitoringPayload) (*ApiMonitoringAcceptedResponse, error) {
	if len(payloads) == 0 {
		return nil, fmt.Errorf("payloads must contain at least one event")
	}
	for i := range payloads {
		if err := payloads[i].validate(); err != nil {
			return nil, fmt.Errorf("payloads[%d]: %w", i, err)
		}
	}
	return doJSON[ApiMonitoringAcceptedResponse](ctx, c, http.MethodPost, "/ingestion/api-monitoring", payloads)
}

func (c *Client) IngestLog(ctx context.Context, payload LogsPayload) (*LogsAcceptedResponse, error) {
	if err := payload.validate(); err != nil {
		return nil, err
	}
	return doJSON[LogsAcceptedResponse](ctx, c, http.MethodPost, "/ingestion/logs", payload)
}

func (c *Client) IngestLogsBatch(ctx context.Context, payloads []LogsPayload) (*LogsAcceptedResponse, error) {
	if len(payloads) == 0 {
		return nil, fmt.Errorf("payloads must contain at least one event")
	}
	for i := range payloads {
		if err := payloads[i].validate(); err != nil {
			return nil, fmt.Errorf("payloads[%d]: %w", i, err)
		}
	}
	return doJSON[LogsAcceptedResponse](ctx, c, http.MethodPost, "/ingestion/logs", payloads)
}

func doJSON[T any](ctx context.Context, c *Client, method, path string, body any) (*T, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", c.apiKey)
	if c.sdkName != "" {
		req.Header.Set("X-PostGhost-SDK-Name", c.sdkName)
	}
	if c.sdkVersion != "" {
		req.Header.Set("X-PostGhost-SDK-Version", c.sdkVersion)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{StatusCode: resp.StatusCode, RawBody: string(respBody)}
		_ = json.Unmarshal(respBody, &apiErr.Detail)
		return nil, apiErr
	}

	var out T
	if len(respBody) == 0 {
		return &out, nil
	}
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response body: %w", err)
	}
	return &out, nil
}
