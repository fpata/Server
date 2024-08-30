package Patient

import (
	"clinic_server/database"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var patients = []Patient{}

func GetPatientById(c *gin.Context) {
	patientId := c.Param("ID")
	var objPatient = Patient{}
	var db *gorm.DB = database.GetDBContext()
	error := db.Model(&objPatient).First(&objPatient, patientId).Error
	if error != nil {
		fmt.Println(error)
	}
	c.IndentedJSON(http.StatusOK, objPatient)
}

func CreatePatient(c *gin.Context) {
	var objPatient Patient
	c.ShouldBind(&objPatient)
	var db *gorm.DB = database.GetDBContext()
	createResult := db.Create(&objPatient)

	if createResult.Error != nil {
		fmt.Println(createResult.Error)
	} else {
		fmt.Println(createResult.RowsAffected)
	}
	c.IndentedJSON(http.StatusOK, objPatient)
}

func UpdatePatient(c *gin.Context) {
	var objPatient Patient
	c.ShouldBind(&objPatient)
	var db *gorm.DB = database.GetDBContext()
	updateResult := db.Save(&objPatient)
	if updateResult.Error != nil {
		fmt.Println(updateResult.Error)
	} else {
		fmt.Println(updateResult.RowsAffected)
	}
	c.IndentedJSON(http.StatusOK, true)
}

func PatchPatient(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
	}

	var db *gorm.DB = database.GetDBContext()
	updateResult := db.Updates(&jsonData)
	if updateResult.Error != nil {
		fmt.Println(updateResult.Error)
	} else {
		fmt.Println(updateResult.RowsAffected)
	}
	c.IndentedJSON(http.StatusOK, patients)
}

func DeletePatient(c *gin.Context) {
	patientId := c.Query("ID")
	var db *gorm.DB = database.GetDBContext()
	db.Delete(&Patient{}, patientId)
	c.IndentedJSON(http.StatusOK, "")
}
