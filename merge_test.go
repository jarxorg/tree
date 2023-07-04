package tree

import (
	"reflect"
	"testing"
)

func TestMergeOption(t *testing.T) {
	tests := []struct {
		is   func() bool
		want bool
	}{
		{is: MergeOptionDefault.isOverrideMap, want: false},
		{is: MergeOptionDefault.isOverrideArray, want: false},
		{is: MergeOptionDefault.isOverrideValue, want: false},
		{is: MergeOptionDefault.isReplaceMap, want: false},
		{is: MergeOptionDefault.isReplaceArray, want: false},
		{is: MergeOptionDefault.isReplaceValue, want: false},
		{is: MergeOptionDefault.isAppend, want: false},
		{is: MergeOptionDefault.isSlurp, want: false},
		{is: MergeOptionOverrideMap.isOverrideMap, want: true},
		{is: MergeOptionOverrideMap.isOverrideArray, want: false},
		{is: MergeOptionOverrideMap.isOverrideValue, want: true},
		{is: MergeOptionOverrideMap.isReplaceMap, want: false},
		{is: MergeOptionOverrideMap.isReplaceArray, want: false},
		{is: MergeOptionOverrideMap.isReplaceValue, want: false},
		{is: MergeOptionOverrideMap.isAppend, want: false},
		{is: MergeOptionOverrideMap.isSlurp, want: false},
		{is: MergeOptionOverrideArray.isOverrideMap, want: false},
		{is: MergeOptionOverrideArray.isOverrideArray, want: true},
		{is: MergeOptionOverrideArray.isOverrideValue, want: true},
		{is: MergeOptionOverrideArray.isReplaceMap, want: false},
		{is: MergeOptionOverrideArray.isReplaceArray, want: false},
		{is: MergeOptionOverrideArray.isReplaceValue, want: false},
		{is: MergeOptionOverrideArray.isAppend, want: false},
		{is: MergeOptionOverrideArray.isSlurp, want: false},
		{is: MergeOptionOverride.isOverrideMap, want: true},
		{is: MergeOptionOverride.isOverrideArray, want: true},
		{is: MergeOptionOverride.isOverrideValue, want: true},
		{is: MergeOptionOverride.isReplaceMap, want: false},
		{is: MergeOptionOverride.isReplaceArray, want: false},
		{is: MergeOptionOverride.isReplaceValue, want: false},
		{is: MergeOptionOverride.isAppend, want: false},
		{is: MergeOptionOverride.isSlurp, want: false},
		{is: MergeOptionReplaceMap.isOverrideMap, want: false},
		{is: MergeOptionReplaceMap.isOverrideArray, want: false},
		{is: MergeOptionReplaceMap.isOverrideValue, want: false},
		{is: MergeOptionReplaceMap.isReplaceMap, want: true},
		{is: MergeOptionReplaceMap.isReplaceArray, want: false},
		{is: MergeOptionReplaceMap.isReplaceValue, want: true},
		{is: MergeOptionReplaceMap.isAppend, want: false},
		{is: MergeOptionReplaceMap.isSlurp, want: false},
		{is: MergeOptionReplaceArray.isOverrideMap, want: false},
		{is: MergeOptionReplaceArray.isOverrideArray, want: false},
		{is: MergeOptionReplaceArray.isOverrideValue, want: false},
		{is: MergeOptionReplaceArray.isReplaceMap, want: false},
		{is: MergeOptionReplaceArray.isReplaceArray, want: true},
		{is: MergeOptionReplaceArray.isReplaceValue, want: true},
		{is: MergeOptionReplaceArray.isAppend, want: false},
		{is: MergeOptionReplaceArray.isSlurp, want: false},
		{is: MergeOptionReplace.isOverrideMap, want: false},
		{is: MergeOptionReplace.isOverrideArray, want: false},
		{is: MergeOptionReplace.isOverrideValue, want: false},
		{is: MergeOptionReplace.isReplaceMap, want: true},
		{is: MergeOptionReplace.isReplaceArray, want: true},
		{is: MergeOptionReplace.isReplaceValue, want: true},
		{is: MergeOptionReplace.isAppend, want: false},
		{is: MergeOptionReplace.isSlurp, want: false},
		{is: MergeOptionAppend.isOverrideMap, want: false},
		{is: MergeOptionAppend.isOverrideArray, want: false},
		{is: MergeOptionAppend.isOverrideValue, want: false},
		{is: MergeOptionAppend.isReplaceMap, want: false},
		{is: MergeOptionAppend.isReplaceArray, want: false},
		{is: MergeOptionAppend.isReplaceValue, want: false},
		{is: MergeOptionAppend.isAppend, want: true},
		{is: MergeOptionAppend.isSlurp, want: false},
		{is: MergeOptionSlurp.isOverrideMap, want: false},
		{is: MergeOptionSlurp.isOverrideArray, want: false},
		{is: MergeOptionSlurp.isOverrideValue, want: false},
		{is: MergeOptionSlurp.isReplaceMap, want: false},
		{is: MergeOptionSlurp.isReplaceArray, want: false},
		{is: MergeOptionSlurp.isReplaceValue, want: false},
		{is: MergeOptionSlurp.isAppend, want: false},
		{is: MergeOptionSlurp.isSlurp, want: true},
		{is: (MergeOptionOverrideMap | MergeOptionAppend).isOverrideMap, want: true},
		{is: (MergeOptionOverrideMap | MergeOptionAppend).isOverrideArray, want: false},
		{is: (MergeOptionOverrideMap | MergeOptionAppend).isOverrideValue, want: true},
		{is: (MergeOptionOverrideMap | MergeOptionAppend).isReplaceMap, want: false},
		{is: (MergeOptionOverrideMap | MergeOptionAppend).isReplaceArray, want: false},
		{is: (MergeOptionOverrideMap | MergeOptionAppend).isReplaceValue, want: false},
		{is: (MergeOptionOverrideMap | MergeOptionAppend).isAppend, want: true},
		{is: (MergeOptionOverrideMap | MergeOptionAppend).isSlurp, want: false},
	}
	for i, test := range tests {
		if got := test.is(); got != test.want {
			t.Errorf("tests[%d] got %v; want %v", i, got, test.want)
		}
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		a    Node
		b    Node
		opts MergeOption
		want Node
	}{
		{
			a:    Map{"a": ToValue(1), "b": ToValue(2)},
			b:    Map{"a": ToValue(3), "c": ToValue(4)},
			want: Map{"a": ToValue(1), "b": ToValue(2), "c": ToValue(4)},
		}, {
			a:    Map{"a": ToValue(1), "b": ToValue(2)},
			b:    Nil,
			want: Map{"a": ToValue(1), "b": ToValue(2)},
		}, {
			a:    ToArrayValues(1, 2),
			b:    ToArrayValues(3, 4, 5),
			want: ToArrayValues(1, 2, 5),
		}, {
			a:    ToArrayValues(1, 2, 3),
			b:    ToArrayValues(4, 5),
			want: ToArrayValues(1, 2, 3),
		}, {
			a:    ToValue("a"),
			b:    ToValue("b"),
			want: ToValue("a"),
		}, {
			a:    Map{"a": ToValue(1), "b": ToValue(2)},
			b:    ToValue("c"),
			want: Map{"a": ToValue(1), "b": ToValue(2)},
		}, {
			a:    ToArrayValues("a"),
			b:    ToValue("b"),
			want: ToArrayValues("a"),
		}, {
			a:    Map{"a": ToValue(1), "b": ToValue(2)},
			b:    Map{"a": ToValue(3)},
			opts: MergeOptionOverrideMap,
			want: Map{"a": ToValue(3), "b": ToValue(2)},
		}, {
			a:    Map{"a": ToValue(1), "b": ToValue(2)},
			b:    Nil,
			opts: MergeOptionOverrideMap,
			want: Nil,
		}, {
			a:    ToValue("a"),
			b:    ToValue("b"),
			opts: MergeOptionOverrideMap,
			want: ToValue("b"),
		}, {
			a:    Map{"a": ToValue(1), "b": ToValue(2), "c": ToValue(3)},
			b:    Map{"a": ToValue(4), "b": ToArrayValues(5, 6), "d": ToValue(7)},
			opts: MergeOptionOverrideMap,
			want: Map{"a": ToValue(4), "b": ToArrayValues(5, 6), "c": ToValue(3), "d": ToValue(7)},
		}, {
			a:    ToArrayValues(1, 2, 3),
			b:    ToArrayValues(4, 5),
			opts: MergeOptionOverrideArray,
			want: ToArrayValues(4, 5, 3),
		}, {
			a:    ToArrayValues(1, 2, 3),
			b:    ToArrayValues(4, 5, 6, 7),
			opts: MergeOptionOverrideArray,
			want: ToArrayValues(4, 5, 6, 7),
		}, {
			a:    ToValue("a"),
			b:    ToValue("b"),
			opts: MergeOptionOverrideArray,
			want: ToValue("b"),
		}, {
			a:    Map{"a": ToValue(1), "b": ToValue(2)},
			b:    Map{"a": ToValue(3)},
			opts: MergeOptionOverride,
			want: Map{"a": ToValue(3), "b": ToValue(2)},
		}, {
			a:    ToArrayValues(1, 2, 3),
			b:    ToArrayValues(4, 5),
			opts: MergeOptionOverride,
			want: ToArrayValues(4, 5, 3),
		}, {
			a:    ToValue("a"),
			b:    ToValue("b"),
			opts: MergeOptionOverride,
			want: ToValue("b"),
		}, {
			a:    Map{"a": ToValue(1), "b": ToValue(2)},
			b:    Map{"a": ToValue(3)},
			opts: MergeOptionReplaceMap,
			want: Map{"a": ToValue(3)},
		}, {
			a:    ToValue("a"),
			b:    ToValue("b"),
			opts: MergeOptionReplaceMap,
			want: ToValue("b"),
		}, {
			a:    ToValue("a"),
			b:    ToValue("b"),
			opts: MergeOptionReplaceMap,
			want: ToValue("b"),
		}, {
			a:    ToArrayValues(1, 2, 3),
			b:    ToArrayValues(4, 5),
			opts: MergeOptionReplaceArray,
			want: ToArrayValues(4, 5),
		}, {
			a:    ToValue("a"),
			b:    ToValue("b"),
			opts: MergeOptionReplaceArray,
			want: ToValue("b"),
		}, {
			a:    ToValue("a"),
			b:    ToValue("b"),
			opts: MergeOptionReplace,
			want: ToValue("b"),
		}, {
			a:    ToArrayValues(1, 2, 3),
			b:    ToArrayValues(4, 5),
			opts: MergeOptionReplace,
			want: ToArrayValues(4, 5),
		}, {
			a:    ToValue("a"),
			b:    ToValue("b"),
			opts: MergeOptionReplace,
			want: ToValue("b"),
		}, {
			a:    ToArrayValues(1, 2, 3),
			b:    ToArrayValues(4, 5),
			opts: MergeOptionAppend,
			want: ToArrayValues(1, 2, 3, 4, 5),
		}, {
			a:    ToArrayValues(1, 2, 3),
			b:    ToArrayValues(4, 5),
			opts: MergeOptionSlurp,
			want: ToArrayValues(1, 2, 3, 4, 5),
		}, {
			a:    ToArrayValues(1, 2, 3),
			b:    ToValue(4),
			opts: MergeOptionSlurp,
			want: ToArrayValues(1, 2, 3, 4),
		}, {
			a:    ToValue(1),
			b:    ToValue(2),
			opts: MergeOptionSlurp,
			want: ToArrayValues(1, 2),
		}, {
			a:    Map{"a": ToValue(1)},
			b:    Map{"a": ToValue(2)},
			opts: MergeOptionSlurp,
			want: Map{"a": ToArrayValues(1, 2)},
		}, {
			a: Map{
				"map":          Map{"a": ToValue(1), "b": ToValue(2)},
				"array":        ToArrayValues(3, 4, 5),
				"arrayOrValue": ToArrayValues(6, 7, 8),
			},
			b: Map{
				"map":          Map{"a": ToValue(9)},
				"array":        ToArrayValues(10, 11),
				"arrayOrValue": ToValue(12),
			},
			opts: MergeOptionOverrideMap | MergeOptionAppend,
			want: Map{
				"map":          Map{"a": ToValue(9), "b": ToValue(2)},
				"array":        ToArrayValues(3, 4, 5, 10, 11),
				"arrayOrValue": ToValue(12),
			},
		}, {
			a: Map{
				"map":          Map{"a": ToValue(1), "b": ToValue(2)},
				"array":        ToArrayValues(3, 4, 5),
				"arrayOrValue": ToArrayValues(6, 7, 8),
			},
			b: Map{
				"map":          Map{"a": ToValue(9)},
				"array":        ToArrayValues(10, 11),
				"arrayOrValue": ToValue(12),
			},
			opts: MergeOptionAppend | MergeOptionSlurp,
			want: Map{
				"map":          Map{"a": ToArrayValues(1, 9), "b": ToValue(2)},
				"array":        ToArrayValues(3, 4, 5, 10, 11),
				"arrayOrValue": ToArrayValues(6, 7, 8, 12),
			},
		}, {
			a: Map{
				"map":   Map{"a": ToValue(1), "b": ToValue(2)},
				"array": ToArrayValues(3, 4, 5),
			},
			b: Map{
				"map":   Map{"a": ToValue(6)},
				"array": ToArrayValues(7, 8),
			},
			opts: MergeOptionReplaceMap | MergeOptionAppend,
			want: Map{
				"map":   Map{"a": ToValue(6)},
				"array": ToArrayValues(7, 8),
			},
		},
	}
	for i, test := range tests {
		a := Clone(test.a)
		b := Clone(test.b)
		got := Merge(a, b, test.opts)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(`tests[%d]: unexpected %v; want %v`, i, got, test.want)
		}
		if !reflect.DeepEqual(a, test.a) {
			t.Errorf(`tests[%d]: unexpected %v; want %v`, i, got, test.want)
		}
		if !reflect.DeepEqual(b, test.b) {
			t.Errorf(`tests[%d]: unexpected %v; want %v`, i, got, test.want)
		}
	}
}
