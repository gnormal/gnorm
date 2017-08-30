//+build ignore

// This is the "makefile" for gnorm.  To build norm, just run go run make.go.

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

const usage = `usage: go run make.go [command]

make.go is the build script for gnorm.

commands:
  install	compile with go install [default]
  build		compile with go build
  help		display this help
  
`

func main() {
	log.SetFlags(0)
	cmd := "install"
	if len(os.Args) > 2 {
		log.Printf("too many arguments: %q\n\n", os.Args[1:])
		log.Fatal(usage)
	}
	if len(os.Args) == 2 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "help":
		log.Print(usage)
		os.Exit(0)
	case "install", "build":
	// ok
	default:
		log.Printf("unknown command %q\n\n", cmd)
		log.Fatal(usage)
	}

	timestamp := time.Now().Format(time.RFC3339)
	hash := run("git", "rev-parse", "HEAD")
	flags := fmt.Sprintf(`-X "gnorm.org/gnorm/cli.timestamp=%s" -X "gnorm.org/gnorm/cli.commitHash=%s"`, timestamp, hash)
	fmt.Print(run("go", cmd, "--ldflags="+flags, "gnorm.org/gnorm"))
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
