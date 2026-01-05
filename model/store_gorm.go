package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type GormOrderStore struct {
	db *gorm.DB
}

func InitGormOrderStore(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}
	if err := db.AutoMigrate(&Order{}); err != nil {
		return err
	}
	if err := InitGormCallbackLogStore(db); err != nil {
		return err
	}
	Store = &GormOrderStore{db: db}
	return nil
}

func (s *GormOrderStore) Create(order *Order) error {
	if order.ID == "" {
		order.ID = generateID()
	}
	now := time.Now()
	if order.CreatedAt.IsZero() {
		order.CreatedAt = now
	}
	order.UpdatedAt = now
	if order.Status == "" {
		order.Status = OrderStatusPending
	}
	return s.db.Create(order).Error
}

func (s *GormOrderStore) GetByID(id string) (*Order, bool) {
	var order Order
	err := s.db.First(&order, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false
		}
		return nil, false
	}
	return &order, true
}

func (s *GormOrderStore) GetByOutTradeNo(outTradeNo string) (*Order, bool) {
	var order Order
	err := s.db.First(&order, "out_trade_no = ?", outTradeNo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false
		}
		return nil, false
	}
	return &order, true
}

func (s *GormOrderStore) UpdateStatus(id string, status OrderStatus, tradeNo string) error {
	updates := map[string]any{
		"status":     status,
		"updated_at": time.Now(),
	}
	if tradeNo != "" {
		updates["trade_no"] = tradeNo
	}
	return s.db.Model(&Order{}).Where("id = ?", id).Updates(updates).Error
}

func (s *GormOrderStore) List() []*Order {
	var orders []*Order
	_ = s.db.Order("created_at desc").Find(&orders).Error
	return orders
}
