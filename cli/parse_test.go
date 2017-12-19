package cli

import (
	"bytes"
	"log"
	"testing"

	"gnorm.org/gnorm/environ"
	"gnorm.org/gnorm/run/data"

	"github.com/BurntSushi/toml"
	"github.com/google/go-cmp/cmp"
)

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

func TestParseConfig(t *testing.T) {
	var stderr, stdout bytes.Buffer
	env := environ.Values{
		Stderr: &stderr,
		Stdout: &stdout,
		Log:    log.New(&stderr, "", 0),
	}
	cfg, err := parseFile(env, "gnorm.toml")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(cfg.Params, map[string]interface{}{"mySpecialValue": "some value"}); diff != "" {
		t.Errorf("Params not copied correctly:\n%s", diff)
	}
	expected := data.ConfigData{
		ConnStr: "dbname=mydb host=127.0.0.1 sslmode=disable user=admin",
		DBType:  "postgres",
		Schemas: []string{"public"},
		PostRun: []string{"echo", "$GNORMFILE"},
		IncludeTables: map[string][]string{
			"public": nil,
		},
		ExcludeTables: map[string][]string{
			"public": []string{"xyzzx"},
		},
		TypeMap: map[string]string{
			"timestamp with time zone": "time.Time",
			"text":              "string",
			"boolean":           "bool",
			"uuid":              "uuid.UUID",
			"character varying": "string",
			"integer":           "int",
			"numeric":           "float64",
		},
		NullableTypeMap: map[string]string{
			"timestamp with time zone": "pq.NullTime",
			"text":              "sql.NullString",
			"boolean":           "sql.NullBool",
			"uuid":              "uuid.NullUUID",
			"character varying": "sql.NullString",
			"integer":           "sql.NullInt64",
			"numeric":           "sql.NullFloat64",
		},
		PluginDirs:       []string{"plugins"},
		OutputDir:        "gnorm",
		StaticDir:        "static",
		NoOverwriteGlobs: []string{"*.perm.go"},
	}
	if diff := cmp.Diff(cfg.ConfigData, expected); diff != "" {
		t.Fatalf("Actual differs from expected:\n%s", diff)
	}

}

func TestParseGnormToml(t *testing.T) {
	c := Config{}
	m, err := toml.DecodeFile("gnorm.toml", &c)
	if err != nil {
		t.Fatal(err)
	}
	undec := m.Undecoded()
	if len(undec) > 0 {
		t.Fatalf("unknown values present in config file: %s", undec)
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
