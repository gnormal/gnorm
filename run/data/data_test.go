package data

import (
	"reflect"
	"testing"
)

func TestStringsSprintf(t *testing.T) {
	s := Strings{"one", "two", "three"}

	s = s.Sprintf("%s 1")

	expected := Strings{"one 1", "two 1", "three 1"}
	for x := range expected {
		if s[x] != expected[x] {
			t.Fatalf("index %v of returned strings expected %q but got %q", x, expected[x], s[x])
		}
	}
}

func TestStringsExcept(t *testing.T) {
	vals := Strings{"one", "two", "three", "four"}.Except([]string{"ONE", "two", "three"})
	expected := Strings{"one", "four"}
	if !reflect.DeepEqual(vals, expected) {
		t.Fatalf("expected %q, but got %q", expected, vals)
	}
}

func TestStringsSort(t *testing.T) {
	vals := Strings{"one", "two", "three", "four"}.Sorted()
	expected := Strings{"four", "one", "three", "two"}
	if !reflect.DeepEqual(vals, expected) {
		t.Fatalf("expected %q, but got %q", expected, vals)
	}
}
