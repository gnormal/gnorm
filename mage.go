//+build mage

// This is the "magefile" for gnorm.  To install mage, run go get github.com/magefile/mage.
// To build gnorm, just mage build.

package main

import (
	"errors"
	"log"
	"os"
	"regexp"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Runs go install for gnorm.  This generates the embedded docs and the version
// info into the binary.
func Build() error {
	mg.Deps(Generate)
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

// Runs go generate.
func Generate() error {
	return sh.Run("go", "generate", "./...")
}

// Generates a new release.  Expects the TAG environment variable to be set,
// which will create a new tag with that name.
func Release() (err error) {
	releaseTag := regexp.MustCompile(`^v1\.[0-9]+\.[0-9]+$`)

	tag := os.Getenv("TAG")
	if !releaseTag.MatchString(tag) {
		return errors.New("TAG environment variable must be in semver v1.x.x format, but was " + tag)
	}

	if err := sh.RunV("git", "tag", "-a", tag, "-m", tag); err != nil {
		return err
	}
	if err := sh.RunV("git", "push", "origin", tag); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			sh.RunV("git", "tag", "--delete", "$TAG")
			sh.RunV("git", "push", "--delete", "origin", "$TAG")
		}
	}()
	return sh.RunV("goreleaser")
}

// Removes generated cruft.  This target shouldn't ever be necessary, since the
// cleanup should happen automatically, but it's here just in case.
func Clean() error {
	return cleanup()
}
