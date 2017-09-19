//+build ignore

// This is the "makefile" for gnorm.  To build norm, just run go run make.go.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const usage = `usage: go run make.go [command]

make.go is the build script for gnorm.

commands:
  install	compile with go install [default]
  all		build for all supported platforms
  help		display this help
`

func main() {
	pack()
	// commented for debugging purposes
	//defer mustRemove("./site/public")

	switch len(os.Args) {
	case 1:
		// (default)
		pack()
		fmt.Print(run("packr", "install", "--ldflags="+flags(), "gnorm.org/gnorm"))
	case 2:
		switch os.Args[1] {
		case "install":
			pack()
			fmt.Print(run("packr", "install", "--ldflags="+flags(), "gnorm.org/gnorm"))
		case "all":
			pack()
			ldf := flags()
			for _, OS := range []string{"windows", "darwin", "linux"} {
				if err := os.Setenv("GOOS", OS); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				for _, ARCH := range []string{"amd64", "386"} {
					if err := os.Setenv("GOOS", OS); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					fmt.Print(run("packr", "build", "-o", "gnorm_"+OS+"_"+ARCH, "--ldflags="+ldf, "gnorm.org/gnorm"))
				}
			}
		case "help":
			fmt.Println(usage)
		default:
			fmt.Println(usage)
			os.Exit(1)
		}
	}
}

func flags() string {
	timestamp := time.Now().Format(time.RFC3339)
	hash := run("git", "rev-parse", "HEAD")
	version := gitTag()
	return fmt.Sprintf(`-X "gnorm.org/gnorm/cli.timestamp=%s" -X "gnorm.org/gnorm/cli.commitHash=%s" -X "gnorm.org/gnorm/cli.version=%s"`, timestamp, hash, version)
}

func run(cmd string, args ...string) string {
	c := exec.Command(cmd, args...)
	c.Stderr = os.Stderr
	b, err := c.Output()
	if err != nil {
		fmt.Print(string(b))
		os.Exit(1)
	}
	return string(b)
}

func gitTag() string {
	c := exec.Command("git", "describe", "--tags")
	b, err := c.Output()
	if err != nil {
		exit, ok := err.(*exec.ExitError)
		if ok && exit.Exited() {
			// probably no git tag
			return "dev"
		}
		fmt.Print(string(b))
		os.Exit(1)
	}

	return strings.TrimSuffix(string(b), "\n")
}

func pack() {
	run("go", "get", "github.com/gobuffalo/packr/...")
	mustRemove("./site/public")
	run("hugo", "-s", "./site")
	// fonts are BIG
	mustRemove("./site/public/fonts")
	mustRemove("./site/public/revealjs/lib/font")
}

func mustRemove(s string) {
	err := os.RemoveAll(s)
	if !os.IsNotExist(err) && err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
