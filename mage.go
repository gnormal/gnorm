//+build mage

// This is the "magefile" for gnorm.  To install mage, run go get github.com/magefile/mage.
// To build gnorm, just mage build.

package main

import (
	"log"
)

// Runs go install for gnorm.  This generates the embedded docs and the version
// info into the binary.
func Build() error {
	if err := genSite(); err != nil {
		return err
	}
	defer cleanup()

	ldf, err := flags()
	if err != nil {
		return err
	}

	log.Print("running go install")
	// use -tags make so we can have different behavior for when we know we've built with mage.
	return run("go", "install", "-tags", "make", "--ldflags="+ldf, "gnorm.org/gnorm")
}

// Generates binaries for all supported versions.  Currently that means a
// combination of windows, linux, and OSX in 32 bit and 64 bit formats.  The
// files will be dumped in the local directory with names according to their
// supported platform.
func All() error {
	if err := genSite(); err != nil {
		return err
	}
	defer cleanup()

	ldf, err := flags()
	if err != nil {
		return err
	}
	for _, OS := range []string{"windows", "darwin", "linux"} {
		for _, ARCH := range []string{"amd64", "386"} {
			log.Printf("running go build for GOOS=%s GOARCH=%s", OS, ARCH)
			env := []string{"GOOS=" + OS, "GOARCH=" + ARCH}
			if err := runWith(env, "go", "build", "-tags", "make", "-o", "gnorm_"+OS+"_"+ARCH, "--ldflags="+ldf); err != nil {
				return err
			}
		}
	}
	return err
}

// Removes generated cruft.  This target shouldn't ever be necessary, since the
// cleanup should happen automatically, but it's here just in case.
func Clean() error {
	return cleanup()
}
