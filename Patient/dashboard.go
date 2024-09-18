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

	userId := c.Query("ID")
	StartDate := c.DefaultQuery("StartDate", firstOfMonth.Format("yyyy-mm-dd"))
	EndDate := c.DefaultQuery("EndDate", lastOfMonth.Format("yyyy-mm-dd"))

	var LoggedInUserRole string
	var db *gorm.DB = database.GetDBContext()
	db.Table("Patient").Where("ID = ?", userId).Select("Role").Scan(&LoggedInUserRole)
	var subQuery string = "Select * from PatientAppointment Where "
	switch LoggedInUserRole {
	case "Patient":
		subQuery = subQuery + " PatientId = " + userId + " AND (ApptDate Between '" + StartDate + "' AND '" + EndDate + "')"
	case "Doctor":
		subQuery = subQuery + " DoctorId = " + userId + " AND (ApptDate Between '" + StartDate + "' AND '" + EndDate + "')"
	case "Admin":
		subQuery = subQuery + " ApptDate Between '" + StartDate + "' AND '" + EndDate + "'"
	}
	var PatientAppointments []*PatientAppointment
	db.Table("PatientAppointment").Raw(subQuery).Scan(&PatientAppointments)
	c.IndentedJSON(http.StatusOK, PatientAppointments)
}
