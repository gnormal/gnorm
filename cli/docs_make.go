//+build make

package cli

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"

	_ "gnorm.org/gnorm/cli/statik"
	"gnorm.org/gnorm/environ"
)

func showDocs(env environ.Values, cmd *cobra.Command, args []string) error {
	// this folder gets briefly copied here during go run make.go
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(statikFS))
	fmt.Fprintln(env.Stdout, "serving docs at http://localhost:8080")
	fmt.Fprintln(env.Stdout, "hit ctrl-C to cancel")
	go func() {
		if err := browser.OpenURL("http://localhost:8080"); err != nil {
			fmt.Println("failed to open browser")
		}
	}()
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return codeErr{errors.WithMessage(err, "can't serve docs"), 1}
	}

	return nil
}
