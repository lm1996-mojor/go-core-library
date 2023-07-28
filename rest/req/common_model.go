package req

import (
	"time"

	"gorm.io/gorm"
)

// CommonModel 公共model
type CommonModel struct {
	Id        int64          `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int64" json:"id,omitempty"`                 //主键id
	CreateBy  int64          `gorm:"column:create_by;type:uint64" json:"createBy,omitempty"`                              //创建人
	UpdateBy  int64          `gorm:"column:update_by;type:uint64" json:"updateBy,omitempty"`                              //更新人
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime" json:"deletedAt,omitempty"`                          //软删标识（有值代表删除）
	CreatedAt *time.Time     `gorm:"<-:create;autoCreateTime;column:created_at;type:datetime" json:"createdAt,omitempty"` //创建时间
	UpdatedAt *time.Time     `gorm:"<-:update;autoUpdateTime;column:updated_at;type:datetime" json:"updatedAt,omitempty"` //更新时间
}

// GetCommonModelColumns 获取公共model定义列
func (commonModel CommonModel) GetCommonModelColumns() []string {
	columns := []string{"id", "create_by", "update_by", "created_at", "updated_at", "deleted_at"}
	return columns
}
