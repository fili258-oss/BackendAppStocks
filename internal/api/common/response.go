package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response estructura estándar de respuesta API
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
}

// ErrorInfo información detallada del error
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// MetaInfo información adicional (paginación, etc)
type MetaInfo struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"page_size,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	TotalItems int `json:"total_items,omitempty"`
}

// SuccessResponse crea una respuesta exitosa
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

// SuccessResponseWithMeta crea una respuesta exitosa con metadata
func SuccessResponseWithMeta(c *gin.Context, statusCode int, data interface{}, meta *MetaInfo) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// ErrorResponse crea una respuesta de error
func ErrorResponse(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

// ErrorResponseWithDetails crea una respuesta de error con detalles
func ErrorResponseWithDetails(c *gin.Context, statusCode int, code, message, details string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// ValidationErrorResponse crea una respuesta de error de validación
func ValidationErrorResponse(c *gin.Context, details string) {
	ErrorResponseWithDetails(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input", details)
}

// NotFoundResponse crea una respuesta de recurso no encontrado
func NotFoundResponse(c *gin.Context, resource string) {
	ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", resource+" not found")
}

// InternalErrorResponse crea una respuesta de error interno
func InternalErrorResponse(c *gin.Context, err error) {
	ErrorResponseWithDetails(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", err.Error())
}
