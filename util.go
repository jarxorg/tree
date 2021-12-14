package tree

import (
	"fmt"
)

// ToValue converts the specified v to a Value as Node.
// Node.Value() returns converted value.
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
		return NumberValue(int64(v.(int)))
	case int64:
		return NumberValue(v.(int64))
	case int32:
		return NumberValue(int64(v.(int32)))
	case float64:
		return NumberValue(v.(float64))
	case float32:
		return NumberValue(float64(v.(float32)))
	case uint64:
		return NumberValue(float64(v.(uint64)))
	case uint32:
		return NumberValue(float64(v.(uint32)))
	case Node:
		return v.(Node)
	}
	// NOTE: Unsupported type.
	return StringValue(fmt.Sprintf("%#v", v))
}

// ToArray converts the specified v to an Array.
func ToArray(v ...interface{}) Array {
	var a Array
	for _, vv := range v {
		a = append(a, ToValue(vv))
	}
	return a
}
