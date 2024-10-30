package PatientCare

import (
	"clinic_server/database"
	"net/http"
	"time"

	"clinic_server/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func GetDashboardInformation(c *gin.Context) {
	var subQuery string = "Select * from PatientAppointment Where "
	var LoggedInUserRole string
	logger.Init(zerolog.InfoLevel)
	var db *gorm.DB = database.GetDBContext()
	var error error
	var PatientAppointments []PatientAppointment

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	userId := c.Query("ID")
	StartDate := c.DefaultQuery("StartDate", firstOfMonth.Format("yyyy-mm-dd"))
	EndDate := c.DefaultQuery("EndDate", lastOfMonth.Format("yyyy-mm-dd"))

	error = db.Table("Patient").Where("ID = ?", userId).Select("Role").Scan(&LoggedInUserRole).Error
	if error != nil {
		logger.Error("Unable to get Patient Appointment information", error)
	} else {
		logger.Info("GetDashboardInformation Complete")
	}

	switch LoggedInUserRole {
	case "Patient":
		subQuery = subQuery + " PatientId = " + userId + " AND (ApptDate Between '" + StartDate + "' AND '" + EndDate + "')"
	case "Doctor":
		subQuery = subQuery + " DoctorId = " + userId + " AND (ApptDate Between '" + StartDate + "' AND '" + EndDate + "')"
	case "Admin":
		subQuery = subQuery + " ApptDate Between '" + StartDate + "' AND '" + EndDate + "'"
	}
	if error == nil {
		error = db.Table("PatientAppointment").Raw(subQuery).Scan(&PatientAppointments).Error
		if error != nil {
			logger.Error("Unable to get Patient Appointment information", error)
		} else {
			logger.Info("GetDashboardInformation Complete")
		}
	}
	if error == nil {
		c.IndentedJSON(http.StatusOK, PatientAppointments)
	} else {
		c.AbortWithError(http.StatusInternalServerError, error)
	}

}
