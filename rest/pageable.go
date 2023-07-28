package rest

// Pageable general page data
type Pageable struct {
	PageNumber int   `json:"pageNumber"` //当前的页数
	PageSize   int   `json:"pageSize"`   //每页的数据数量
	TotalItems int64 `json:"totalItems"` //数据总数量
	TotalPages int64 `json:"totalPages"` //总页数
}

// NewPageable return new Pageable info
func NewPageable(pageNumber int, pageSize int, totalItems int64) (page Pageable) {
	if pageNumber >= 1 {
		pageNumber = pageNumber + 1
	}
	if pageNumber <= 0 {
		pageNumber = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	totalPages := totalItems / int64(pageSize)
	if totalItems%int64(pageSize) > 0 {
		totalPages++
	}
	page = Pageable{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
	return
}
