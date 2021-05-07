package http

import (
	"fmt"
	"net/http"

	"apm.dev/go-simple-chat/chats"
	"apm.dev/go-simple-chat/users"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func register(c *gin.Context) {
	var user users.User
	// validation added in User struct
	c.Bind(&user)
	// add user
	err := users.Add(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Message: fmt.Sprintf("Welcome %s!", user.Name),
	})
}

func allUsers(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Data: users.GetAll(),
	})
}

func getChats(c *gin.Context) {
	uname := c.Param("sender")
	if len(uname) < 3 {
		c.JSON(http.StatusUnprocessableEntity, Response{
			Message: "sender is required",
		})
	}

	chats, err := chats.GetChatsOfUser(uname)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{Data: chats})
}

func getChatMessages(c *gin.Context) {
	uname := c.Param("sender")
	cname := c.Param("receiver")
	if len(uname) < 3 || len(cname) < 3 {
		c.JSON(http.StatusUnprocessableEntity, Response{
			Message: "sender and receiver are required",
		})
	}

	msgs, err := chats.GetUserChatMessages(uname, cname)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{Data: msgs})
}

func startChat(c *gin.Context) {
	from := c.Param("sender")
	to := c.Param("receiver")

	if len(from) < 3 || len(to) < 3 {
		c.JSON(http.StatusUnprocessableEntity, Response{
			Message: "sender and receiver are required",
		})
	}

	chats.ServeWs(c.Writer, c.Request, from, to)
}
