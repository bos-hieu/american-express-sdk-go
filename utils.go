package americanexpress

import (
	"net/url"
	"reflect"
	"strconv"
)

// encodeQuery converts a struct to URL query values
func encodeQuery(v interface{}) (url.Values, error) {
	values := url.Values{}
	
	if v == nil {
		return values, nil
	}
	
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	if val.Kind() != reflect.Struct {
		return values, nil
	}
	
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		
		// Get the tag value
		tag := fieldType.Tag.Get("url")
		if tag == "" || tag == "-" {
			continue
		}
		
		// Skip empty values
		if field.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}
		
		// Get the actual value
		var value string
		switch field.Kind() {
		case reflect.String:
			if field.String() != "" {
				value = field.String()
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Int() != 0 {
				value = strconv.FormatInt(field.Int(), 10)
			}
		case reflect.Bool:
			value = strconv.FormatBool(field.Bool())
		case reflect.Float32, reflect.Float64:
			if field.Float() != 0 {
				value = strconv.FormatFloat(field.Float(), 'f', -1, 64)
			}
		}
		
		if value != "" {
			values.Add(tag, value)
		}
	}
	
	return values, nil
}