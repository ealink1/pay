package model

import (
	"sync"
	"time"
)

type OrderStore interface {
	Create(order *Order) error
	GetByID(id string) (*Order, bool)
	GetByOutTradeNo(outTradeNo string) (*Order, bool)
	UpdateStatus(id string, status OrderStatus, tradeNo string) error
	List() []*Order
}

type InMemoryOrderStore struct {
	mu     sync.RWMutex
	orders map[string]*Order
}

var Store OrderStore = &InMemoryOrderStore{orders: make(map[string]*Order)}

func (s *InMemoryOrderStore) Create(order *Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order.ID = generateID()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.Status = OrderStatusPending

	s.orders[order.ID] = order
	return nil
}

func (s *InMemoryOrderStore) GetByID(id string) (*Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, exists := s.orders[id]
	return order, exists
}

func (s *InMemoryOrderStore) GetByOutTradeNo(outTradeNo string) (*Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, order := range s.orders {
		if order.OutTradeNo == outTradeNo {
			return order, true
		}
	}
	return nil, false
}

func (s *InMemoryOrderStore) UpdateStatus(id string, status OrderStatus, tradeNo string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[id]
	if !exists {
		return nil
	}

	order.Status = status
	order.UpdatedAt = time.Now()
	if tradeNo != "" {
		order.TradeNo = tradeNo
	}

	return nil
}

func (s *InMemoryOrderStore) List() []*Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]*Order, 0, len(s.orders))
	for _, order := range s.orders {
		orders = append(orders, order)
	}
	return orders
}

func generateID() string {
	return time.Now().Format("20060102150405") + randomString(6)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
