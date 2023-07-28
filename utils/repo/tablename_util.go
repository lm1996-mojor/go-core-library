package repo

import (
	"fmt"
	"reflect"
	"strings"
)

var typeTableNameMap = make(map[string]string)

type TableName interface {
	TableName() string
}

func registerTableName(tPtr *reflect.Type, tableName string) {
	key := obtainTableNameMapKey(tPtr)
	typeTableNameMap[key] = tableName
}

func obtainTableNameMapKey(tPtr *reflect.Type) string {
	t := *tPtr
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return fmt.Sprintf("%s/%s", t.PkgPath(), t.Name())
}

func ObtainTableNameByEntityPtr(entityPtr interface{}) (tableName string) {
	t := reflect.TypeOf(entityPtr).Elem()
	tableName = ObtainTableNameByType(&t)
	return
}

func ObtainTableNameByType(tPtr *reflect.Type) (tableName string) {
	key := obtainTableNameMapKey(tPtr)
	name, exists := typeTableNameMap[key]
	if exists {
		return name
	}
	t := *tPtr
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	valObj := reflect.New(t)
	v, ok := valObj.Interface().(TableName)
	if ok {
		tableName = v.TableName()
		registerTableName(tPtr, v.TableName())
		return
	}
	tableName = strings.ToLower(t.Name())
	registerTableName(tPtr, tableName)
	return
}
