package database
import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
var DB *gorm.DB // 定义 mysql 连接实例
var err_db error
// 初始化 mysql DB 连接实例
func init() {
	// 单例模式获取数据库连接 实例
	DB, err_db = gorm.Open("mysql", "root:123456@/imdatabase?charset=utf8&parseTime=True&loc=Local")
	if err_db != nil {
		panic(err_db)
	}
	//DB = DB.Debug() // 全局显示sql详情
	// defer DB.close() // 持久连接 就不需要关闭了
}
