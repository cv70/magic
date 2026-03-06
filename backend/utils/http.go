package utils

import (
	"backend/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespSuccess[T any](c *gin.Context, data T) {
	resp := types.HTTPResponse[T]{
		Code: http.StatusOK,
		Data: data,
	}
	c.JSON(http.StatusOK, resp)
}

func RespSuccessNone(c *gin.Context) {
	resp := types.HTTPResponse[int]{
		Code: http.StatusOK,
		Data: 0,
	}
	c.JSON(http.StatusOK, resp)
}

func RespError(c *gin.Context, code int, msg string) {
	resp := types.HTTPResponse[int]{
		Code: code,
		Msg:  msg,
		Data: 0,
	}
	c.JSON(http.StatusOK, resp)
}
