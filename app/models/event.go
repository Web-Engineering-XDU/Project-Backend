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

func InsertEvent(event *Event) error {
	return DB().Create(event).Error
}

func SelectHashCount(hash string, id int) (ret int64) {
	DB().Where("content_hash = ? AND src_agent_id = ?", hash, id).Count(&ret)
	return ret
}

func SelectEventList(limit, offset int) (ret []Event) {
	if offset+limit < 1 {
		return []Event{}
	}
	DB().Model(&Event{}).
		Order("create_at desc").
		Limit(limit).Offset(offset).Find(&ret)
	return ret
}

func SelectEventListByAgentID(id, limit, offset int) (ret []Event) {
	if offset+limit < 1 {
		return []Event{}
	}
	DB().Model(&Event{}).
		Where("src_agent_id = ?", id).
		Order("create_at desc").
		Limit(limit).Offset(offset).Find(&ret)
	return ret
}
