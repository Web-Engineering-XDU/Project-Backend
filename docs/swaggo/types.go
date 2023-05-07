package swaggo

import "github.com/Web-Engineering-XDU/Project-Backend/app/models"

type GetAgentListResponseResult struct {
	Content []models.AgentDetail `json:"content"`
	Count   int                  `json:"count"`
}

type NewAgentResponseResult struct {
	Id int `json:"id"`
}

type StateInfo struct {
	Code int    `json:"code"    example:"200"`
	Msg  string `json:"msg"     example:"ok"`
}

type GetAgentListResponse struct {
	StateInfo
	Result GetAgentListResponseResult `json:"result"`
}

type NewAgentResponse struct {
	StateInfo
	Result NewAgentResponseResult `json:"result"`
}

type GetEventListResponseResult struct {
	Content []models.Event `json:"content"`
	Count   int            `json:"count"`
}

type GetEventListResponse struct {
	StateInfo
	Result GetEventListResponseResult `json:"result"`
}

type GetRelationsResponse struct {
	StateInfo
	Result []models.AgentRelation `json:"result"`
}
