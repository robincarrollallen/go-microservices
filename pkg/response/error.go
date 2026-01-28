package response

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// APIError 接口 - 业务错误需要实现此接口
type APIError interface {
	error // 实现 error 接口 (必须：可以作为 error 类型使用)
	GetCode() int    // 返回错误码
	GetMessage() string // 返回错误消息
	GetStatus() int    // 返回 HTTP 状态码
}

// AppError 应用错误类型
type AppError struct {
	Code    int    // 业务错误码
	Message string // 错误信息
	Status  int    // HTTP 状态码
}

// NewAppError 创建应用错误
func NewAppError(code int, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// Error 实现 error 接口 (必须：实现 error 类型方法)
func (e *AppError) Error() string {
	return e.Message
}

// GetCode 返回错误码
func (e *AppError) GetCode() int {
	return e.Code
}

// GetMessage 返回错误消息
func (e *AppError) GetMessage() string {
	return e.Message
}

// GetStatus 返回 HTTP 状态码
func (e *AppError) GetStatus() int {
	return e.Status
}

// FormatValidationError 格式化 validator 验证错误 支持单个和多个验证错误的处理，返回用户友好的错误信息
func FormatValidationError(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		if len(validationErrors) == 0 {
			return err.Error()
		}

		// 构建友好的错误消息
		messages := make([]string, len(validationErrors))
		for i, e := range validationErrors {
			// 输出格式: "Field 'Name' validation failed on tag 'required'"
			messages[i] = fmt.Sprintf("Field '%s' validation failed on tag '%s'", e.Field(), e.Tag())
		}
		return strings.Join(messages, "; ")
	}

	// 如果不是 validator 错误，尝试从错误字符串中提取
	errStr := err.Error()
	if idx := strings.Index(errStr, "Error:"); idx != -1 {
		return errStr[idx+6:] // 跳过 "Error:" 前缀
	}

	return errStr
}
