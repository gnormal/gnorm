package run

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMakeTable(t *testing.T) {
	people := []struct {
		Name string
		Age  int
	}{
		{
			Name: "Bob",
			Age:  30,
		},
		{
			Name: "Samantha",
			Age:  3,
		},
	}

	s, err := makeTable(people, "{{.Name}}|{{.Age}}", "Name", "Age")
	if err != nil {
		t.Fatal(err)
	}
	expected := `
+----------+-----+
|   Name   | Age |
+----------+-----+
| Bob      |  30 |
| Samantha |   3 |
+----------+-----+
`[1:]
	if s != expected {
		t.Fatal(cmp.Diff(s, expected))
	}
}
