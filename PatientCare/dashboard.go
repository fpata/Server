package PatientCare

import (
	"clinic_server/database"
	"net/http"
	"time"

	"clinic_server/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func GetDashboardInformation(c *gin.Context) {
	logger.Init(zerolog.InfoLevel)

	db := database.GetDBContext()
	if db == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "database connection failed"})
		return
	}

	startDate, endDate, err := getDateRange(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.Query("ID")
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	var loggedInUserRole string
	if err := db.Table("Patient").Where("ID = ?", userId).Select("Role").Scan(&loggedInUserRole).Error; err != nil {
		logger.Error("Unable to get user role", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get user role"})
		return
	}

	query := `SELECT * FROM PatientAppointment WHERE ApptDate BETWEEN ? AND ?`
	params := []interface{}{startDate, endDate}

	switch loggedInUserRole {
	case "Patient":
		query += " AND PatientId = ?"
		params = append(params, userId)
	case "Doctor":
		query += " AND DoctorId = ?"
		params = append(params, userId)
	case "Admin":
		// No additional filters for admin
	default:
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid user role"})
		return
	}

	var patientAppointments []PatientAppointment
	if err := db.Raw(query, params...).Scan(&patientAppointments).Error; err != nil {
		logger.Error("Unable to get Patient Appointment information", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch appointments"})
		return
	}

	logger.Info("GetDashboardInformation Complete")
	c.JSON(http.StatusOK, patientAppointments)
}

func getDateRange(c *gin.Context) (startDate, endDate string, err error) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	startDate = c.DefaultQuery("StartDate", firstOfMonth.Format("2006-01-02"))
	endDate = c.DefaultQuery("EndDate", lastOfMonth.Format("2006-01-02"))

	_, err = time.Parse("2006-01-02", startDate)
	if err != nil {
		return "", "", err
	}
	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		return "", "", err
	}

	return startDate, endDate, nil
}
