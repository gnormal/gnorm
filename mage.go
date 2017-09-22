//+build mage

// This is the "magefile" for gnorm.  To install mage, run go get github.com/magefile/mage.
// To build gnorm, just mage build.

package main

import (
	"log"
	"os"
)

// Runs go install for gnorm.  This generates the embedded docs and the version
// info into the binary.
func Build() error {
	Deps()
	cleanup := genSite()
	defer cleanup()

	log.Print("running go install")
	return run("go", "install", "-tags", "make", "--ldflags="+flags(), "gnorm.org/gnorm")
}

// Generates binaries for all supported versions.  Currently that means a
// combination of windows, linux, and OSX in 32 bit and 64 bit formats.  The
// files will be dumped in the local directory with names according to their
// supported platform.
func All() error {
	Deps()
	cleanup := genSite()
	defer cleanup()

	ldf := flags()
	for _, OS := range []string{"windows", "darwin", "linux"} {
		if err := os.Setenv("GOOS", OS); err != nil {
			return err
		}
		for _, ARCH := range []string{"amd64", "386"} {
			if err := os.Setenv("GOOS", OS); err != nil {
				return err
			}
			log.Printf("running go build for GOOS=%s GOARCH=%s", OS, ARCH)
			err := run("go", "build", "-tags", "make", "-o", "gnorm_"+OS+"_"+ARCH, "--ldflags="+ldf)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Downloads commands needed to build gnorm.
func Deps() {
	log.Print("downloading hugo")
	must(run("go", "get", "github.com/gohugoio/hugo"))
	log.Print("downloading statik")
	must(run("go", "get", "github.com/rakyll/statik"))
}
