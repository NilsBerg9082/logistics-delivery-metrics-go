package main

import (
	"fmt"
	"time"
)

type Delivery struct {
	ID          string
	Region      string
	Service     string
	DeliveredAt time.Time
	PromisedAt  time.Time
}

func reportDelivery(client *MetricsClient, delivery Delivery) error {
	status := "on_time"
	if delivery.DeliveredAt.After(delivery.PromisedAt) {
		status = "late"
	}
	payload := map[string]any{
		"type":            "counter",
		"name":            "logistics.delivery.completed",
		"value":           1,
		"tags":            map[string]string{"region": delivery.Region, "service": delivery.Service, "status": status},
		"idempotency_key": "delivery:" + delivery.ID,
	}
	if err := client.Report(payload); err != nil {
		return fmt.Errorf("report delivery %s: %w", delivery.ID, err)
	}
	return nil
}

func main() {
	client, err := NewMetricsClient()
	if err != nil {
		panic(err)
	}
	delivery := Delivery{
		ID: "pkg-2026-0007", Region: "east", Service: "standard",
		DeliveredAt: time.Date(2026, 8, 9, 15, 30, 0, 0, time.UTC),
		PromisedAt:  time.Date(2026, 8, 9, 16, 0, 0, 0, time.UTC),
	}
	if err := reportDelivery(client, delivery); err != nil {
		panic(err)
	}
	fmt.Println("reported delivery.completed for", delivery.ID)
}
