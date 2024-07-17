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
	patientRouter.GET("/", Patient.GetAllPatients)
	patientRouter.GET("/:Id", Patient.GetPatientById)
	patientRouter.POST("/SearchByParams/", Patient.GetPatientByParams)
	patientRouter.GET("/GetBatchResult/", Patient.GetAllPatientsWithPaging)
	patientRouter.POST("/", Patient.CreatePatient)
	patientRouter.PUT("/", Patient.UpdatePatient)
	patientRouter.PATCH("/", Patient.PatchPatient)
	patientRouter.DELETE("/", Patient.DeletePatient)
	router.GET("/albums", albums.GetAlbums)

	router.Run("localhost:8088")
}
