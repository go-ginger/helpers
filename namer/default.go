package namer

import (
	"bytes"
	"github.com/jinzhu/inflection"
	"reflect"
	"strings"
)

type Default struct {
	INamer

	sMap *safeMap
}

func (n *Default) Initialize() {
	n.sMap = newSafeMap()
}

func (n *Default) GetNameByType(reflectType reflect.Type) string {
	for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	return inflection.Plural(n.getName(reflectType.Name()))
}

func (n *Default) GetName(value interface{}) string {
	reflectType := reflect.ValueOf(value).Type()
	return n.GetNameByType(reflectType)
}

func (n *Default) getName(name string) string {
	const (
		lower = false
		upper = true
	)

	if v := n.sMap.Get(name); v != "" {
		return v
	}

	if name == "" {
		return ""
	}

	replacer := GetReplacer()
	var (
		value                                    = replacer.Replace(name)
		buf                                      = bytes.NewBufferString("")
		lastCase, currCase, nextCase, nextNumber bool
	)

	for i, v := range value[:len(value)-1] {
		nextCase = value[i+1] >= 'A' && value[i+1] <= 'Z'
		nextNumber = value[i+1] >= '0' && value[i+1] <= '9'

		if i > 0 {
			if currCase == upper {
				if lastCase == upper && (nextCase == upper || nextNumber == upper) {
					buf.WriteRune(v)
				} else {
					if value[i-1] != '_' && value[i+1] != '_' {
						buf.WriteRune('_')
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
				if i == len(value)-2 && (nextCase == upper && nextNumber == lower) {
					buf.WriteRune('_')
				}
			}
		} else {
			currCase = upper
			buf.WriteRune(v)
		}
		lastCase = currCase
		currCase = nextCase
	}

	buf.WriteByte(value[len(value)-1])

	s := strings.ToLower(buf.String())
	n.sMap.Set(name, s)
	return s
}
