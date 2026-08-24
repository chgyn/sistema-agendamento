package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedData struct {
	Items      interface{} `json:"items"`
	TotalCount int64       `json:"total_count"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
}

// Success envia resposta HTTP 200 com payload
func Success(c *gin.Context, data interface{}, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

// Created envia resposta HTTP 201 com payload criado
func Created(c *gin.Context, data interface{}, message ...string) {
	msg := "Criado com sucesso"
	if len(message) > 0 {
		msg = message[0]
	}
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

// Error envia resposta de erro com código HTTP
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
	})
}

// BadRequest envia resposta 400
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Unauthorized envia resposta 401
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// Forbidden envia resposta 403
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}

// NotFound envia resposta 404
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// Conflict envia resposta 409 (ex: duplo agendamento / colisão de horário)
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, message)
}

// InternalServerError envia resposta 500
func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}
