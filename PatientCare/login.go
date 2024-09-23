package PatientCare

import (
	"clinic_server/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var login = LoginModel{}

func ValidateLogin(c *gin.Context) {
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var db *gorm.DB = database.GetDBContext()
	err = db.Table("Login").Where("Username = ? AND Password = ?", login.UserName, login.Password).Scan(&login).Error

	if err == nil {
		c.IndentedJSON(http.StatusOK, login)
	} else {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}
