package model

import "time"

type CallbackLog struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	Provider    string    `gorm:"type:varchar(32);index"`
	Path        string    `gorm:"type:varchar(255)"`
	Method      string    `gorm:"type:varchar(16)"`
	RemoteIP    string    `gorm:"type:varchar(64)"`
	TraceID     string    `gorm:"type:varchar(64);index"`
	AppID       string    `gorm:"type:varchar(64);index"`
	OutTradeNo  string    `gorm:"type:varchar(64);index"`
	TradeNo     string    `gorm:"type:varchar(64);index"`
	TradeStatus string    `gorm:"type:varchar(64);index"`
	NotifyID    string    `gorm:"type:varchar(128);index"`
	Sign        string    `gorm:"type:text"`
	VerifyOK    bool      `gorm:"index"`
	VerifyError string    `gorm:"type:text"`
	ParamsJSON  string    `gorm:"type:longtext"`
	HeadersJSON string    `gorm:"type:longtext"`
	ReceivedAt  time.Time `gorm:"index"`
	CreatedAt   time.Time `gorm:"index"`
}

func (CallbackLog) TableName() string {
	return "callback_log"
}
