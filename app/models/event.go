package models

import (
    "time"

    "gorm.io/plugin/soft_delete"
)

type Event struct {
    ID          int		`gorm:"primaryKey"              form:"id"           json:"id"           example:"15"`
    SrcAgentId  int		`gorm:"not null"                form:"srcAgentId"   json:"srcAgentId"   example:"1145"`
    JsonStr     string	`gorm:"type:TEXT;not null"      form:"jsonStr"      json:"jsonStr"      example:"{\"ip_without_dot\":\"11125117081\",\"ip\":\"111.251.170.81\",\"latitude\":\"25.0504\"}"`     
    ContentHash string	`gorm:"type:CHAR(16);not null"  form:"contentHash"  json:"contentHash"  example:"47ed5d26b86a7519"`
    Error       bool	`gorm:"not null"                form:"error"        json:"error"        example:"false"`
    Log         string	`gorm:"type:TEXT;not null"      form:"log"          json:"log"          example:""`

    CreateAt    time.Time	`gorm:"not null"    json:"createAt"     example:"2023-04-11T05:07:53+08:00"`
    DeleteAt    time.Time	`gorm:"not null"    json:"deleteAt"     example:"2023-04-11T06:07:53+08:00"`

    Deleted     soft_delete.DeletedAt `gorm:"softDelete:flag;type:TINYINT;not null" swaggerignore:"true"`
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
