package models

import "gorm.io/plugin/soft_delete"

const agentTypeTableName = "agent_types"

type AgentType struct {
	ID          int     `gorm:"primaryKey"                  form:"id"`
	Name        string	`gorm:"type:VARCHAR(128);not null"  form:"name"`
	AllowInput  bool	`gorm:"not null"                    form:"allowInput"`
	AllowOutput bool 	`gorm:"not null"                    form:"allowOutput"`
    
	Deleted     soft_delete.DeletedAt `gorm:"softDelete:flag;type:TINYINT;not null"`
}

func (*AgentType) TableName() string {
	return agentTypeTableName
}

func SelectAgentTypeList() (ret []AgentType){
	DB().Find(&ret)
	return ret
}
