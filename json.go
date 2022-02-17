package tree

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// MarshalJSON returns the JSON encoding of the specified node.
func MarshalJSON(n Node) ([]byte, error) {
	return json.Marshal(n)
}

// DecodeJSON decodes JSON as a node using the provided decoder.
func DecodeJSON(dec *json.Decoder) (Node, error) {
	t, err := dec.Token()
	if err != nil || t == nil {
		return nil, err
	}
	switch t.(type) {
	case string:
		return StringValue(t.(string)), nil
	case float64:
		return NumberValue(t.(float64)), nil
	case bool:
		return BoolValue(t.(bool)), nil
	case json.Delim:
		d := t.(json.Delim)
		switch ds := d.String(); ds {
		case "{":
			m := Map{}
			if err := jsonMap(dec, &m); err != nil {
				return nil, err
			}
			return m, nil
		case "[":
			return jsonArray(dec, Array{})
		}
	}
	return nil, fmt.Errorf("Unknown token %#v", t)
}

// UnmarshalJSON parses the JSON-encoded data to a Node.
func UnmarshalJSON(data []byte) (Node, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	return DecodeJSON(dec)
}

// UnmarshalJSON is an implementation of json.Unmarshaler.
func (n *Map) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	t, err := dec.Token()
	if err != nil {
		return err
	}
	d, ok := t.(json.Delim)
	if !ok || d.String() != "{" {
		return fmt.Errorf("Unknown token %#v", t)
	}
	if *n == nil {
		*n = make(Map)
	}
	return jsonMap(dec, n)
}

// UnmarshalJSON is an implementation of json.Unmarshaler.
func (n *Array) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	t, err := dec.Token()
	if err != nil {
		return err
	}
	d, ok := t.(json.Delim)
	if !ok || d.String() != "[" {
		return fmt.Errorf("Unknown token %#v", t)
	}
	*n, err = jsonArray(dec, *n)
	return err
}

func jsonMap(dec *json.Decoder, m *Map) error {
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if d, ok := t.(json.Delim); ok {
		if d.String() == "}" {
			return nil
		}
		return fmt.Errorf("Unknown token %#v", t)
	}

	key, ok := t.(string)
	if !ok {
		return fmt.Errorf("Unknown token %#v", t)
	}

	t, err = dec.Token()
	if err != nil {
		return err
	}
	if d, ok := t.(json.Delim); ok {
		switch ds := d.String(); ds {
		case "{":
			mm := Map{}
			if err := jsonMap(dec, &mm); err != nil {
				return err
			}
			(*m)[key] = mm
			return jsonMap(dec, m)
		case "[":
			aa, err := jsonArray(dec, Array{})
			if err != nil {
				return err
			}
			(*m)[key] = aa
			return jsonMap(dec, m)
		}
	}

	(*m)[key] = jsonValue(t)
	return jsonMap(dec, m)
}

func jsonArray(dec *json.Decoder, a Array) (Array, error) {
	t, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if d, ok := t.(json.Delim); ok {
		switch ds := d.String(); ds {
		case "]":
			return a, nil
		case "{":
			mm := Map{}
			if err := jsonMap(dec, &mm); err != nil {
				return nil, err
			}
			return jsonArray(dec, append(a, mm))
		case "[":
			aa, err := jsonArray(dec, Array{})
			if err != nil {
				return nil, err
			}
			return jsonArray(dec, append(a, aa))
		}
	}
	return jsonArray(dec, append(a, jsonValue(t)))
}

func jsonValue(t json.Token) Node {
	if t == nil {
		return nil
	}
	switch t.(type) {
	case string:
		return StringValue(t.(string))
	case bool:
		return BoolValue(t.(bool))
	case float64:
		return NumberValue(t.(float64))
	}
	return StringValue(fmt.Sprintf("%#v", t))
}
