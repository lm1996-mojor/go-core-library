package req

// DelModel 公共请求删除结构体
type DelModel struct {
	Ids []int64 `json:"ids"` //ids
}

// PageParam 分页请求参数
type PageParam struct {
	PageNumber int `json:"pageNumber"` //当前的页数
	PageSize   int `json:"pageSize"`   //每页的数据数量
}
