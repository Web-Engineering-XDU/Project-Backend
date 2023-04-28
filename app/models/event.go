package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Event struct {
	ID          int `gorm:"primaryKey"`
	SrcAgentId  int
	JsonStr     string
	ContentHash string
	Error       bool
	Log         string
	CreateAt    time.Time
	DeleteAt    time.Time
	Deleted     soft_delete.DeletedAt `gorm:"softDelete:flag"`
}

func InsertEvent(event *Event) error{
	return DB().Create(event).Error
}

