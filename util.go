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

// ToArrayValues calss ToValues for each provided vs and returns them as an Array.
func ToArrayValues(vs ...interface{}) Array {
	var a Array
	for _, v := range vs {
		a = append(a, ToValue(v))
	}
	return a
}

// ToNode converts the specified v to an Node.
func ToNode(v interface{}) Node {
	if v == nil {
		return nil
	}
	switch v.(type) {
	case Node:
		return v.(Node)
	case []interface{}:
		a := v.([]interface{})
		aa := Array{}
		for _, vv := range a {
			aa = append(aa, ToNode(vv))
		}
		return aa
	case map[string]interface{}:
		m := v.(map[string]interface{})
		mm := Map{}
		for k := range m {
			mm[k] = ToNode(m[k])
		}
		return mm
	case map[interface{}]interface{}:
		m := v.(map[interface{}]interface{})
		mm := Map{}
		for k := range m {
			mm[fmt.Sprintf("%v", k)] = ToNode(m[k])
		}
		return mm
	}
	return ToValue(v)
}
