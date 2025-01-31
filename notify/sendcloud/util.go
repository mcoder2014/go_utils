package sendcloud

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// StructToURLValues 将结构体指针转换为 url.Values
func StructToURLValues(obj any) (*url.Values, error) {
	// 检查传入的是否为指针
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return nil, fmt.Errorf("input must be a non-nil pointer")
	}
	// 获取指针指向的值
	value = value.Elem()
	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must point to a struct")
	}

	values := url.Values{}
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if !field.CanInterface() {
			continue
		}
		// 获取字段的 json 标签
		jsonTag := typ.Field(i).Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		// 处理 json 标签中的选项
		name := strings.Split(jsonTag, ",")[0]
		if name == "-" {
			continue
		}
		omitempty := strings.Contains(jsonTag, "omitempty")

		// 根据字段类型处理值
		switch field.Kind() {
		case reflect.String:
			if omitempty && field.String() == "" {
				continue
			}
			values.Add(name, field.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if omitempty && field.Int() == 0 {
				continue
			}
			values.Add(name, strconv.FormatInt(field.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if omitempty && field.Uint() == 0 {
				continue
			}
			values.Add(name, strconv.FormatUint(field.Uint(), 10))
		case reflect.Bool:
			if omitempty && !field.Bool() {
				continue
			}
			values.Add(name, strconv.FormatBool(field.Bool()))
		default:
			// 可以根据需要添加更多类型的处理逻辑
			continue
		}
	}
	return &values, nil
}
