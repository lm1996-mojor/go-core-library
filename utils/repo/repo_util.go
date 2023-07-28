package repo

import (
	"reflect"

	"mojor/go-core-library/exception"

	"gorm.io/gorm"
)

const (
	ErrFoundMessage = "entity found err"
	ErrCountMessage = "entity count err"
)

func GetByCondition(db *gorm.DB, tableTypePtr *reflect.Type, itemPtr interface{}, funcSlice []func(db *gorm.DB) *gorm.DB) {
	tableName := ObtainTableNameByType(tableTypePtr)
	err := db.Table(tableName).Scopes(funcSlice...).Find(itemPtr).Error
	exception.CheckErr(err, 500, ErrFoundMessage)
	return
}

func GetMainAndSubCond(db *gorm.DB, itemPtr interface{}, funcSlice []func(db *gorm.DB) *gorm.DB) (total int64) {
	if len(funcSlice) != 0 {
		db = db.Scopes(funcSlice...)
	}
	ty := reflect.TypeOf(itemPtr).Elem().Elem()
	db.Model(reflect.New(ty).Interface()).Count(&total)
	err := db.Find(itemPtr).Error
	exception.CheckErr(err, 500, ErrFoundMessage)
	return
}

func ExistsByCondition(db *gorm.DB, tableTypePtr *reflect.Type, funcSlice []func(db *gorm.DB) *gorm.DB) (exist bool) {
	num := CountByCondition(db, tableTypePtr, funcSlice)
	exist = num > 0
	return
}

func CountByCondition(db *gorm.DB, tableTypePtr *reflect.Type, funcSlice []func(db *gorm.DB) *gorm.DB) (num int64) {
	tableName := ObtainTableNameByType(tableTypePtr)
	if len(funcSlice) == 0 {
		db.Table(tableName).Count(&num)
		return
	}
	err := db.Table(tableName).Scopes(funcSlice...).Count(&num).Error
	exception.CheckErr(err, 500, ErrCountMessage)
	return
}

func DeleteByCondition(tx *gorm.DB, tableTypePtr *reflect.Type, funcSlice []func(db *gorm.DB) *gorm.DB) {
	if len(funcSlice) == 0 {
		return
	}
	tableName := ObtainTableNameByType(tableTypePtr)
	tx.Table(tableName).Scopes(funcSlice...).Delete(reflect.New(*tableTypePtr))
}

func CreateAll(db *gorm.DB, tableTypePtr *reflect.Type, itemPtr interface{}) {
	tableName := ObtainTableNameByType(tableTypePtr)
	valSlice := reflect.ValueOf(itemPtr)
	db.Table(tableName).CreateInBatches(itemPtr, valSlice.Len())
}

func ObtainImpossibleScope() (fnSlice []func(db *gorm.DB) *gorm.DB) {
	fn := func(db *gorm.DB) *gorm.DB { return db.Where("1=2") }
	fnSlice = make([]func(db *gorm.DB) *gorm.DB, 1)
	fnSlice[0] = fn
	return
}

func ObtainInevitableScope() (fnSlice []func(db *gorm.DB) *gorm.DB) {
	fn := func(db *gorm.DB) *gorm.DB { return db.Where("1=1") }
	fnSlice = make([]func(db *gorm.DB) *gorm.DB, 1)
	fnSlice[0] = fn
	return
}

// Paginate db分页封装
func Paginate(pageNumber int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := pageNumber
		pageSize := pageSize
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
