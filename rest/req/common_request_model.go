package req

// DelModel 公共请求删除结构体
type DelModel struct {
	Ids []int64 `json:"ids" url:"ids"` //ids
}

// PageParam 分页请求参数
type PageParam struct {
	PageNumber int `url:"pageNumber"` //当前的页数
	PageSize   int `url:"pageSize"`   //每页的数据数量
}
