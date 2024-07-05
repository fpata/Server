package main

import (
	"clinic_server/albums"
	"clinic_server/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	userRouter := router.Group("/users")
	userRouter.GET("/", user.GetAllUsers)
	userRouter.GET("/:Id", user.GetUserById)
	userRouter.POST("/SearchByParams/", user.GetUserByParams)
	userRouter.GET("/GetBatchResult/", user.GetAllUsersWithPaging)
	userRouter.POST("/", user.CreateUser)
	userRouter.PUT("/", user.UpdateUser)
	userRouter.PATCH("/", user.PatchUser)
	userRouter.DELETE("/", user.DeleteUser)
	router.GET("/albums", albums.GetAlbums)

	router.Run("localhost:8088")
}
