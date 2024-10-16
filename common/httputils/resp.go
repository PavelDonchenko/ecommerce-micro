package httputils

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PavelDonchenko/ecommerce-micro/common/logger"
)

type ResponseErr struct {
	Status int           `json:"status"`
	Error  []interface{} `json:"error"`
}

func NewBadRequestErrorWithLog(c *gin.Context, err interface{}, log *logger.Logger) {
	log.Error(c.Request.Context(), "http error", slog.Any("bad request error", err))
	res := &ResponseErr{
		Status: http.StatusBadRequest,
		Error:  []interface{}{err},
	}

	c.JSON(http.StatusBadRequest, res)
}

func NewInternalErrorWIthLog(c *gin.Context, err interface{}, log *logger.Logger) {
	log.Error(c.Request.Context(), "http error", slog.Any("internal error", err))
	res := &ResponseErr{
		Status: http.StatusInternalServerError,
		Error:  []interface{}{err},
	}

	c.JSON(http.StatusInternalServerError, res)
}
