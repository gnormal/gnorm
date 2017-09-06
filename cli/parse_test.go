package cli

import "testing"

func TestParseTables(t *testing.T) {
	tables := []string{"table", "schema.table2"}
	schemas := []string{"schema", "schema2"}

	m, err := parseTables(tables, schemas)
	if err != nil {
		t.Fatal(err)
	}

	if len(m) != 2 {
		t.Fatalf("Expected returned map to have one entry per schema (2), but got %v", len(m))
	}

	list := m["schema"]
	if len(list) != 2 {
		t.Errorf(`expcted "schema" table list to have 2 items but has %v`, len(list))
	}
	if !contains(list, "table") {
		t.Error(`"schema" table list should have included schema-unspecified "table", but did not`)
	}
	if !contains(list, "table2") {
		t.Error(`"schema" table list should have included schema-unspecified "table", but did not`)
	}

	list2 := m["schema2"]
	if len(list2) != 1 {
		t.Errorf(`expcted "schema2" table list to have 1 item but has %v`, len(list2))
	}
	if !contains(list2, "table") {
		t.Error(`"schema" table list should have included schema-unspecified "table", but did not`)
	}
}

func contains(list []string, s string) bool {
	for x := range list {
		if list[x] == s {
			return true
		}
	}
	return false
}
