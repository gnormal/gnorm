package cli

import (
	"testing"

	"github.com/BurntSushi/toml"
)

func TestExampleToml(t *testing.T) {
	c := Config{}
	m, err := toml.DecodeFile("gnorm.toml", &c)
	if err != nil {
		t.Fatal("error parsing config file", err)
	}
	undec := m.Undecoded()
	if len(undec) > 0 {
		t.Fatal("Warning: unknown values present in config file:", undec)
	}
}
