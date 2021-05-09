package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"apm.dev/go-simple-chat/src/data/repository"
	"apm.dev/go-simple-chat/src/data/storage/memory"
	"apm.dev/go-simple-chat/src/domain/authing"
	"apm.dev/go-simple-chat/src/presentation/rest"
	"apm.dev/go-simple-chat/src/presentation/rest/controllers"
	"github.com/joho/godotenv"
)

func main() {
	// configs, env variables
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	jwtKey, jwtExp := getJwtKeyAndExp()
	httpServeAddr := getEnvOrPanic("HTTP_HOST")

	// DataBases, MQs, APIs client
	userDS := memory.NewUserDS()

	// Repositories
	userRepo := repository.NewUserRepo(userDS)

	// Pseudo services
	jwtMng := authing.NewJWTManager(jwtKey, jwtExp)

	// Services
	authSvc := authing.NewService(userRepo, jwtMng)

	// Handlers(controllers,entry points)
	authCtrl := controllers.NewAuthController(authSvc)

	// Servers
	restSrv := rest.NewServer(authCtrl)

	// Start
	restSrv.Start(httpServeAddr)

	// Listen to OS interrupt signal to stop app
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	// Stop
	restSrv.Stop()
}

func getJwtKeyAndExp() (string, time.Duration) {
	jwtKey := getEnvOrPanic("JWT_KEY")
	jwtExpStr := getEnvOrPanic("JWT_EXP")

	jwtExpInt, err := strconv.ParseInt(jwtExpStr, 10, 64)
	if err != nil {
		panic(err)
	}
	return jwtKey, time.Minute * time.Duration(jwtExpInt)
}

func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("%s env var is required", key))
	}
	return value
}
