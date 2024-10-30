package PatientCare

import (
	"clinic_server/database"
	"net/http"
	"reflect"

	"clinic_server/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func GetPatientById(c *gin.Context) {
	var err error
	logger.Init(zerolog.InfoLevel)

	patientId := c.Param("ID")
	var patientViewModel PatientViewModel = PatientViewModel{}
	var db *gorm.DB = database.GetDBContext()
	err = db.Find(&patientViewModel.Patient, patientId).Error
	if err != nil {
		logger.Error("Unable to get patient information", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = db.Where("PatientID = ?", patientId).Find(&patientViewModel.PatientAppointments).Error
	if err != nil {
		logger.Error("Unable to get patient Appointment", err)
	}

	err = db.Where("PatientID = ?", patientId).Find(&patientViewModel.PatientTreatments).Error
	if err != nil {
		logger.Error("Unable to get patient PatientTreatments", err)
	}

	err = db.Where("PatientID = ?", patientId).Find(&patientViewModel.PatientTreatmentDetails).Error
	if err != nil {
		logger.Error("Unable to get patient PatientTreatmentDetails", err)
	}

	err = db.Where("PatientID = ?", patientId).Find(&patientViewModel.PatientReports).Error
	if err != nil {
		logger.Error("Unable to get patient PatientReports", err)
	}

	if err == nil {
		c.IndentedJSON(http.StatusOK, patientViewModel)
	} else {
		c.IndentedJSON(http.StatusFailedDependency, patientViewModel)
	}
}

func CreatePatient(c *gin.Context) {
	logger.Init(zerolog.InfoLevel)
	var err error
	var patientViewModel PatientViewModel
	err = c.ShouldBind(&patientViewModel)
	if err != nil {
		logger.Error("Unable to bind Patient JSON", err)
	}
	var db *gorm.DB = database.GetDBContext()
	err = db.Create(&patientViewModel.Patient).Error
	if err != nil {
		logger.Error("Unable to get Patient Appointment information", err)
	}
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
	logger.Init(zerolog.InfoLevel)
	c.ShouldBind(&patientViewModel)
	var db *gorm.DB = database.GetDBContext()
	var err = db.Save(&patientViewModel.Patient).Error
	if err != nil {
		logger.Error("Unable to Update Patient information", err)
	}
	SavePatientArrays(patientViewModel.PatientAppointments, db, patientViewModel)
	SavePatientArrays(patientViewModel.PatientReports, db, patientViewModel)
	SavePatientArrays(patientViewModel.PatientTreatments, db, patientViewModel)
	SavePatientArrays(patientViewModel.PatientTreatmentDetails, db, patientViewModel)
	c.IndentedJSON(http.StatusOK, &patientViewModel)
}

func DeletePatient(c *gin.Context) {
	patientId := c.Query("ID")
	logger.Init(zerolog.InfoLevel)
	var db *gorm.DB = database.GetDBContext()
	var err = db.Delete(&Patient{}, patientId).Error
	if err != nil {
		logger.Error("Unable to delete Patient information", err)
	}
	c.IndentedJSON(http.StatusOK, "")
}

func UpdatePatientIDInArrays[pa PatientArray](patientArray []*pa, patientId int64, pvm PatientViewModel) {
	if len(patientArray) > 0 {
		for _, arrayVal := range patientArray {
			reflect.ValueOf(arrayVal).Elem().FieldByName("PatientID").SetInt(patientId)
		}
	}
}

func SavePatientArrays[pa PatientArray](patientArray []*pa, db *gorm.DB, pvm PatientViewModel) {
	var intVal int64 = 0
	var initialPatientTreatmentId int64 = 0
	var err error
	logger.Init(zerolog.InfoLevel)
	if len(patientArray) > 0 {
		var isPatientTreatment bool = (reflect.TypeOf(patientArray[0]).String() == "*Patient.PatientTreatment")
		for _, arrayVal := range patientArray {
			var value = reflect.ValueOf(arrayVal).Elem().FieldByName("ID").Int()
			if value <= intVal {
				if isPatientTreatment {
					initialPatientTreatmentId = reflect.ValueOf(arrayVal).Elem().FieldByName("ID").Int()
				}
				reflect.ValueOf(arrayVal).Elem().FieldByName("ID").SetInt(intVal)
			}
			err = db.Save(&arrayVal).Error
			if err != nil {
				logger.Error("Unable to save "+(reflect.TypeOf(arrayVal)).String(), err)
			}
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
