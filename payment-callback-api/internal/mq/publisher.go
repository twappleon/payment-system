package mq

import (
	"context"
	"encoding/json"
	"log/slog"
)

type Event struct {
	EventID         string `json:"event_id"`
	EventType       string `json:"event_type"`
	ChannelCode     string `json:"channel_code"`
	PlatformOrderNo string `json:"platform_order_no"`
	ChannelOrderNo  string `json:"channel_order_no"`
	OccurredAt      string `json:"occurred_at"`
}

type Publisher struct {
	logger *slog.Logger
}

func NewPublisher(logger *slog.Logger) *Publisher {
	return &Publisher{logger: logger}
}

func (p *Publisher) Publish(ctx context.Context, routingKey string, event Event) error {
	payload, _ := json.Marshal(event)
	p.logger.InfoContext(ctx, "publish event placeholder", "routing_key", routingKey, "payload", string(payload))
	return nil
}

