package tree

import (
	"gopkg.in/yaml.v2"
)

// MarshalYAML returns the YAML encoding of the specified node.
func MarshalYAML(n Node) ([]byte, error) {
	return yaml.Marshal(n)
}

// UnmarshalYAML returns the YAML encoding of the specified node.
func UnmarshalYAML(data []byte) (Node, error) {
	var v interface{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return ToNode(v), nil
}

// UnmarshalYAML is an implementation of yaml.Unmarshaler.
func (n *Map) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v interface{}
	if err := unmarshal(&v); err != nil {
		return err
	}
	if *n == nil {
		*n = make(Map)
	}
	for k, v := range ToNode(v).Map() {
		(*n)[k] = v
	}
	return nil
}

// UnmarshalYAML is an implementation of yaml.Unmarshaler.
func (n *Array) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v interface{}
	if err := unmarshal(&v); err != nil {
		return err
	}
	for _, v := range ToNode(v).Array() {
		*n = append(*n, v)
	}
	return nil
}
