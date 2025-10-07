package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, message string, data interface{}) {
	response := Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

func Error(c *gin.Context, status int, message string, err error) {
	response := Response{
		Status:  "error",
		Message: message,
	}
	if err != nil {
		response.Data = err.Error()
	}
	c.JSON(status, response)
}
