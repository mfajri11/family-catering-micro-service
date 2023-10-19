package HandlerHttp

import (
	"fmt"
	"net/http"

	"github.com/mfajri11/family-catering-micro-service/auth-service/internal/core/domain"
	"github.com/mfajri11/family-catering-micro-service/auth-service/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type AuthHTTPHandler struct {
	authSvc ports.IAuthService
}

func (h AuthHTTPHandler) Login(c *gin.Context) {
	req := LoginRequest{}

	defer c.Request.Body.Close()
	err := c.ShouldBindJSON(&req)
	if err != nil {
		fmt.Printf("HandlerHTTP.AuthHTTPHandler.Login: fail bind json: %s", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, nil) // TODO: implement error response
		return
	}

	resp, err := h.authSvc.Login(c, domain.LoginRequest{
		Email:             req.Email,
		EncryptedPassword: req.Password,
	})

	if err != nil {
		fmt.Printf("HandlerHTTP.AuthHTTPHandler.Login: fail bind json: %s", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, nil) // TODO: implement error response
		return
	}

	c.JSON(http.StatusOK, resp)

}
