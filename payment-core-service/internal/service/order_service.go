package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/company/payment-core-service/internal/model"
	"github.com/company/payment-core-service/internal/statemachine"
)

type CreatePaymentOrderRequest struct {
	MerchantNo      string `json:"merchant_no"`
	MerchantOrderNo string `json:"merchant_order_no"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	NotifyURL       string `json:"notify_url"`
	ClientIP        string `json:"client_ip"`
}

type HandleCallbackRequest struct {
	ChannelCode     string `json:"channel_code"`
	PlatformOrderNo string `json:"platform_order_no"`
	ChannelOrderNo  string `json:"channel_order_no"`
	ChannelStatus   string `json:"channel_status"`
	Amount          string `json:"amount"`
	RawBody         string `json:"raw_body"`
}

type OrderService struct {
	mu     sync.Mutex
	orders map[string]model.PaymentOrder
	index  map[string]string
}

func NewOrderService() *OrderService {
	return &OrderService{
		orders: make(map[string]model.PaymentOrder),
		index:  make(map[string]string),
	}
}

func (s *OrderService) CreatePaymentOrder(req CreatePaymentOrderRequest) (model.PaymentOrder, error) {
	if req.MerchantNo == "" || req.MerchantOrderNo == "" || req.Amount == "" || req.Currency == "" {
		return model.PaymentOrder{}, errors.New("missing required field")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	indexKey := req.MerchantNo + ":" + req.MerchantOrderNo
	if existingOrderNo, ok := s.index[indexKey]; ok {
		return s.orders[existingOrderNo], nil
	}

	now := time.Now().UTC()
	order := model.PaymentOrder{
		MerchantNo:      req.MerchantNo,
		MerchantOrderNo: req.MerchantOrderNo,
		PlatformOrderNo: "P" + now.Format("20060102150405") + randomHex(4),
		Amount:          req.Amount,
		Currency:        req.Currency,
		Status:          string(statemachine.StatusPaying),
		NotifyStatus:    "PENDING",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	s.orders[order.PlatformOrderNo] = order
	s.index[indexKey] = order.PlatformOrderNo
	return order, nil
}

func (s *OrderService) QueryPaymentOrder(platformOrderNo string) (model.PaymentOrder, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[platformOrderNo]
	return order, ok
}

func (s *OrderService) HandleChannelCallback(req HandleCallbackRequest) (model.PaymentOrder, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[req.PlatformOrderNo]
	if !ok {
		return model.PaymentOrder{}, errors.New("order not found")
	}

	next := channelStatusToPaymentStatus(req.ChannelStatus)
	current := statemachine.PaymentStatus(order.Status)
	if order.Status == string(next) {
		return order, nil
	}
	if !statemachine.CanTransit(current, next) {
		return model.PaymentOrder{}, errors.New("invalid payment status transition")
	}

	order.ChannelCode = req.ChannelCode
	order.ChannelOrderNo = req.ChannelOrderNo
	order.Status = string(next)
	order.UpdatedAt = time.Now().UTC()
	s.orders[order.PlatformOrderNo] = order
	return order, nil
}

func channelStatusToPaymentStatus(status string) statemachine.PaymentStatus {
	switch status {
	case "SUCCESS", "PAID":
		return statemachine.StatusSuccess
	case "FAILED":
		return statemachine.StatusFailed
	default:
		return statemachine.StatusUnknown
	}
}

func randomHex(size int) string {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "0000"
	}
	return hex.EncodeToString(buf)
}

