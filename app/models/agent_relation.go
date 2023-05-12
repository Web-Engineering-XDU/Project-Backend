package models

import "gorm.io/gorm"

const agentRelationTableName = "agent_relations"

type AgentRelation struct {
	ID         int `gorm:"primaryKey"  form:"id"          json:"id"`
	SrcAgentId int `gorm:"not null"    form:"srcAgentId"  json:"srcAgentId"`
	DstAgentId int `gorm:"not null"    form:"dstAgentId"  json:"dstAgentId"`
}

func (*AgentRelation) TableName() string {
	return agentRelationTableName
}

func InsertAgentRelation(agentRelations []*AgentRelation) error {
	return DB().Create(agentRelations).Error
}

func DeleteAgentRelation(agentRelation AgentRelation) {
	DB().
		Where("src_agent_id = ? AND dst_agent_id = ?",
			agentRelation.SrcAgentId,
			agentRelation.DstAgentId).
		Delete(&AgentRelation{})
}

func DeleteAllAgentRelationAbout(id int) error {
	return DB().
		Where("src_agent_id = ? OR dst_agent_id = ?", id, id).
		Delete(&AgentRelation{}).Error
}

func SelectAllAgentRelations() (ret []AgentRelation) {
	DB().Find(&ret)
	return ret
}

func SelectAgentRelationsAbout(id int) (ret []AgentRelation) {
	DB().Where("src_agent_id = ? OR dst_agent_id = ?", id, id).Find(&ret)
	return ret
}

func SetAgentRelation(id int, srcs, dsts []int) error {
	return DB().Transaction(func(tx *gorm.DB) error {
		err := DeleteAllAgentRelationAbout(id)
		if err != nil {
			return err
		}
		relations := make([]*AgentRelation, len(srcs)+len(dsts))
		i := 0
		for _, src := range srcs {
			relations[i] = &AgentRelation{SrcAgentId: src, DstAgentId: id}
			i++
		}
		for _, dst := range dsts {
			relations[i] = &AgentRelation{SrcAgentId: id, DstAgentId: dst}
			i++
		}
		err = InsertAgentRelation(relations)
		return err
	})
}
