// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package models

import (
	"time"
)

type Agent struct {
	ID          int32
	Name        string
	Enable      bool
	TypeID      int32
	EventMaxAge int64
	PropJsonStr string
	CreateAt    time.Time
	Deleted     bool
	Description string
}

type AgentRelation struct {
	ID         int32
	SrcAgentID int32
	DstAgentID int32
}

type AgentType struct {
	ID          int32
	Name        string
	AllowInput  bool
	AllowOutput bool
	Deleted     bool
}

type Event struct {
	ID          int32
	SrcAgentID  int32
	JsonStr     string
	ContentHash string
	Error       bool
	Log         string
	CreateAt    time.Time
	DeleteAt    time.Time
	Deleted     bool
}
