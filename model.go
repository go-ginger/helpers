package helpers

import (
	gm "github.com/go-ginger/models"
	"reflect"
	"strings"
)

type iBeforeDump interface {
	BeforeDump(request gm.IRequest, data interface{})
}

func clear(value *reflect.Value, _type *reflect.Type) {
	if value == nil {
		return
	}
	if value.IsValid() {
		if value.CanSet() {
			value.Set(reflect.Zero(*_type))
		}
	}
}

func BeforeDump(request gm.IRequest, data interface{}) {
	s, ok := data.(reflect.Value)
	if !ok {
		s = reflect.ValueOf(data)
	}
	kind := s.Kind()
	if kind == reflect.Ptr {
		s = s.Elem()
		kind = s.Kind()
	}
	sType := s.Type()
	switch kind {
	case reflect.Slice:
		for ind := 0; ind < s.Len(); ind++ {
			BeforeDump(request, s.Index(ind))
		}
		break
	case reflect.Struct:
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			ff := sType.Field(i)
			if IsEmptyValue(f) {
				continue
			}
			switch f.Type().Kind() {
			case reflect.Ptr:
				if f.IsNil() {
					break
				}
				BeforeDump(request, f.Elem())
				break
			case reflect.Interface:
				i := f.Interface()
				BeforeDump(request, i)
				break
			case reflect.Struct:
				BeforeDump(request, f)
				break
			case reflect.Slice:
				for ind := 0; ind < f.Len(); ind++ {
					BeforeDump(request, f.Index(ind))
				}
				break
			}
			if cs, ok := ff.Tag.Lookup("c"); ok {
				continueCheck := true
				csParts := strings.Split(cs, ",")
				for _, csPart := range csParts {
					if csPart == "load_only" {
						continueCheck = false
						break
					}
				}
				if !continueCheck {
					continue
				}
			}
			tag, ok := ff.Tag.Lookup("read_roles")
			if ok {
				canRead := false
				auth := request.GetAuth()
				if auth != nil {
					tagParts := strings.Split(tag, ",")
					for _, role := range tagParts {
						if auth.HasRole(role) || (role == "id" &&
							auth.GetCurrentAccountId(request) == request.GetIDString()) {
							canRead = true
							break
						}
					}
				}
				if !canRead {
					clear(&f, &ff.Type)
				}
			}
		}
		if s.CanAddr() {
			addr := s.Addr()
			if addr.IsValid() && addr.CanInterface() {
				mv := addr.Interface()
				if baseModel, ok := mv.(gm.IBaseModel); ok {
					baseModel.Populate(request)
				}
				if cls, ok := mv.(iBeforeDump); ok {
					cls.BeforeDump(request, mv)
				}
			}
		}
		break
	}
}
