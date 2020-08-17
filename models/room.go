package models

import (
	"github.com/wuyan94zl/IM/database"
	"time"
)

func init() {
	database.DB.AutoMigrate(&Room{})
}

type Room struct {
	Id 			int
	Number		string
	Name		string
	CreateAt	time.Time
}
