package tree

import (
	"fmt"
)

func ToValue(v interface{}) Node {
	if v == nil {
		return nil
	}
	switch v.(type) {
	case string:
		return StringValue(v.(string))
	case bool:
		return BoolValue(v.(bool))
	case int:
		return Int64Value(int64(v.(int)))
	case int64:
		return Int64Value(v.(int64))
	case int32:
		return Int64Value(int64(v.(int32)))
	case float64:
		return Float64Value(v.(float64))
	case float32:
		return Float64Value(float64(v.(float32)))
	case uint64:
		return Float64Value(float64(v.(uint64)))
	case uint32:
		return Float64Value(float64(v.(uint32)))
	}
	return StringValue(fmt.Sprintf("%v", v))
}

func ToArray(v ...interface{}) Array {
	var a Array
	for _, vv := range v {
		a = append(a, ToValue(vv))
	}
	return a
}
