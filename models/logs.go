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
	UserName 	string
	Message 	string
	RoomNumber 	string
	Type 		int
	CreatedAt 	time.Time
	Time 		string	`gorm:"-"`
}

func RightLog(name string,message string,roomNumber string,t int){
	log := Log{UserName: name,Message: message,RoomNumber: roomNumber,Type: t,CreatedAt: time.Now()}
	database.DB.Create(&log)
}