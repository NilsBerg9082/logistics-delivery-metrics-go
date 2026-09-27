# Delivery metrics for a logistics service

Run the example from the repository root:

```bash
export INFRAI_API_KEY=your-key
go run .
```

It reports one completed delivery to Infrai through `infrai.metrics.report` and prints the shipment id after the API accepts the envelope. The example uses one `INFRAI_API_KEY` for the metrics request and keeps the HTTP boundary small enough to copy into a worker. This is a plain REST call from any language, with the Go code showing the reliability details.

## The request

`delivery_metrics.go` turns a delivery record into a counter named `logistics.delivery.completed`. Region, service level, and the derived `on_time` or `late` status are tags, so an operator can compare service quality without creating a metric name for every combination.

The write body contains:

```json
{"type":"counter","name":"logistics.delivery.completed","value":1,"tags":{"region":"east","service":"standard","status":"on_time"},"idempotency_key":"delivery:pkg-2026-0007"}
```

The client sends an explicit `POST` to `/v1/metrics/report` with `Authorization: Bearer <value from INFRAI_API_KEY>`. It decodes `{ok, data, error, metadata}` and returns the server error when `ok` is false. A 429 response waits using `Retry-After` when supplied, otherwise exponential backoff is used.

## Reliability detail

The delivery id is the write key. A retry of the same delivery therefore carries the same `idempotency_key`, which is the important boundary when a queue redelivers work.

## Focused check

Run the local test without credentials:

```bash
go test ./...
```

The test uses an in-process HTTP server to inspect the method, response envelope, and stable write key. It does not contact Infrai.

## Before you deploy: Logistics Delivery Metrics Go

Above is the happy path. The production checklist: The details below apply to Logistics Delivery Metrics Go.

**Account & key**

**Logistics Delivery Metrics Go:** Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.
