package main

import (
	"apm.dev/go-simple-chat/chats"
	"apm.dev/go-simple-chat/http"
	"github.com/gin-gonic/gin"
)

func main() {

	h := chats.GetHub()
	go h.Run()

	r := gin.Default()

	http.RegisterRoutes(r)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}
