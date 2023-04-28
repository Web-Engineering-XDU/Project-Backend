package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

const agentTableName = "agents"

type AgentBasic struct {
	ID     int `gorm:"primaryKey"`
	Enable bool
	TypeId int
}

type AgentExtra struct {
	Name        string
	Description string
	CreateAt    time.Time
}

type Agent struct {
	AgentBasic
	AgentExtra

	EventForever bool
	EventMaxAge  time.Duration
	PropJsonStr  string

	Deleted soft_delete.DeletedAt `gorm:"softDelete:flag"`
}

func (*Agent) TableName() string {
	return agentTableName
}

type AgentRuntime struct {
	AgentBasic

	EventForever bool
	EventMaxAge  time.Duration
	PropJsonStr  string

	AllowInput  bool `gorm:"column:allow_input"`
	AllowOutput bool `gorm:"column:allow_output"`
}

func InsertAgent(agent *Agent) error {
	return DB().Create(agent).Error
}

func DeleteAgent(agent *Agent) {
	DB().Delete(agent)
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

func SelectAgentDetailByID(id int) (ret Agent, ok bool) {
	if DB().Omit("deleted").First(&ret, "id = ?", id).RowsAffected == 1{
		return ret, true
	} else {
		return Agent{}, false
	}
}

func SelectAgentList() (ret []Agent) {
	DB().Find(&ret)
	return ret
}