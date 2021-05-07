package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {

	r.POST("/users", register)
	r.GET("/users", allUsers)
	r.GET("/chats/:sender", getChats)
	r.GET("/chats/:sender/:receiver", getChatMessages)
	r.GET("/ws/:sender/:receiver", startChat)
}
