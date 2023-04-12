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
	Enable      int32
	TypeID      int32
	EventMaxAge int64
	PropJsonStr string
	CreateAt    time.Time
	Deleted     int32
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
	AllowInput  int32
	AllowOutput int32
	Deleted     int32
}

type Event struct {
	ID         int32
	SrcAgentID int32
	JsonStr    string
	Error      int32
	Log        string
	CreateAt   time.Time
	DeleteAt   time.Time
	Deleted    int32
}
