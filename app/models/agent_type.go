package models

import "gorm.io/plugin/soft_delete"

const agentTypeTableName = "agent_types"

type AgentType struct {
	ID          int `gorm:"primaryKey"`
	Name        string	`gorm:"type:VARCHAR(128);not null"`
	AllowInput  bool	`gorm:"not null"`
	AllowOutput bool 	`gorm:"not null"`
	Deleted     soft_delete.DeletedAt `gorm:"softDelete:flag;type:TINYINT;not null"`
}

func (*AgentType) TableName() string {
	return agentTypeTableName
}

func SelectAgentTypeList() (ret []AgentType){
	DB().Find(&ret)
	return ret
}
