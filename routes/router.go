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
	process := new(controllers.ProcessController)

	router := gin.Default()

	// recover from panics and respond with internal server error
	router.Use(gin.Recovery())

	// enabling cors
	config := cors.DefaultConfig()
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

    v1 := router.Group("/api/v1/platform")

    imageGroup := v1.Group("/image")
	imageGroup.GET("", image.GetAllImages)
    imageGroup.GET("/:corId", image.GetImageByCorId)
    imageGroup.GET("/name/:creatorName", image.GetImageByCreatorName)

    challengeGroup := v1.Group("/challenge")
	challengeGroup.GET("", challenge.GetAllChallenges)
    challengeGroup.GET("/:corId", challenge.GetChallengeByCorID)
    challengeGroup.GET("/name/:creatorName", challenge.GetChallengeByCreatorName)

    processGroup := v1.Group("/process")
	processGroup.GET("", process.GetAllProcesses)
    processGroup.GET("/:corId", process.GetProcessByCorID)
    processGroup.GET("/name/:creatorName", process.GetProcessByCreatorName)
    processGroup.GET("/image/:imageName", process.GetProcessByImageName)

	router.Run(":" + PORT)
}
