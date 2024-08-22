package Patient

import (
	"clinic_server/database"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var patients = []Patient{}

func GetAllPatientsWithPaging(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	var db *gorm.DB = database.GetDBContext()
	db.Table("Patient").Limit(limit).Offset(offset).Find(&patients)
	c.IndentedJSON(http.StatusOK, patients)
}

func GetAllPatients(c *gin.Context) {
	var db *gorm.DB = database.GetDBContext()
	db.Table("Patient").Find(&patients)
	c.IndentedJSON(http.StatusOK, patients)
}

func GetPatientById(c *gin.Context) {
	patientId := c.Param("Id")
	var objPatient = Patient{}
	var db *gorm.DB = database.GetDBContext()
	error := db.Model(&objPatient).Preload("PatientTreatments").Preload("PatientTreatments.PatientTreatmentDetails").
		Preload("PatientReports").Preload("PatientAppointments").
		First(&objPatient, patientId).Error
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
	c.IndentedJSON(http.StatusOK, objPatient)
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
	patientId := c.Query("Id")
	var db *gorm.DB = database.GetDBContext()
	db.Delete(&Patient{}, patientId)
	c.IndentedJSON(http.StatusOK, "")
}
