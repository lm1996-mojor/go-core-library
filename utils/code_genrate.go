package utils

import (
	"strings"

	"github.com/lm1996-mojor/go-core-library/rest/req"
	"github.com/lm1996-mojor/go-core-library/utils/repo"

	uuid "github.com/satori/go.uuid"
)

// 编码前缀集
//
//	var codePrefix = map[int]string{
//		//1: "APP",      // 应用
//		//2: "SER_FUNC", // 业务
//		//3: "C",        // 控件
//		//4: "S",        // 子集
//		//5: "LK",       // 审批
//		//6: "sms",      // 短信
//	}

type CodePrefix struct {
	req.CommonModel
	PrefixStr string `gorm:"column:prefix_str;type:string" json:"prefixStr"`
	Remark    string `gorm:"column:remark;type:string" json:"remark"`
}

// 获取编码前缀
func obtainCodePrefixText(codeType int) string {
	PrefixStr := ""
	repo.ObtainCustomDbByDbName("platform_management").Table("code_prefix").Where("id = ?", codeType).Select("prefix_str").Scan(&PrefixStr)
	if PrefixStr == "" {
		PrefixStr = "link_ease"
	}
	return PrefixStr
}

// GenerateCodeByUUID 根据UUID生成编码
func GenerateCodeByUUID(codeType int) (code string) {
	uuId := uuid.NewV4()
	idStr := uuId.String()
	idStr = strings.ToUpper(strings.ReplaceAll(idStr, "-", ""))
	return obtainCodePrefixText(codeType) + idStr
}
