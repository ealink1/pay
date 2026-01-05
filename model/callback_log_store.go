package model

import (
	"errors"
	"sync"
	"time"

	"gorm.io/gorm"
)

type CallbackLogStore interface {
	Create(log *CallbackLog) error
}

type NoopCallbackLogStore struct{}

func (s *NoopCallbackLogStore) Create(log *CallbackLog) error {
	return nil
}

type InMemoryCallbackLogStore struct {
	mu   sync.Mutex
	logs []*CallbackLog
}

func (s *InMemoryCallbackLogStore) Create(log *CallbackLog) error {
	if log == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logs = append(s.logs, log)
	return nil
}

type GormCallbackLogStore struct {
	db *gorm.DB
}

func (s *GormCallbackLogStore) Create(log *CallbackLog) error {
	if log == nil {
		return nil
	}
	if log.ReceivedAt.IsZero() {
		log.ReceivedAt = time.Now()
	}
	return s.db.Create(log).Error
}

var CallbackLogs CallbackLogStore = &NoopCallbackLogStore{}

func InitGormCallbackLogStore(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}
	if err := db.AutoMigrate(&CallbackLog{}); err != nil {
		return err
	}
	CallbackLogs = &GormCallbackLogStore{db: db}
	return nil
}
