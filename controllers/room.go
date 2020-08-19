package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/wuyan94zl/IM/database"
	"github.com/wuyan94zl/IM/models"
	"html/template"
	rand2 "math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func RoomList(c *gin.Context)  {
	rooms := []models.Room{}
	database.DB.Find(&rooms)
	c.HTML(http.StatusOK, "rooms.html", gin.H{
		"rooms": rooms,
	})
}

func RoomAdd(c *gin.Context) {
	rand2.Seed(time.Now().UnixNano())
	n := rand2.Intn(999999)
	number := 888000 + n
	room := models.Room{Number: strconv.Itoa(number),Name: c.PostForm("name"),UserNum: 0,CreateAt: time.Now()}
	database.DB.Create(&room)
	c.Redirect(http.StatusMovedPermanently,"/rooms")
}

func RoomInfo(c *gin.Context) {
	number := c.Param("number")
	name := c.Param("name")
	if number == "" || name == "" {
		return
	}
	user := models.User{}
	database.DB.Where("name = ?",name).First(&user)
	if user.Name == name {
		rHU := models.RoomHasUser{}
		database.DB.Where("user_name = ? and room_number = ?",name,number).First(&rHU)
		if rHU.UserName == name {
			//c.Redirect(http.StatusMovedPermanently,"/rooms")
			return
		}
	}else{
		user := models.User{Name: name,Email: "ds",Password: "ddd",CreateAt: time.Now()}
		database.DB.Create(&user)
	}
	room := models.Room{}
	database.DB.Where("number = ?",number).First(&room)
	roomHasUser := []models.RoomHasUser{}
	database.DB.Where("room_number = ?",number).Find(&roomHasUser)
	roomHasUser = append(roomHasUser,models.RoomHasUser{Id: 999999999,UserName: name})
	if room.Number != number {
		return
	}

	logs := []models.Log{}
	database.DB.Offset(0).Limit(30).Order("id desc").Find(&logs)

	for k,v := range logs{
		logs[k].Time = v.CreatedAt.Format("2006-01-02 15:04:05")
		logs[k].HtmlMessage = template.HTML(v.Message)
	}
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].Id < logs[j].Id
	})
	c.HTML(http.StatusOK, "index.html", gin.H{
		"room": room,
		"name": name,
		"logs": logs,
		"userNum": len(roomHasUser),
		"roomHasUser": roomHasUser,
	})
}
