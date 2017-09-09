package cli

import (
	"io/ioutil"
	"testing"
)

func TestSampleEmbedded(t *testing.T) {
	b, err := ioutil.ReadFile("gnorm.toml")
	if err != nil {
		t.Fatal("Unexpected error reading gnorm.toml file:", err)
	}
	if sample != string(b) {
		t.Fatal("gnorm.toml file on disk differs from embedded contents in sample.go. Did you run go generate?")
	}
}
