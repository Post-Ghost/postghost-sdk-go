# PostGhost Go SDK

Typed Go SDK for PostGhost ingestion endpoints.

## Install

```bash
go get github.com/Post-Ghost/postghost-sdk-go
```

## Initialize Client

```go
package main

import (
	"context"
	"log"

	postghost "github.com/Post-Ghost/postghost-sdk-go"
)

func main() {
	client, err := postghost.NewClient("pgk_xxx", false) // false => https://api.postghost.dev/api/v1
	if err != nil {
		log.Fatal(err)
	}

	latency := 120
	_, err = client.IngestAPIMonitoring(context.Background(), postghost.ApiMonitoringPayload{
		ServiceName: "stripe",
		Method:      "GET",
		URL:         "https://api.stripe.com/v1/customers",
		StatusCode:  200,
		LatencyMS:   &latency,
	})
	if err != nil {
		log.Fatal(err)
	}
}
```

Use local mode:

```go
client, err := postghost.NewClient("pgk_xxx", true) // true => http://localhost:8080/api/v1
```

On client initialization, the SDK performs an auth check:

- Sends `GET /auth-check` to the configured base URL.
- Expects HTTP 204 when the API key is valid.
- For 4xx responses with a JSON body containing `message`, the constructor returns an error and no client is created.

SDK metadata headers are set automatically by this SDK on every request:
- `X-PostGhost-SDK-Name`
- `X-PostGhost-SDK-Version`

## Available Methods

- `IngestPulse(ctx, externalID, payload)`
- `IngestAPIMonitoring(ctx, payload)`
- `IngestAPIMonitoringBatch(ctx, payloads)`
- `IngestLog(ctx, payload)`
- `IngestLogsBatch(ctx, payloads)`

All methods accept typed structs and return typed responses.
