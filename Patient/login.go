package Patient

import (
	"clinic_server/database"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var login = LoginModel{}

func ValidateLogin(c *gin.Context) {
	err := c.ShouldBindJSON(&login)
	if err != nil {
		fmt.Println(err)
		c.Error(err)
		c.Abort()
		return
	}
	var db *gorm.DB = database.GetDBContext()
	db.Table("Login").Where("Username = ? AND Password = ?", login.UserName, login.Password).Scan(&login)
	c.IndentedJSON(200, login)
}
