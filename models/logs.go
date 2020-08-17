package models

import (
	"github.com/wuyan94zl/IM/database"
	"time"
)

func init() {
	database.DB.AutoMigrate(&Log{})
}

type Log struct {
	Id 			int
	UserId 		int
	Message 	string
	CreatedAt 	time.Time
}