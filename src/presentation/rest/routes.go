package rest

import "github.com/gin-gonic/gin"

func (s *RestServer) registerRoutes(r *gin.Engine) {

	r.POST("/register", s.auth.Register)

	r.POST("/login", s.auth.Login)
}
