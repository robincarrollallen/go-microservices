package response

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// AppError 应用错误类型
type AppError struct {
	Code       int    // 业务错误码
	Message    string // 错误信息
	HTTPStatus int    // HTTP 状态码
}

// NewAppError 创建应用错误
func NewAppError(code int, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	return e.Message
}

// 预定义的常用错误
var (
	ErrDuplicateName      = NewAppError(4001, "Tenant name already exists", http.StatusConflict)
	ErrTenantNotFound     = NewAppError(4004, "Tenant not found", http.StatusNotFound)
	ErrInvalidRequest     = NewAppError(4000, "Invalid request", http.StatusBadRequest)
	ErrInternalServer     = NewAppError(5000, "Internal server error", http.StatusInternalServerError)
	ErrDomainConflict     = NewAppError(4002, "Domain already exists", http.StatusConflict)
)

// FormatValidationError 格式化 validator 验证错误
// 支持单个和多个验证错误的处理，返回用户友好的错误信息
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

