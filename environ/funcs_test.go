package environ

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"testing"
	"text/template"
)

const testEnv = "GO_TEST_ENV"

func TestPlugin(t *testing.T) {
	table := []struct {
		cmd, function string
		ctx           interface{}
		pass          bool
		expect        map[string]interface{}
	}{
		{"nix", "echo", "hello,world", true,
			map[string]interface{}{
				"data": "hello,world",
			},
		},
	}

	for _, v := range table {
		i, err := toJSON(v.ctx)
		if err != nil {
			t.Fatal(err)
		}
		o, err := execJSON(testRunner, v.cmd, v.function, string(i))
		if err != nil {
			t.Fatal(err)
		}
		for k, e := range o {
			ev := v.expect[k]
			if !reflect.DeepEqual(e, ev) {
				t.Errorf("expected %v got %v", ev, e)
			}
		}
	}

	tpl, err := template.New("plugin").Funcs(
		template.FuncMap{
			"plugin": func(name, function string, ctx interface{}) (interface{}, error) {
				return callPlugin(testRunner, name, function, ctx)
			},
		},
	).Parse(`{{range plugin "nix" "echoPlugin" . }}{{.}}{{end}}`)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err = tpl.Execute(&buf, "Hello,World"); err != nil {
		t.Fatal(err)
	}
	expect := "Hello,Worldnix echoPlugin"
	got := buf.String()
	if got != expect {
		t.Errorf("expected %s got %s", expect, got)
	}
}

func toJSON(ctx interface{}) ([]byte, error) {
	d := make(map[string]interface{})
	d["data"] = ctx
	return json.Marshal(d)
}

func testRunner(name string, args ...string) ([]byte, error) {
	v := []string{name}
	return helperCMD(append(v, args...)...).CombinedOutput()
}

func helperCMD(args ...string) *exec.Cmd {
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = []string{"GO_TEST_ENV=command"}
	return cmd
}

func TestMain(t *testing.M) {
	switch os.Getenv(testEnv) {
	case "command":
		args := os.Args[1:]
		switch args[0] {
		case "nix":
			switch args[1] {
			case "echo":
				fmt.Println(args[2])
			case "echoPlugin":
				c := args[2]
				d := make(map[string]interface{})
				if err := json.Unmarshal([]byte(c), &d); err != nil {
					log.Fatal(err)
				}
				data := d["data"].(string)
				d["data"] = data + "nix echoPlugin"
				v, err := toJSON(d)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(v))
			}
		}
	default:
		os.Exit(t.Run())
	}
}
