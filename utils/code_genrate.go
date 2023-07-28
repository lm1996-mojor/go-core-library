package utils

import (
	"strings"
	"time"

	"mojor/go-core-library/utils/repo"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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
	Id        int            `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int" json:"id,omitempty"`                   //主键id
	CreateBy  int64          `gorm:"column:create_by;type:uint64" json:"createBy,omitempty"`                              //创建人
	UpdateBy  int64          `gorm:"column:update_by;type:uint64" json:"updateBy,omitempty"`                              //更新人
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime" json:"deletedAt,omitempty"`                          //软删标识（有值代表删除）
	CreatedAt *time.Time     `gorm:"<-:create;autoCreateTime;column:created_at;type:datetime" json:"createdAt,omitempty"` //创建时间
	UpdatedAt *time.Time     `gorm:"<-:update;autoUpdateTime;column:updated_at;type:datetime" json:"updatedAt,omitempty"` //更新时间
	PrefixStr string         `gorm:"column:prefix_str;type:string" json:"prefixStr"`
	Remark    string         `gorm:"column:remark;type:string" json:"remark"`
}

// 获取编码前缀
func obtainCodePrefixText(codeType int) string {
	PrefixStr := ""
	repo.ObtainCustomDbByDbName("platform_management").Table("code_prefix").Where("search_code = ?", codeType).Select("prefix_str").Scan(&PrefixStr)
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
