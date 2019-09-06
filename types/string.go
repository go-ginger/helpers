package types

import "fmt"

type str string

var String = new(str)

func (lv *str) Ptr(value interface{}) *string {
	if value == nil {
		return nil
	}
	v := fmt.Sprintf("%v", value)
	return &v
}
