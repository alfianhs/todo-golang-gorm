package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"golang-gorm/app/config"
	"golang-gorm/app/delivery/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName(".env")
	viper.AddConfigPath("./")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func main() {
	// timeout
	timeoutContext := time.Duration(viper.GetInt("TIMEOUT")) * time.Second

	// init logger
	writers := make([]io.Writer, 0)
	if viper.GetString("LOG_TO_STDOUT") == "true" {
		writers = append(writers, io.Writer(os.Stdout))
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(io.MultiWriter(writers...))

	// set gin writer to logrus
	gin.DefaultWriter = logrus.StandardLogger().Writer()

	// gin mode release if env production
	if viper.GetString("GO_ENV") == "production" || viper.GetString("GO_ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// connect db
	db := config.GetDB(viper.GetString("DB_URL"))

	// init gin
	ginEngine := gin.New()

	// add logger
	ginEngine.Use(middleware.Logger(io.MultiWriter(writers...)))

	// set bootstrap
	bootstrapConfig := config.BootstrapConfig{
		GinEngine: ginEngine,
		DB:        db,
		Validator: config.NewValidator(),
		Timeout:   timeoutContext,
	}
	config.Bootstrap(bootstrapConfig)

	// cors
	ginEngine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}))

	// get port
	port := viper.GetString("PORT")

	// run gin
	ginEngine.Run(":" + port)
}
