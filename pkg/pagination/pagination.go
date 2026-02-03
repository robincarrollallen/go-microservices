package pagination

import "math"

type PaginationParams struct {
	Page     int
	PageSize int
}

// NewPaginationParams 创建默认分页参数
func NewPaginationParams() PaginationParams {
	return PaginationParams{
		Page:     1,
		PageSize: 10,
	}
}

// NewPaginationParamsWithValues 从用户输入创建分页参数，缺失的字段使用默认值
func NewPaginationParamsWithValues(page, pageSize int) PaginationParams {
	params := NewPaginationParams()

	if page > 0 {
		params.Page = page
	}
	if pageSize > 0 {
		params.PageSize = pageSize
	}

	return params
}

// NewPaginationParamsFromPointers 从指针创建分页参数，处理可能为 nil 的值
func NewPaginationParamsFromPointers(page, pageSize *int) PaginationParams {
	params := NewPaginationParams()

	if page != nil && *page > 0 {
		params.Page = *page
	}
	if pageSize != nil && *pageSize > 0 {
		params.PageSize = *pageSize
	}

	return params
}

// CalculatePaginationParams 计算分页参数（从 PaginationParams 结构体）
func CalculatePaginationParams(params PaginationParams) (offset int, limit int) {
	page := params.Page
	pageSize := params.PageSize

	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page <= 0 {
		page = 1
	}

	offset = (page - 1) * pageSize
	return offset, pageSize
}

// CalculateTotalPages 计算总页数
func CalculateTotalPages(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(pageSize)))
}
