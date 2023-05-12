package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

const agentTableName = "agents"

type AgentBasic struct {
	ID     int  `gorm:"primaryKey"  form:"id"       json:"id"       example:"8848"`
	Enable bool `gorm:"not null"    form:"enable"   json:"enable"   example:"true"`
	TypeId int  `gorm:"not null"    form:"typeId"   json:"typeId"   example:"1"`
}

type AgentExtra struct {
	Name        string `gorm:"type:VARCHAR(128);not null"  form:"name"         json:"name"          example:"5s timer"`
	Description string `gorm:"type:TEXT;not null"          form:"description"  json:"description"   example:"I'm a schedule agent"`

	CreateAt time.Time `gorm:"not null"  json:"createAt"    example:"2023-04-11T05:07:53+08:00"`
}

type Agent struct {
	AgentBasic
	AgentExtra

	EventForever bool `gorm:"not null"  form:"eventForever"   json:"eventForever"`
	EventMaxAge  int  `gorm:"not null"  form:"eventMaxAge"    json:"eventMaxAge"   example:"3600000000000"`

	PropJsonStr string `gorm:"type:TEXT;not null"  form:"propJsonStr"  json:"propJsonStr"`

	Deleted soft_delete.DeletedAt `gorm:"softDelete:flag;type:TINYINT;not null"`
}

func (*Agent) TableName() string {
	return agentTableName
}

func (u *Agent) ToUpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"enable":        u.Enable,
		"name":          u.Name,
		"description":   u.Description,
		"event_forever": u.EventForever,
		"event_max_age": u.EventMaxAge,
		"prop_json_str": u.PropJsonStr,
	}
}

type AgentRuntime struct {
	AgentBasic

	EventForever bool   `json:"eventForever"    example:"false"`
	EventMaxAge  int    `json:"eventMaxAge"     example:"0"`
	PropJsonStr  string `json:"propJsonStr"     example:"{\"cron\":\"*/5 * * * * *\"}"`

	AllowInput  bool `json:"allowInput"      example:"false"`
	AllowOutput bool `json:"allowOutput"     example:"true"`
}

type AgentDetail struct {
	AgentRuntime

	AgentExtra
	TypeName string `json:"typeName"    example:"Schedule Agent"`
}

func InsertAgent(agent *Agent) error {
	agent.CreateAt = time.Now()
	return DB().Create(agent).Error
}

func DeleteAgent(id int) error {
	return DB().Delete(&Agent{}, id).Error
}

func DeleteAgentAndRelationsAbout(id int) (bool, error) {
	err := DB().Transaction(func(tx *gorm.DB) error {
		err := DeleteAgent(id)
		if err != nil {
			return err
		}
		err = DeleteAllAgentRelationAbout(id)
		if err != nil {
			return err
		}
		return nil
	})

	return err == nil, err
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

func SelectAgentDetailList(limit, offset int) (ret []AgentDetail, totalCount int64) {

	tx := DB().
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
		Joins("INNER JOIN agent_types ON agents.type_id = agent_types.id").Count(&totalCount)

	if offset+limit < 1 {
		ret = []AgentDetail{}
	} else {
		tx.Limit(limit).Offset(offset).
			Scan(&ret)
	}

	return ret, totalCount
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

type AgentIdAndName struct {
	Id   int
	Name string
}

// func (*AgentIdAndName) TableName() string {
// 	return agentTableName
// }

func SelectAgentIdAndNames(targets []int) (ret []AgentIdAndName){
	DB().Model(&Agent{}).
		Select("id, name").
		Find(&ret, targets)
	return ret
}

func SelectPrevableAgents(id int, key string, limit, offset int) (ret []AgentIdAndName, totalCount int64, err error) {
	tx := DB().Model(&Agent{}).
		Where("agents.id <> ? && allow_output = true && agents.name LIKE ?", id, "%"+key+"%").
		Select("agents.id, agents.name").
		Joins("INNER JOIN agent_types ON agents.type_id = agent_types.id").
		Count(&totalCount)
	if offset+limit < 1 || totalCount == 0 {
		ret = []AgentIdAndName{}
	} else {
		tx.Limit(limit).Offset(offset).Find(&ret)
	}

	return ret, totalCount, nil
}

func SelectNextableAgents(id int, key string, limit, offset int) (ret []AgentIdAndName, totalCount int64, err error) {
	tx := DB().Model(&Agent{}).
		Where("agents.id <> ? && allow_input = true && agents.name LIKE ?", id, "%"+key+"%").
		Select("agents.id, agents.name").
		Joins("INNER JOIN agent_types ON agents.type_id = agent_types.id").
		Count(&totalCount)
	if offset+limit < 1 || totalCount == 0 {
		ret = []AgentIdAndName{}
	} else {
		tx.Limit(limit).Offset(offset).Find(&ret)
	}

	return ret, totalCount, nil
}

func UpdateAgent(agent *Agent) error {
	return DB().Model(agent).Updates(agent.ToUpdateMap()).Error
}
