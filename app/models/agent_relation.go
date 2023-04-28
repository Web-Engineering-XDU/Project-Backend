package models

const agentRelationTableName = "agent_relations"

type AgentRelation struct {
	ID         int `gorm:"primaryKey"`
	SrcAgentId int
	DstAgentId int
}

func (*AgentRelation) TableName() string {
	return agentRelationTableName
}

func InsertAgentRelation(agentRelation *AgentRelation) error {
	return DB().Create(agentRelation).Error
}

func DeleteAgentRelation(agentRelation AgentRelation) {
	DB().
		Where("src_agent_id = ? AND dst_agent_id = ?",
			agentRelation.SrcAgentId,
			agentRelation.DstAgentId).
		Delete(&AgentRelation{})
}

func DeleteAllAgentRelationAbout(id int) {
	DB().
		Where("src_agent_id = ? OR dst_agent_id = ?", id, id).
		Delete(&AgentRelation{})
}

func SelectAgentRelationList() (ret []AgentRelation) {
	DB().Find(&ret)
	return ret
}