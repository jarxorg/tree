package tree

import (
	"errors"
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
	a := make(Array, len(vs))
	for i, v := range vs {
		a[i] = ToValue(v)
	}
	return a
}

// ToNodeValues calss ToValues for each provided vs and returns them as []Node.
func ToNodeValues(vs ...interface{}) []Node {
	ns := make([]Node, len(vs))
	for i, v := range vs {
		ns[i] = ToValue(v)
	}
	return ns
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
		aa := make(Array, len(a))
		for i, vv := range a {
			aa[i] = ToNode(vv)
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

// SkipWalk is used as a return value from WalkFunc to indicate that
// the node and that children in the call is to be skipped.
// It is not returned as an error by any function.
var SkipWalk = errors.New("skip")

// WalkFunc is the type of the function called by Walk to visit each nodes.
//
// The keys argument contains that parent keys and the node key that
// type is int (array index) or string (map key).
type WalkFunc func(n Node, keys []interface{}) error

// Walk walks the node tree rooted at root, calling fn for each node or
// that children in the tree, including root.
func Walk(n Node, fn WalkFunc) error {
	return walk(n, []interface{}{}, fn)
}

func walk(n Node, lastKeys []interface{}, fn WalkFunc) error {
	if err := fn(n, lastKeys); err != nil {
		if err == SkipWalk {
			return nil
		}
		return err
	}

	last := len(lastKeys)
	keys := make([]interface{}, last+1)
	copy(keys, lastKeys)

	if a := n.Array(); a != nil {
		for i, v := range a {
			keys[last] = i
			if err := walk(v, keys, fn); err != nil {
				return err
			}
		}
		return nil
	}
	if m := n.Map(); m != nil {
		for _, k := range m.Keys() {
			keys[last] = k
			if err := walk(m[k], keys, fn); err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}
