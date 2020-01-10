package helpers

import (
	"reflect"
)

// ReflectMethod Reflect if an interface is either a struct or a pointer to a struct
func ReflectMethod(obj interface{}, methodName string) *reflect.Value {
	value := reflect.ValueOf(obj)

	// Check if the passed interface is a pointer
	if value.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface, so we have a pointer to work with
		value = reflect.New(reflect.TypeOf(obj))
	}

	// Get the method by name
	Method := value.MethodByName(methodName)
	if !Method.IsValid() {
		return nil
	}
	return &Method
}

func ExecuteFunction(f interface{}, funcName string) interface{} {
	rf := reflect.TypeOf(f)
	if rf.Kind() != reflect.Func {
		return nil
	}
	vf := reflect.ValueOf(f)
	ReflectMethod(f, funcName).Call([]reflect.Value{reflect.ValueOf(1)})
	return vf.Interface()
}

type Value struct {
	reflect.Value
}

func CreateArray(t reflect.Type) reflect.Value {
	var arrayType reflect.Type
	arrayType = reflect.SliceOf(t)
	return reflect.Zero(arrayType)
}

// NewInstance returns new instance with fresh properties of given value
func NewInstance(value interface{}) interface{} {
	var result interface{}
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		result = reflect.New(rv.Elem().Type()).Interface()
	} else {
		result = reflect.New(reflect.TypeOf(value)).Elem().Interface()
	}
	return result
}

func NewInstanceOfType(typ reflect.Type) interface{} {
	var result interface{}
	result = reflect.New(typ).Interface()
	return result
}

func NewSliceInstanceOfType(typ reflect.Type) interface{} {
	var result interface{}
	result = reflect.MakeSlice(reflect.SliceOf(typ), 0, 0).Interface()
	return result
}

func NewSliceInstanceOfTypePtr(typ reflect.Type) interface{} {
	var result interface{}
	result = reflect.New(reflect.SliceOf(typ)).Interface()
	return result
}

func AppendToSlice(arr interface{}, valuePtr interface{}) (result interface{}) {
	arrValue := reflect.ValueOf(arr)
	result = reflect.Append(arrValue, reflect.ValueOf(valuePtr).Elem()).Interface()
	return
}

func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
