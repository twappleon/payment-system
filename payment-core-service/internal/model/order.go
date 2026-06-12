package model

import "time"

type PaymentOrder struct {
	MerchantNo      string    `json:"merchant_no"`
	MerchantOrderNo string    `json:"merchant_order_no"`
	PlatformOrderNo string    `json:"platform_order_no"`
	ChannelCode     string    `json:"channel_code,omitempty"`
	ChannelOrderNo  string    `json:"channel_order_no,omitempty"`
	Amount          string    `json:"amount"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	NotifyStatus    string    `json:"notify_status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

