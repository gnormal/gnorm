// +build mage

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func flags() string {
	timestamp := time.Now().Format(time.RFC3339)
	hash := output("git", "rev-parse", "HEAD")
	version := gitTag()
	return fmt.Sprintf(`-X "gnorm.org/gnorm/cli.timestamp=%s" -X "gnorm.org/gnorm/cli.commitHash=%s" -X "gnorm.org/gnorm/cli.version=%s"`, timestamp, hash, version)
}

func gitTag() string {
	c := exec.Command("git", "describe", "--tags")
	c.Stderr = os.Stderr
	b, err := c.Output()
	if err != nil {
		exit, ok := err.(*exec.ExitError)
		if ok && exit.Exited() {
			// probably no git tag
			return "dev"
		}
		must(err)
	}

	return strings.TrimSuffix(string(b), "\n")
}

func genSite() (cleanup func()) {
	log.Print("cleaning up any existing hugo generated files")
	mustRemove("./cli/public")
	log.Print("generating docs site")
	must(run("hugo", "-s", "./site", "-d", "../cli/public"))
	log.Print("removing fonts from generated site")
	// fonts are BIG
	mustRemove("./cli/public/fonts")
	mustRemove("./cli/public/revealjs/lib/font")

	log.Print("generating statik embedded files")
	must(run("statik", "-f", "-src", "./cli/public", "-dest", "./cli"))
	return func() {
		log.Print("removing generated hugo site")
		mustRemove("./cli/public")
		log.Print("removing generated statik package")
		mustRemove("./cli/statik")
	}
}

func mustRemove(s string) {
	err := os.RemoveAll(s)
	if !os.IsNotExist(err) && err != nil {
		log.Fatal(err)
	}
}

func run(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stderr = os.Stderr
	if os.Getenv("MAGEFILE_VERBOSE") != "" {
		c.Stdout = os.Stdout
	}
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}

func must(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func output(cmd string, args ...string) string {
	c := exec.Command(cmd, args...)
	c.Stderr = os.Stderr
	b, err := c.Output()
	must(err)
	return string(b)
}
