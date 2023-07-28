package common_repsitory

// CheckModel 用于数据唯一性检查的结构体
type CheckModel struct {
	TableName       string      //需要检查的表
	SingleParamName string      `json:"singleParamName"` //需要检查的唯一性字段，数据库同名
	Value           interface{} `json:"value"`           //判断值(任意值)
}
