package model

import "time"

type OrderStatus string

const (
	OrderStatusPending OrderStatus = "pending"
	OrderStatusPaid    OrderStatus = "paid"
	OrderStatusFailed  OrderStatus = "failed"
	OrderStatusClosed  OrderStatus = "closed"
)

type Order struct {
	ID          string      `json:"id" gorm:"primaryKey;type:varchar(64)"`
	OutTradeNo  string      `json:"out_trade_no" gorm:"uniqueIndex;type:varchar(64)"`
	TotalAmount string      `json:"total_amount" gorm:"type:varchar(32)"`
	Subject     string      `json:"subject" gorm:"type:varchar(255)"`
	Body        string      `json:"body" gorm:"type:text"`
	QrCode      string      `json:"qr_code" gorm:"type:longtext"`
	Status      OrderStatus `json:"status" gorm:"type:varchar(16);index"`
	TradeNo     string      `json:"trade_no,omitempty" gorm:"type:varchar(64)"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
