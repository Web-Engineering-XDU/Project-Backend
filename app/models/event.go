package models

import (
    "time"

    "gorm.io/plugin/soft_delete"
)

type Event struct {
    ID          int		`gorm:"primaryKey"              form:"id"           json:"id"`
    SrcAgentId  int		`gorm:"not null"                form:"srcAgentId"   json:"srcAgentId"`
    JsonStr     string	`gorm:"type:TEXT;not null"      form:"jsonStr"      json:"jsonStr"`
    ContentHash string	`gorm:"type:CHAR(16);not null"  form:"contentHash"  json:"contentHash"`
    Error       bool	`gorm:"not null"                form:"error"        json:"error"`
    Log         string	`gorm:"type:TEXT;not null"      form:"log"          json:"log"`

    CreateAt    time.Time	`gorm:"not null"    json:"createAt"`
    DeleteAt    time.Time	`gorm:"not null"    json:"deleteAt"`

    Deleted     soft_delete.DeletedAt `gorm:"softDelete:flag;type:TINYINT;not null"`
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
