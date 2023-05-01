package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

const agentTableName = "agents"

type AgentBasic struct {
	ID     int  `gorm:"primaryKey"  form:"id"`
	Enable bool `gorm:"not null"    form:"enable"`
	TypeId int  `gorm:"not null"    form:"typeId"`
}

type AgentExtra struct {
	Name        string    `gorm:"type:VARCHAR(128);not null"  form:"name"`
	Description string    `gorm:"type:TEXT;not null"          form:"description"`
	CreateAt    time.Time `gorm:"not null"`
}

type Agent struct {
	AgentBasic
	AgentExtra

	EventForever bool           `gorm:"not null"  form:"eventForever"`
	EventMaxAge  time.Duration  `gorm:"not null"  form:"eventMaxAge"`

	PropJsonStr  string         `gorm:"type:TEXT;not null"`

	Deleted soft_delete.DeletedAt `gorm:"softDelete:flag;type:TINYINT;not null"`
}

func (*Agent) TableName() string {
	return agentTableName
}

func (u *Agent) ToUpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"enable":       u.Enable,
		"name":         u.Name,
		"description":  u.Description,
		"eventForever": u.EventForever,
		"eventMaxAge":  u.EventMaxAge,
		"propJsonStr":  u.PropJsonStr,
	}
}

type AgentRuntime struct {
	AgentBasic

	EventForever bool
	EventMaxAge  time.Duration
	PropJsonStr  string

	AllowInput  bool
	AllowOutput bool
}

type AgentDetail struct {
	AgentRuntime

	AgentExtra
	TypeName string
}

func InsertAgent(agent *Agent) error {
	agent.CreateAt = time.Now()
	return DB().Create(agent).Error
}

func DeleteAgent(id int) bool {
	return DB().Delete(&Agent{}, id).RowsAffected > 0
}

func SelectAgentRuntimeList() (ret []AgentRuntime) {
	DB().
		Model(&Agent{}).
		Select(`agents.id id,
		agents.enable enable,
		agents.event_forever event_forever,
		agents.event_max_age event_max_age,
		agents.prop_json_str prop_json_str,
		agents.type_id type_id,
		agent_types.allow_input allow_input,
		agent_types.allow_output allow_output`).
		Joins("INNER JOIN agent_types ON agents.type_id = agent_types.id").
		Scan(&ret)
	return ret
}

func SelectAgentBasicInfoList(limit, offset int) (ret []Agent) {
	DB().Select("id, name, enable, type_id").
		Limit(limit).Offset(offset).
		Find(&ret)
	return ret
}

func SelectAgentDetailList(limit, offset int) (ret []AgentDetail) {
	if offset+limit < 1 {
		return []AgentDetail{}
	}

	DB().
		Model(&Agent{}).
		Select(`agents.id id,
		agents.name name,
		agents.enable enable,
		agents.event_forever event_forever,
		agents.event_max_age event_max_age,
		agents.prop_json_str prop_json_str,
		agents.create_at create_at,
		agents.type_id type_id,
		agent_types.name type_name,
		agent_types.allow_input allow_input,
		agent_types.allow_output allow_output`).
		Joins("INNER JOIN agent_types ON agents.type_id = agent_types.id").
		Limit(limit).Offset(offset).
		Scan(&ret)
	return ret
}

func SelectAgentDetailByID(id int) (ret AgentDetail, ok bool) {
	if DB().
		Model(&Agent{}).
		Where("agents.id = ?", id).
		Select(`agents.id id,
		agents.name name,
		agents.enable enable,
		agents.event_forever event_forever,
		agents.event_max_age event_max_age,
		agents.prop_json_str prop_json_str,
		agents.create_at create_at,
		agents.type_id type_id,
		agent_types.name type_name,
		agent_types.allow_input allow_input,
		agent_types.allow_output allow_output`).
		Joins("INNER JOIN agent_types ON agents.type_id = agent_types.id").
		First(&ret).RowsAffected == 1 {

		return ret, true
	} else {
		return ret, false
	}
}

func UpdateAgent(agent *Agent) bool {
	return DB().Model(agent).Updates(agent.ToUpdateMap()).RowsAffected > 0
}
