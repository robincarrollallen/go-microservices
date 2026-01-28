package dto

type CreateDomainRequest struct {
	TenantID uint   `json:"tenantID" binding:"required"`             // 关联商户ID: 必传
	Domain   string `json:"domain" binding:"required,min=1,max=100"` // 域名: 必传, 长度1-100
	Status   *uint8 `json:"status"`                                  // 指针类型，用于判断是否传值
}
