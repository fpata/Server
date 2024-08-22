package Patient

import (
	"clinic_server/database"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var login = Login{}

func ValidateLogin(c *gin.Context) {
	err := c.ShouldBindJSON(&login)
	if err != nil {
		fmt.Println(err)
		c.Error(err)
		c.Abort()
		return
	}
	var Id int64 = 0
	var db *gorm.DB = database.GetDBContext()
	db.Table("Login").Where("Username = ? AND Password = ?", login.UserName, login.Password).Select("Id").Find(&Id)
	c.IndentedJSON(200, gin.H{"Id": strconv.Itoa((int)(Id))})
}
