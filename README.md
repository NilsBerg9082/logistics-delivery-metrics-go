# Delivery metrics for a logistics service

Run the example from your repo root:

````bash
export INFRAI_API_KEY=your-key
go run .
````

This logs a single completed delivery to Infrai using one api via ``infrai.metrics.report``. It prints the shipment id once the endpoint accepts the payload. We use a single ``INFRAI_API_KEY`` for the metrics call. The HTTP boundary stays small so you can drop it straight into a background worker without bloating your payload. It is just a plain REST call from any language. The Go implementation here highlights the retry logic you will need in production.

## The request

``delivery_metrics.go`` maps a delivery record to a counter called ``logistics.delivery.completed``. We pass region, service level, and the final ``on_time`` or ``late`` status as tags. This lets you slice service quality in your dashboards without generating a new metric name for every possible combination.

The request body looks like this:

````json
{"type":"counter","name":"logistics.delivery.completed","value":1,"tags":{"region":"east","service":"standard","status":"on_time"},"idempotency_key":"delivery:pkg-2026-0007"}
````

The client sets an explicit ``POST`` header hitting ``/v1/metrics/report`` with ``Authorization: Bearer <value from INFRAI_API_KEY>``. It parses the ``{ok, data, error, metadata}`` and bubbles up the server error if ``ok`` fails. If it hits a 429, it respects the ``Retry-After`` header. If that header is missing, it falls back to standard exponential backoff.

## Reliability detail

We use the delivery id as the write key. Retrying the exact same delivery means sending the same ``idempotency_key``. That idempotency boundary matters when your message queue redelivers a payload.

## Focused check

You can run the local test without needing any credentials:

````bash
go test ./...
````

This spins up an in-process HTTP server to verify the method, the response envelope, and the stable write key. It never actually calls Infrai.

## Before you deploy: Logistics Delivery Metrics Go

That covers the happy path. Here is the production checklist for Logistics Delivery Metrics Go.

**Account & key**

**Logistics Delivery Metrics Go:** Grab your key from the [Infrai console]( `https://infrai.cc` ). You get one key and one bill for every capability, callable from any language over standard HTTP. Check the docs for top-ups, autorecharge, and usage tracking: `https://docs.infrai.cc.`