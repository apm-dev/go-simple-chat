package controllers

import (
	"errors"
	"net/http"

	"apm.dev/go-simple-chat/src/domain"
	"apm.dev/go-simple-chat/src/domain/authing"
	request "apm.dev/go-simple-chat/src/presentation/rest/requests"
	"apm.dev/go-simple-chat/src/presentation/rest/responses"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	svc authing.Service
}

func NewAuthController(svc authing.Service) *AuthController {
	return &AuthController{svc}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var data request.Register
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.Make(http.StatusUnprocessableEntity, err.Error(), nil))
		return
	}

	id, err := ctrl.svc.Register(data.Name, data.Email, data.Password)
	if err != nil {
		code := http.StatusBadRequest
		if errors.Is(err, domain.ErrInternalServer) {
			code = http.StatusInternalServerError
		}
		c.JSON(code, responses.Make(code, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, responses.Make(http.StatusOK, "welcome", gin.H{"id": id}))
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var data request.Login
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.Make(http.StatusUnprocessableEntity, err.Error(), nil))
		return
	}

	token, err := ctrl.svc.Login(data.Email, data.Password)
	if err != nil {
		code := http.StatusBadRequest
		if errors.Is(err, domain.ErrInternalServer) {
			code = http.StatusInternalServerError
		}
		c.JSON(code, responses.Make(code, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, responses.Make(http.StatusOK, "", gin.H{"token": token}))
}
