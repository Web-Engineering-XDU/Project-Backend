package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Web-Engineering-XDU/Project-Backend/config"
	_ "github.com/go-sql-driver/mysql"
)

var quries *Queries
var db *sql.DB

func InitDB(config config.MySQL) (*Queries, error) {
	if quries == nil {
		var err error
		db, err = sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/huggo?parseTime=true", config.User, config.Password, config.Host, config.Port))
		if err != nil {
			return nil, err
		}
		quries = New(db)
		return quries, nil
	}
	return nil, errors.New("Cannot init DB twice")
}

func Query() *Queries {
	return quries
}

func DB() *sql.DB {
	return db
}
