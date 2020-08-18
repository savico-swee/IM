package models

import (
	"github.com/wuyan94zl/IM/database"
	"time"
)

func init() {
	database.DB.AutoMigrate(&Room{})
	database.DB.AutoMigrate(&RoomHasUser{})
}

type Room struct {
	Id 			int
	Number		string
	Name		string
	UserNum		int
	CreateAt	time.Time
}

type RoomHasUser struct {
	Id 			int
	RoomNumber 	string
	UserId		int
	UserName	string
	CreatedAt 	time.Time
}
