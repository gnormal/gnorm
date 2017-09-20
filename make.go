//+build ignore

// This is the "makefile" for gnorm.  To build norm, just run go run make.go.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const usage = `
usage: go run make.go [flags] [command]

make.go is the build script for gnorm.

flags:

  -v		show verbose output

commands:
  install	compile with go install [default]
  all		build for all supported platforms
  help		display this help
`

var verbose = false

func main() {
	log.SetFlags(0)
	flag.BoolVar(&verbose, "v", false, "")
	flag.Parse()
	if !verbose {
		log.SetOutput(ioutil.Discard)
	}
	if len(flag.Args()) > 1 {
		Fatal("invalid args" + usage)
	}

	genSite()
	log.Print("downloading statik")
	run("go", "get", "github.com/rakyll/statik")
	log.Print("generating statik embedded files")
	run("statik", "-f", "-src", "./cli/public", "-dest", "./cli")
	defer func() {
		log.Print("removing generated hugo site")
		mustRemove("./cli/public")
		log.Print("removing generated statik package")
		mustRemove("./cli/statik")
	}()
	switch len(flag.Args()) {
	case 0:
		// (default)
		log.Print("running go install")
		run("go", "install", "-tags", "make", "--ldflags="+flags(), "gnorm.org/gnorm")
	case 1:
		switch flag.Args()[1] {
		case "install":
			log.Print("running go install")
			run("go", "install", "-tags", "make", "--ldflags="+flags(), "gnorm.org/gnorm")
		case "all":
			makeAll()
		case "help":
			log.Print(usage)
		default:
			Fatal(usage)
		}
	default:
		// we already checked for this, but just in case
		Fatal("invalid args" + usage)
	}
}

func makeAll() {
	ldf := flags()
	for _, OS := range []string{"windows", "darwin", "linux"} {
		if err := os.Setenv("GOOS", OS); err != nil {
			Fatal(err)
		}
		for _, ARCH := range []string{"amd64", "386"} {
			if err := os.Setenv("GOOS", OS); err != nil {
				Fatal(err)
			}
			log.Printf("running go build for GOOS=%s GOARCH=%s", OS, ARCH)
			run("go", "build", "-tags", "make", "-o", "gnorm_"+OS+"_"+ARCH, "--ldflags="+ldf)
		}
	}
}

func flags() string {
	timestamp := time.Now().Format(time.RFC3339)
	hash := output("git", "rev-parse", "HEAD")
	version := gitTag()
	return fmt.Sprintf(`-X "gnorm.org/gnorm/cli.timestamp=%s" -X "gnorm.org/gnorm/cli.commitHash=%s" -X "gnorm.org/gnorm/cli.version=%s"`, timestamp, hash, version)
}

func run(cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.Stderr = os.Stderr
	if verbose {
		c.Stdout = os.Stdout
	}
	err := c.Run()
	if err != nil {
		Fatal(err)
	}
}

func output(cmd string, args ...string) string {
	c := exec.Command(cmd, args...)
	c.Stderr = os.Stderr
	b, err := c.Output()
	if err != nil {
		log.Print(string(b))
		Fatal(err)
	}
	return string(b)
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
		Fatal(string(b))
	}

	return strings.TrimSuffix(string(b), "\n")
}

func genSite() {
	log.Print("cleaning up any existing hugo generated files")
	mustRemove("./cli/public")
	log.Print("downloading hugo")
	run("go", "get", "github.com/gohugoio/hugo")
	log.Print("generating docs site")
	run("hugo", "-s", "./site", "-d", "../cli/public")
	log.Print("removing fonts from generated site")
	// fonts are BIG
	mustRemove("./cli/public/fonts")
	mustRemove("./cli/public/revealjs/lib/font")
}

func mustRemove(s string) {
	err := os.RemoveAll(s)
	if !os.IsNotExist(err) && err != nil {
		log.Fatal(err)
	}
}

func Fatal(args ...interface{}) {
	log.SetOutput(os.Stdout)
	log.Fatal(args...)
}
