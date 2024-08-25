package Patient

import (
	"clinic_server/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDashboardInformation(c *gin.Context) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	userId := c.Query("Id")
	StartDate := c.DefaultQuery("StartDate", firstOfMonth.Format("dd-mmm-yyyy"))
	EndDate := c.DefaultQuery("EndDate", lastOfMonth.Format("dd-mmm-yyyy"))

	var LoggedInUserRole string
	var db *gorm.DB = database.GetDBContext()
	db.Table("Patient").Where("Id = ?", userId).Select("Role").Scan(&LoggedInUserRole)
	var subQuery string
	switch LoggedInUserRole {
	case "Patient":
		subQuery = "PatientId = ? AND (ApptDate Between ? AND ?)"
	case "Doctor":
		subQuery = "DoctorId = ? AND (ApptDate Between ? AND ?)"
	case "Admin":
		subQuery = "ApptDate Between ? AND ?"
	}
	var PatientAppointments []*PatientAppointment
	if LoggedInUserRole == "Admin" {
		db.Table("PatientAppointment").Where(subQuery, StartDate, EndDate).Find(&PatientAppointments)
	} else {
		db.Table("PatientAppointment").Where(subQuery, userId, StartDate, EndDate).Find(&PatientAppointments)
	}
	c.IndentedJSON(http.StatusOK, PatientAppointments)
}
