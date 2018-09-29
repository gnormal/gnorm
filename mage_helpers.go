// +build mage

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func flags() (string, error) {
	timestamp := time.Now().Format(time.RFC3339)
	hash, err := output("git", "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	version, err := gitTag()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`-X "gnorm.org/gnorm/cli.timestamp=%s" -X "gnorm.org/gnorm/cli.commitHash=%s" -X "gnorm.org/gnorm/cli.version=%s"`, timestamp, hash, version), nil
}

func gitTag() (string, error) {
	s, err := output("git", "describe", "--tags")
	if err != nil {
		ee, ok := errors.Cause(err).(*exec.ExitError)
		if ok && ee.Exited() {
			// probably no git tag
			return "dev", nil
		}
		return "", err
	}

	return strings.TrimSuffix(s, "\n"), nil
}

func cleanup() error {
	log.Print("removing generated hugo site")
	err := rm("./cli/public")
	if err != nil {
		fmt.Println("error removing generated hugo folder:", err)
	}
	log.Print("removing generated statik package")
	if err := rm("./cli/statik"); err != nil {
		fmt.Println("error removing statik folder:", err)
	}
	if err := rm("dist"); err != nil {
		fmt.Println("error removing release folder:", err)
	}
	return nil
}

func genSite() error {
	log.Print("cleaning up any existing hugo generated files")

	if err := rm("./cli/public"); err != nil {
		return err
	}
	log.Print("generating docs site")
	if err := run("hugo", "-s", "./site", "-d", "../cli/public"); err != nil {
		return err
	}
	log.Print("removing fonts from generated site")
	// fonts are BIG
	// if err := rm("./cli/public/fonts"); err != nil {
	// 	return err
	// }
	if err := rm("./cli/public/revealjs/lib/font"); err != nil {
		return err
	}

	log.Print("generating statik embedded files")
	return run("statik", "-f", "-src", "./cli/public", "-dest", "./cli")
}

func rm(s string) error {
	err := os.RemoveAll(s)
	if os.IsNotExist(err) {
		return nil
	}
	return errors.Wrapf(err, `failed to remove %s`, s)
}

func run(cmd string, args ...string) error {
	return runWith(nil, cmd, args...)
}

func runWith(env []string, cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Env = os.Environ()
	for _, v := range env {
		c.Env = append(c.Env, v)
	}
	c.Stderr = os.Stderr
	if os.Getenv("MAGEFILE_VERBOSE") != "" {
		c.Stdout = os.Stdout
	}
	return errors.Wrapf(c.Run(), `failed to run %v %q`, cmd, args)
}

func output(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	c.Env = os.Environ()
	c.Stderr = os.Stderr
	b, err := c.Output()
	if err != nil {
		return "", errors.Wrapf(err, `failed to run %v %q`, cmd, args)
	}
	return string(b), nil
}
