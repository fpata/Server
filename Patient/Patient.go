package Patient

import (
	"clinic_server/database"
	"net/http"
	"reflect"

	"clinic_server/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type dummy struct {
	inty types.NullInt64
}

func GetPatientById(c *gin.Context) {
	patientId := c.Param("ID")
	var patientViewModel PatientViewModel = PatientViewModel{}
	var db *gorm.DB = database.GetDBContext()
	db.Find(&patientViewModel.Patient, patientId)
	db.Where("PatientID = ?", patientId).Find(&patientViewModel.PatientAppointments)
	db.Where("PatientID = ?", patientId).Find(&patientViewModel.PatientTreatments)
	db.Where("PatientID = ?", patientId).Find(&patientViewModel.PatientTreatmentDetails)
	db.Where("PatientID = ?", patientId).Find(&patientViewModel.PatientReports)
	c.IndentedJSON(http.StatusOK, patientViewModel)
}

func CreatePatient(c *gin.Context) {
	var patientViewModel PatientViewModel
	c.ShouldBind(&patientViewModel)
	var db *gorm.DB = database.GetDBContext()
	db.Create(&patientViewModel.Patient)
	var patientId = patientViewModel.Patient.ID
	UpdatePatientIDInArrays(patientViewModel.PatientAppointments, patientId, patientViewModel)
	UpdatePatientIDInArrays(patientViewModel.PatientReports, patientId, patientViewModel)
	UpdatePatientIDInArrays(patientViewModel.PatientTreatments, patientId, patientViewModel)
	UpdatePatientIDInArrays(patientViewModel.PatientTreatmentDetails, patientId, patientViewModel)
	SavePatientArrays(patientViewModel.PatientAppointments, db, patientViewModel)
	SavePatientArrays(patientViewModel.PatientReports, db, patientViewModel)
	SavePatientArrays(patientViewModel.PatientTreatments, db, patientViewModel)
	SavePatientArrays(patientViewModel.PatientTreatmentDetails, db, patientViewModel)
	c.IndentedJSON(http.StatusOK, &patientViewModel)
}

func UpdatePatient(c *gin.Context) {
	var patientViewModel PatientViewModel
	c.ShouldBind(&patientViewModel)
	var db *gorm.DB = database.GetDBContext()
	db.Save(&patientViewModel.Patient)
	SavePatientArrays(patientViewModel.PatientAppointments, db, patientViewModel)
	SavePatientArrays(patientViewModel.PatientReports, db, patientViewModel)
	SavePatientArrays(patientViewModel.PatientTreatments, db, patientViewModel)
	SavePatientArrays(patientViewModel.PatientTreatmentDetails, db, patientViewModel)
	c.IndentedJSON(http.StatusOK, &patientViewModel)
}

func DeletePatient(c *gin.Context) {
	patientId := c.Query("ID")
	var db *gorm.DB = database.GetDBContext()
	db.Delete(&Patient{}, patientId)
	c.IndentedJSON(http.StatusOK, "")
}

func UpdatePatientIDInArrays[pa PatientArray](patientArray []*pa, patientId int64, pvm PatientViewModel) {
	if patientArray != nil || len(patientArray) > 0 {
		for _, arrayVal := range patientArray {
			reflect.ValueOf(arrayVal).Elem().FieldByName("PatientID").SetInt(patientId)
		}
	}
}

func SavePatientArrays[pa PatientArray](patientArray []*pa, db *gorm.DB, pvm PatientViewModel) {
	var intVal int64 = 0
	var initialPatientTreatmentId int64 = 0
	if patientArray != nil || len(patientArray) > 0 {
		var isPatientTreatment bool = (reflect.TypeOf(patientArray[0]).String() == "*Patient.PatientTreatment")
		for _, arrayVal := range patientArray {
			var value = reflect.ValueOf(arrayVal).Elem().FieldByName("ID").Int()
			if value <= intVal {
				if isPatientTreatment {
					initialPatientTreatmentId = reflect.ValueOf(arrayVal).Elem().FieldByName("ID").Int()
				}
				reflect.ValueOf(arrayVal).Elem().FieldByName("ID").SetInt(intVal)
			}
			db.Save(&arrayVal)
			if isPatientTreatment {
				for _, ptd := range pvm.PatientTreatmentDetails {
					if ptd.PatientTreatmentID == initialPatientTreatmentId && initialPatientTreatmentId <= 0 {
						ptd.PatientTreatmentID = reflect.ValueOf(arrayVal).Elem().FieldByName("ID").Int()
					}
				}
			}
		}
	}
}

type PatientArray interface {
	PatientAppointment | PatientReport | PatientTreatment | PatientTreatmentDetail
}
