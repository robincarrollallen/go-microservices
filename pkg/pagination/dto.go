package pagination

type PaginationRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"pageSize" binding:"omitempty,min=1,max=100"`
}

// ToPaginationParams 将 PaginationRequest 转换为 PaginationParams
func (r *PaginationRequest) ToPaginationParams() PaginationParams {
	return NewPaginationParamsWithValues(r.Page, r.PageSize)
}

type PaginationResponse[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
}
