package PatientCare
import (
	"clinic_server/database"
	"clinic_server/logger"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateClinic(c *gin.Context) {
	var clinic database.Clinic
	if err := c.ShouldBindJSON(&createClinicModel); err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&clinic).Error; err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, clinic)
}