package model

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusFailed    OrderStatus = "failed"
	OrderStatusClosed    OrderStatus = "closed"
)

type Order struct {
	ID          string      `json:"id"`
	OutTradeNo  string      `json:"out_trade_no"`
	TotalAmount string      `json:"total_amount"`
	Subject     string      `json:"subject"`
	Body        string      `json:"body"`
	QrCode      string      `json:"qr_code"`
	Status      OrderStatus `json:"status"`
	TradeNo     string      `json:"trade_no,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
