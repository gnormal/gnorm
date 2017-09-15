package database

import "testing"

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
