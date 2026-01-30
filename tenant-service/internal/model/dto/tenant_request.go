package dto

type CreateTenantRequest struct {
	Name    string   `form:"name" binding:"required,min=1,max=100"`     // 商户名: 必传, 长度1-100
	Domains []string `form:"domains" binding:"omitempty,dive,required"` // 域名: 可不传或为空数组，但如果非空数组, 则元素必须是非空字符串
	Status  *uint8   `form:"status"`                                    // 指针类型，用于判断是否传值
}
