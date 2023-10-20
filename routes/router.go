package routes

import (
	"os"
	"platform_api/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRoutes() {
	PORT := os.Getenv("SERVER_PORT")

	challenge := new(controllers.ChallengeBuilderController)
	image := new(controllers.ImageController)

	router := gin.Default()

	// recover from panics and respond with internal server error
	router.Use(gin.Recovery())

	// enabling cors
	config := cors.DefaultConfig()
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.GET("/api/v1/platform/image", image.GetAllImages)
	router.GET("/api/v1/platform/image/:id", image.GetImageById)
	//router.GET("/api/v1/platform/image/:email", challenge.GetChallengeByEmail)
	router.GET("/api/v1/platform/challenge", challenge.GetAllChallenges)
	router.GET("/api/v1/platform/challenge/:corId", challenge.GetChallengeByCorID)

	router.Run(":" + PORT)
}
