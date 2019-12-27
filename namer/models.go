package namer

import "reflect"

type INamer interface {
	Initialize()
	GetNameByType(reflectType reflect.Type) string
	GetName(value interface{}) string
}
