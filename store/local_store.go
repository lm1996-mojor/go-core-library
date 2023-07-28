package store

import (
	"sync"
)

var storeMap sync.Map
var PoInterKey = MapPoInterKey{}

type MapPoInterKey struct {
	UserId   string
	ClientId string
}

// Set store value
func Set(key string, value interface{}) {
	//tls.Set(key, tls.MakeData(value))
	storeMap.Store(key, value)
	return
}

// Get value by key
func Get(key string) (value interface{}, ok bool) {
	//d, ok := tls.Get(key)
	//if ok {
	//	value = d.Value()
	//}
	value, ok = storeMap.Load(key)
	return
}

// Del delete value by key
func Del(key interface{}) {
	//tls.Del(key)
	storeMap.Delete(key)
}

// Clean empty local map
func Clean() {
	storeMap.Range(walk)
}

// 删除并带检测
func walk(key, value interface{}) bool {
	storeMap.Delete(key)
	_, ok := storeMap.Load(key)
	return !ok
}
