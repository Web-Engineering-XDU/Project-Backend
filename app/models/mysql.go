package models

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

var quries *Queries

func InitDB(dataSourceName string) (*Queries, error) {
	if quries == nil {
		db, err := sql.Open("mysql", "root:z0013@/huggo?parseTime=true")
		if err != nil {
			return nil, err
		}
		quries = New(db)
		return quries, nil
	}
	return nil, errors.New("Cannot init DB twice")
}

func DB() *Queries {
	return quries
}


