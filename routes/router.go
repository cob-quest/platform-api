package routes

import (
	"platform_api/configs"
	"platform_api/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRoutes() {

	challenge := new(controllers.ChallengeController)
	image := new(controllers.ImageController)
	process := new(controllers.ProcessController)
	attempt := new(controllers.AttemptController)

	router := gin.Default()

	// recover from panics and respond with internal server error
	router.Use(gin.Recovery())

	// enabling cors
	config := cors.DefaultConfig()
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	v1 := router.Group("/api/v1")

	// platform api
	platform := v1.Group("/platform")

	platformImage := platform.Group("/image")
	platformImage.GET("", image.GetAllImages)
	platformImage.GET("/:corId", image.GetImageByCorId)
	platformImage.GET("/name/:creatorName", image.GetImageByCreatorName)
	platformImage.GET("/status", process.GetProcessStatusByCorId)

	platformChallenge := platform.Group("/challenge")
	platformChallenge.GET("", challenge.GetAllChallenges)
	platformChallenge.GET("/:corId", challenge.GetChallengeByCorID)
	platformChallenge.GET("/name/:creatorName", challenge.GetChallengeByCreatorName)
	platformChallenge.GET("/status", process.GetProcessStatusByCorId)

	platformProcess := platform.Group("/process")
	platformProcess.GET("", process.GetAllProcesses)
	platformProcess.GET("/:corId", process.GetProcessByCorID)
	platformProcess.GET("/name/:creatorName", process.GetProcessByCreatorName)

    platformAttempt := platform.Group("/attempt")
    platformAttempt.GET("/status", process.GetProcessStatusByCorId)

	// trigger api
	trigger := v1.Group("/trigger")
	triggerImage := trigger.Group("/image")
	triggerImage.POST("", image.UploadImage)

	// challenge
	triggerChallenge := trigger.Group("/challenge")
	triggerChallenge.POST("", challenge.CreateChallenge)

	// start challenge
	triggerStartChallenge := trigger.Group("/startchallenge")
	triggerStartChallenge.GET("", attempt.GetAttempt)
	triggerStartChallenge.POST("", attempt.StartAttempt)

	router.Run(":" + configs.PORT)
}