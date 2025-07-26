package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse 统一的API响应结构
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// SuccessResponseWithMessage 带自定义消息的成功响应
func SuccessResponseWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, APIResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequestResponse 400错误响应
func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, 400, message)
}

// UnauthorizedResponse 401错误响应
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, 401, message)
}

// ForbiddenResponse 403错误响应
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, 403, message)
}

// NotFoundResponse 404错误响应
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, 404, message)
}

// InternalServerErrorResponse 500错误响应
func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, 500, message)
}

// PaginationResponse 分页响应结构
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int64       `json:"total_pages"`
}

// SuccessPaginationResponse 分页成功响应
func SuccessPaginationResponse(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	response := PaginationResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "success",
		Data:    response,
	})
}