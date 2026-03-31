package app

import (
	"log"
	"os"

	"api/src/common/filters"
	"api/src/common/i18n"
	"api/src/common/interceptors"
	"api/src/config"
	"api/src/modules/auth"
	"api/src/modules/setlists"
	"api/src/modules/songs"
	"api/src/modules/users"
	"api/src/modules/ws"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

func SetupApp() *gin.Engine {
	config.LoadEnv()

	config.ConnectDatabase()
	config.ConnectRedis()

	if err := i18n.Initialize("locales", language.English); err != nil {
		log.Fatalf("Failed to initialize i18n: %v", err)
	}

	if err := i18n.InitValidator(); err != nil {
		log.Fatalf("Failed to initialize validator i18n: %v", err)
	}

	router := gin.Default()

	router.Use(config.SetupCors())

	router.Use(i18n.Middleware())

	router.Use(interceptors.LoggerInterceptor())
	router.Use(filters.ErrorHandler())

	api := router.Group("/api/v1")

	users.InitModule(api)
	auth.InitModule(api)
	songs.InitModule(api)
	setlists.InitModule(api)
	ws.InitModule(api)

	return router
}

func Bootstrap() {
	router := SetupApp()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Application is starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
