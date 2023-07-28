package common_repsitory

import (
	"gorm.io/gorm"
)

// SingleCheck 统一数据唯一性检查工具
func SingleCheck(db *gorm.DB, checkInfo CheckModel) bool {
	var resultMap map[string]interface{}
	db.Table(checkInfo.TableName).Where(checkInfo.SingleParamName+" = ?", checkInfo.Value).Find(&resultMap)
	if id, ok := resultMap["id"]; ok && id != 0 {
		return true
	}
	return false
}
