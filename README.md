# PostGhost Go SDK

Typed Go SDK for PostGhost ingestion endpoints.

## Install

```bash
go get postghost-sdk-go
```

## Initialize Client

```go
package main

import (
	"context"
	"log"

	postghost "postghost-sdk-go"
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
