package swaggo

import "github.com/Web-Engineering-XDU/Project-Backend/app/models"

type GetAgentListResponseResult struct {
	Content    []models.AgentDetail `json:"content"`
	Count      int                  `json:"count"`
	TotalCount int                  `json:"totalCount"`
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
	Content    []models.Event `json:"content"`
	Count      int            `json:"count"`
	TotalCount int            `json:"totalCount"`
}

type GetEventListResponse struct {
	StateInfo
	Result GetEventListResponseResult `json:"result"`
}

type GetRelationsResponse struct {
	StateInfo
	Result []models.AgentRelation `json:"result"`
}

type DryRunResponse struct {
	StateInfo
	Result []map[string]string `json:"result"`
}

type GetRelationableAgentsResp struct {
	StateInfo
	Result GetRelationableAgentsRespResult `json:"result"`
}

type GetRelationableAgentsRespResult struct {
	Content    []models.AgentIdAndName `json:"content"`
	Count      int                     `json:"count"`
	TotalCount int                     `json:"totalCount"`
}

type GetRelationsForEditResp struct {
	StateInfo
	Result struct {
		Src []models.AgentIdAndName `json:"srcs"`
		Dst []models.AgentIdAndName `json:"dsts"`
	} `json:"result"`
}
