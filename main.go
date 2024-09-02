package main

import (
	"clinic_server/Patient"
	"clinic_server/albums"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	patientRouter := router.Group("/patients")
	patientRouter.GET("/:ID", Patient.GetPatientById)
	patientRouter.POST("/SearchByParams/", Patient.GetPatientByParams)
	patientRouter.POST("/", Patient.CreatePatient)
	patientRouter.PUT("/", Patient.UpdatePatient)
	patientRouter.DELETE("/", Patient.DeletePatient)
	router.GET("/albums", albums.GetAlbums)
	router.POST("/login/", Patient.ValidateLogin)
	router.GET("/dashboard", Patient.GetDashboardInformation)
	router.Run("localhost:8088")
}
