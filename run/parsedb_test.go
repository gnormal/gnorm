package run

import "testing"

func TestMakeFilter(t *testing.T) {
	var include, exclude map[string][]string
	f := makeFilter(include, exclude)
	if !f("anything", "any other thing") {
		t.Fatalf("nil maps should return a filter that doesn't filter")
	}
	include = map[string][]string{}
	exclude = map[string][]string{}
	f = makeFilter(include, exclude)
	if !f("anything", "any other thing") {
		t.Fatalf("empty maps should return a filter that doesn't filter")
	}
}
