package tree

import (
	"errors"
	"fmt"
	"regexp"
	"sync"
)

// ToValue converts the specified v to a Value as Node.
// Node.Value() returns converted value.
func ToValue(v interface{}) Node {
	if v == nil {
		return nil
	}
	switch tv := v.(type) {
	case string:
		return StringValue(tv)
	case bool:
		return BoolValue(tv)
	case int:
		return NumberValue(int64(tv))
	case int64:
		return NumberValue(tv)
	case int32:
		return NumberValue(int64(tv))
	case float64:
		return NumberValue(tv)
	case float32:
		return NumberValue(float64(tv))
	case uint64:
		return NumberValue(float64(tv))
	case uint32:
		return NumberValue(float64(tv))
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
	switch tv := v.(type) {
	case Node:
		return tv
	case []interface{}:
		a := make(Array, len(tv))
		for i, vv := range tv {
			a[i] = ToNode(vv)
		}
		return a
	case map[string]interface{}:
		m := Map{}
		for k := range tv {
			m[k] = ToNode(tv[k])
		}
		return m
	case map[interface{}]interface{}:
		m := Map{}
		for k := range tv {
			m[fmt.Sprintf("%v", k)] = ToNode(tv[k])
		}
		return m
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

	return n.Each(func(key interface{}, v Node) error {
		if key == nil {
			return nil
		}
		keys[last] = key
		return walk(v, keys, fn)
	})
}

var regexpPool = sync.Pool{
	New: func() interface{} {
		return map[string]*regexp.Regexp{}
	},
}

func pooledRegexp(expr string) (*regexp.Regexp, error) {
	cache := regexpPool.Get().(map[string]*regexp.Regexp)
	defer regexpPool.Put(cache)

	if re, ok := cache[expr]; ok {
		return re, nil
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	cache[expr] = re
	return re, nil
}

func regexpMatchString(expr, value string) (bool, error) {
	re, err := pooledRegexp(expr)
	if err != nil {
		return false, err
	}
	return re.MatchString(value), nil
}
