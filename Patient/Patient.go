package Patient

import (
	"bytes"
	"clinic_server/database"
	"clinic_server/types"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Patient struct {
	Id                       types.NullInt64  `gorm:"PrimaryKey"`
	FirstName                types.NullString `json:"FirstName"`
	LastName                 types.NullString `json:"LastName"`
	MiddleName               types.NullString `json:"MiddleName"`
	Age                      types.NullInt64  `json:"Age"`
	Gender                   types.NullInt64  `json:"Gender"`
	UserName                 types.NullString `json:"UserName"`
	Password                 types.NullString `json:"Password"`
	PermAddress1             types.NullString `json:"PermAddress1"`
	PermAddress2             types.NullString `json:"PermAddress2"`
	PermCity                 types.NullString `json:"PermCity"`
	PermState                types.NullString `json:"PermState"`
	PermCountry              types.NullString `json:"PermCountry"`
	PermPostalCode           types.NullString `json:"PermPostalCode"`
	CorrAddress1             types.NullString `json:"CorrAddress1"`
	CorrAddress2             types.NullString `json:"CorrAddress2"`
	CorrCity                 types.NullString `json:"CorrCity"`
	CorrState                types.NullString `json:"CorrState"`
	CorrCountry              types.NullString `json:"CorrCountry"`
	CorrPostalCode           types.NullString `json:"CorrPostalCode"`
	PrimaryPhone             types.NullString `json:"PrimaryPhone"`
	PrimaryEmail             types.NullString `json:"PrimaryEmail"`
	SecondaryPhone           types.NullString `json:"SecondaryPhone"`
	SecondaryEmail           types.NullString `json:"SecondaryEmail"`
	EmergencyContactName     types.NullString `json:"EmergencyContactName"`
	EmergencyContactEmail    types.NullString `json:"EmergencyContactEmail"`
	EmergencyContactPhone    types.NullString `json:"EmergencyContactPhone"`
	EmergencyContactRelation types.NullString `json:"EmergencyContactRelation"`
	ExistingDieases          types.NullString `json:"ExistingDieases"`
	Medications              types.NullString `json:"Medications"`
	Allergies                types.NullString `json:"Allergies"`
	FatherMedicalHistory     types.NullString `json:"FatherMedicalHistory"`
	MotherMedicalHistory     types.NullString `json:"MotherMedicalHistory"`
	PatientReports           []*PatientReport
	PatientTreatments        []*PatientTreatment
}

type PatientReport struct {
	Id             types.NullInt64  `gorm:"PrimaryKey"`
	PatientID      types.NullInt64  `json:"PatientId"`
	ReportDate     types.NullInt64  `json:"ReportDate"`
	ReportName     types.NullString `json:"ReportName"`
	RepoprtFinding types.NullString `json:"RepoprtFinding"`
	DoctorName     types.NullString `json:"DoctorName"`
}

type PatientTreatment struct {
	Id                      types.NullInt64  `gorm:"PrimaryKey"`
	PatientID               types.NullInt64  `json:"PatientId"`
	ChiefComplaint          types.NullString `json:"ChiefComplaint"`
	Observation             types.NullString `json:"Observation"`
	TreatmentPlan           types.NullString `json:"TreatmentPlan"`
	PatientTreatmentDetails []*PatientTreatmentDetail
}

type PatientTreatmentDetail struct {
	Id                 types.NullInt64  `gorm:"PrimaryKey"`
	PatientID          types.NullInt64  `json:"PatientId"`
	PatientTreatmentID types.NullInt64  `json:"PatientTreatmentID"`
	Tooth              types.NullString `json:"Tooth"`
	Procedure          types.NullString `json:"Procedure"`
	Advice             types.NullString `json:"Advice"`
}

type SearchResult struct {
	Id           types.NullInt64  `gorm:"PrimaryKey"`
	FirstName    types.NullString `json:"FirstName"`
	LastName     types.NullString `json:"LastName"`
	PrimaryPhone types.NullString `json:"PrimaryPhone"`
	PrimaryEmail types.NullString `json:"PrimaryEmail"`
	PermCity     types.NullString `json:"PermCity"`
}

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

func GetPatientByParams(c *gin.Context) {
	var searchCondition SearchResult
	var searchResult []SearchResult
	if err := c.ShouldBindJSON(&searchCondition); err != nil {
		fmt.Println(err)
		c.Error(err)
		c.Abort()
		return
	}
	var query = getWhereClausenBasedOnSearch(searchCondition)
	var db *gorm.DB = database.GetDBContext()
	db.Raw(query).Scan(&searchResult)
	c.IndentedJSON(http.StatusOK, searchResult)
}

func GetPatientById(c *gin.Context) {
	patientId := c.Param("Id")
	var objPatient = Patient{}
	var db *gorm.DB = database.GetDBContext()
	error := db.Model(&objPatient).Preload("PatientTreatments").Preload("PatientTreatments.PatientTreatmentDetails").Preload("PatientReports").First(&objPatient, patientId).Error
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

func getWhereClausenBasedOnSearch(searchCondition SearchResult) string {
	var putAndCondition bool = false
	var sqlQuery bytes.Buffer
	sqlQuery.WriteString("Select Id,FirstName,LastName,PrimaryPhone,PrimaryEmail,PermCity from user Where ")
	if len(searchCondition.FirstName.String) != 0 {
		sqlQuery.WriteString("FirstName like '%")
		sqlQuery.WriteString(searchCondition.FirstName.String)
		sqlQuery.WriteString("%'")
		putAndCondition = true
	}
	if len(searchCondition.LastName.String) != 0 {
		if putAndCondition {
			sqlQuery.WriteString(" And ")
		}
		sqlQuery.WriteString(" LastName like '%")
		sqlQuery.WriteString(searchCondition.LastName.String)
		sqlQuery.WriteString("%'")
		putAndCondition = true
	}
	if len(searchCondition.PrimaryEmail.String) != 0 {
		if putAndCondition {
			sqlQuery.WriteString(" And ")
		}
		sqlQuery.WriteString(" PrimaryEmail like '%")
		sqlQuery.WriteString(searchCondition.PrimaryEmail.String)
		sqlQuery.WriteString("%'")
		putAndCondition = true
	}
	if len(searchCondition.PrimaryPhone.String) != 0 {
		if putAndCondition {
			sqlQuery.WriteString(" And ")
		}
		sqlQuery.WriteString(" PrimaryPhone like '%")
		sqlQuery.WriteString(searchCondition.PrimaryPhone.String)
		sqlQuery.WriteString("%'")
		putAndCondition = true
	}
	if len(searchCondition.PermCity.String) != 0 {
		if putAndCondition {
			sqlQuery.WriteString(" And ")
		}
		sqlQuery.WriteString(" PermCity like '%")
		sqlQuery.WriteString(searchCondition.PermCity.String)
		sqlQuery.WriteString("%'")
		putAndCondition = true
	}
	fmt.Println(sqlQuery.String())
	return sqlQuery.String()
}
