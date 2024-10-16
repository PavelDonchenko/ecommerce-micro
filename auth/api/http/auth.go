package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PavelDonchenko/ecommerce-micro/common/httputils"
	"github.com/PavelDonchenko/ecommerce-micro/common/logger"
	trace "github.com/PavelDonchenko/ecommerce-micro/common/trace/otel"
	"github.com/PavelDonchenko/ecommerce-micro/common/validation"

	"github.com/PavelDonchenko/ecommerce-micro/auth/service"
)

type Auth struct {
	log     *logger.Logger
	srvAuth service.Auth
}

func NewAuth(log *logger.Logger, srvAuth service.Auth) *Auth {
	return &Auth{log: log, srvAuth: srvAuth}
}

func (a *Auth) SignUp(c *gin.Context) {
	ctx, span := trace.NewSpan(c.Request.Context(), "AuthController.SignUp")
	defer span.End()

	var err error
	var newUser NewUser
	if err = c.ShouldBind(&newUser); err != nil {
		httputils.NewBadRequestErrorWithLog(c, err.Error(), a.log)
		return
	}
	validErr := validation.Validate(newUser)
	if validErr != nil {
		httputils.NewBadRequestErrorWithLog(c, validErr, a.log)
		return
	}

	u, err := toServiceNewUser(newUser)
	if err != nil {
		httputils.NewBadRequestErrorWithLog(c, err.Error(), a.log)
		return
	}

	a.log.Info(ctx, "HELOOOO")
	createdUser, err := a.srvAuth.CreateUser(ctx, u)
	if err != nil {
		httputils.NewInternalErrorWIthLog(c, err.Error(), a.log)
		return
	}

	c.JSON(http.StatusCreated, toAppUser(createdUser))
}
