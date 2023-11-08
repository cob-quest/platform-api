package routes

import (
	"platform_api/configs"
	"platform_api/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func InitRoutes() {

	challenge := controllers.NewChallengeController(configs.Client)
	image := controllers.NewImageController(configs.Client)
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
	platformImage.GET("/status/:corId", process.GetProcessStatusByCorId)
	platformImage.POST("", image.UploadImage)

	platformChallenge := platform.Group("/challenge")
	platformChallenge.GET("", challenge.GetAllChallenges)
	platformChallenge.GET("/:corId", challenge.GetChallengeByCorID)
	platformChallenge.GET("/name/:creatorName", challenge.GetChallengeByCreatorName)
	platformChallenge.GET("/status/:corId", process.GetProcessStatusByCorId)
	platformChallenge.POST("", challenge.CreateChallenge)

	platformProcess := platform.Group("/process")
	platformProcess.GET("", process.GetAllProcesses)
	platformProcess.GET("/:corId", process.GetProcessByCorID)
	platformProcess.GET("/name/:creatorName", process.GetProcessByCreatorName)

	platformAttempt := platform.Group("/attempt")
	platformAttempt.POST("", attempt.StartAttempt)
	platformAttempt.GET("/status/:corId", process.GetProcessStatusByCorId)
	platformAttempt.GET("/:token", attempt.GetOneAttemptByToken)
	platformAttempt.POST("/submit", attempt.SubmitAttemptByToken)


	platformResult := platform.Group("/result")
	platformResult.POST("/:token", )

	// // trigger api
	// trigger := v1.Group("/trigger")
	// triggerImage := trigger.Group("/image")
	// triggerImage.POST("", image.UploadImage)

	// // challenge
	// triggerChallenge := trigger.Group("/challenge")
	// triggerChallenge.POST("", challenge.CreateChallenge)

	// // start challenge
	// triggerStartChallenge := trigger.Group("/startchallenge")
	// triggerStartChallenge.GET("", attempt.GetAttempt)
	// triggerStartChallenge.POST("", attempt.StartAttempt)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":" + configs.PORT)
}
