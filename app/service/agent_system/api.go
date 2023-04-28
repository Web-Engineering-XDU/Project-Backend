package agentsystem

import (
	"context"

	"github.com/Web-Engineering-XDU/Project-Backend/app/models"
)

func (as *AgentSystem) AddAgent(ctx context.Context, a models.AddAgentParams, rels []models.AddAgentRelationParams) error {
	tx, err := models.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := models.Query().WithTx(tx)

	if err = qtx.AddAgent(ctx, a); err!=nil {
		return nil
	}

	if rels!=nil {
		for _, v := range rels{
			if err = qtx.AddAgentRelation(ctx, v); err!=nil {
				return nil
			}
		}
	}

	return nil
}

func UpdateAgent() {}

func DryRunAgent() {}

func DeleteAgent() {}

func GetAgentList() {}

func GetAgentDetail() {}

func GetAgentEventList() {}

func GetEventList() {}

func GetPendingEventNumber() {}

func GetRunningAgentList() {}
