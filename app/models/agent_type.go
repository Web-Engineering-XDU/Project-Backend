package models

import "gorm.io/plugin/soft_delete"

const agentTypeTableName = "agent_types"

type AgentType struct {
	ID          int `gorm:"primaryKey"`
	Name        string
	AllowInput  bool
	AllowOutput bool
	Deleted     soft_delete.DeletedAt `gorm:"softDelete:flag"`
}

func (*AgentType) TableName() string {
	return agentTypeTableName
}

func SelectAgentTypeList() (ret []AgentType){
	DB().Find(&ret)
	return ret
}
