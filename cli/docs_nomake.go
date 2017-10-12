//+build !make

package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"gnorm.org/gnorm/environ"
)

func showDocs(env environ.Values, cmd *cobra.Command, args []string) error {
	fmt.Fprintln(env.Stdout, "docs not available, you need to build with mage build")

	return nil
}
