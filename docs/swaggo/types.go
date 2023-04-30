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
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type GetAgentListResponse struct {
	StateInfo
	Result GetAgentListResponseResult `json:"result"`
}

type NewAgentResponse struct {
	StateInfo
	Result NewAgentResponseResult `json:"result"`
}

type DeleteAgentResponse struct {
	StateInfo
}
