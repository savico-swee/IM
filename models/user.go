package models

import (
	"github.com/wuyan94zl/IM/database"
	"time"
)

func init() {
	database.DB.AutoMigrate(&User{})
}

type User struct {
	Id 			int
	Email		string
	Password	string
	Name 		string
	CreateAt	time.Time
}
