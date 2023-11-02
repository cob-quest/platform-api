package main

import (
	"platform_api/configs"
	"platform_api/mq"
	"platform_api/routes"
	"platform_api/services"
	_ "platform_api/docs"
)

func main() {
	configs.InitEnv()   // init env
	services.Init()     // init s3
	mq.Init()           // init rabbitmq connection
	routes.InitRoutes() // init controller routes
}
