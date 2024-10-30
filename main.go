package main

import (
	"clinic_server/PatientCare"
	"clinic_server/albums"
	"clinic_server/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	patientRouter := router.Group("/patients")
	patientRouter.GET("/:ID", PatientCare.GetPatientById)
	patientRouter.POST("/SearchByParams/", PatientCare.GetPatientByParams)
	patientRouter.POST("/", PatientCare.CreatePatient)
	patientRouter.PUT("/", PatientCare.UpdatePatient)
	patientRouter.DELETE("/", PatientCare.DeletePatient)
	router.GET("/albums", albums.GetAlbums)
	router.POST("/login/", PatientCare.ValidateLogin)
	router.GET("/dashboard", PatientCare.GetDashboardInformation)
	router.NoRoute(func(c *gin.Context) {
		c.String(404, "Route Not Found")
	})

	config := config.GetConfiguration()
	urlp := config.Server.ServerUrl + ":" + config.Server.Port
	router.Run(urlp)
}
