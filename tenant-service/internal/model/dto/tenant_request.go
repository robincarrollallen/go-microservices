package dto

import (
	"time"

	"shared.local/pkg/pagination"
)

type CreateTenantRequest struct {
	Name    string   `form:"name" binding:"required,min=1,max=100"`     // 商户名: 必传, 长度1-100
	Domains []string `form:"domains" binding:"omitempty,dive,required"` // 域名: 可不传或为空数组，但如果非空数组, 则元素必须是非空字符串
	Status  *uint8   `form:"status"`                                    // 指针类型，用于判断是否传值
}

type ListTenantRequest struct {
	pagination.PaginationRequest
	ID        *uint     `form:"id"`        // 商户ID: 可选
	Domain    string    `form:"domain"`    // 关联域名: 可选
	StartTime *time.Time `form:"startTime"` // 创建时间范围开始: 可选
	EndTime   *time.Time `form:"endTime"`   // 创建时间范围结束: 可选
}
