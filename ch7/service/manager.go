package service

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var manager *Manager

type Manager struct {
	db *gorm.DB
}

func NewManager(db *gorm.DB) *Manager {
	return &Manager{db: db}
}

func GetManager() *Manager {
	return manager
}

func InitManagerFromDsn(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	manager = &Manager{db: db}
	return nil
}
