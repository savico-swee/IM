package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wuyan94zl/IM/database"
	"github.com/wuyan94zl/IM/models"
	"net/http"
	"time"
)

func RoomList(c *gin.Context)  {
	rooms := []models.Room{}
	database.DB.Find(&rooms)
	fmt.Println(rooms)
	c.HTML(http.StatusOK, "rooms.html", gin.H{
		"rooms": rooms,
	})
}

func RoomAdd(c *gin.Context) {
	room := models.Room{Number: "ddsasd",Name: c.DefaultQuery("name","w"),CreateAt: time.Now()}
	fmt.Println(room)
	database.DB.Create(&room)
	c.Redirect(http.StatusMovedPermanently,"/rooms")
}

func RoomDel(c *gin.Context) {

}
