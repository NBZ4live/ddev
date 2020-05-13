package cmd

import (
	"fmt"
	"github.com/drud/ddev/pkg/ddevapp"
	"github.com/drud/ddev/pkg/util"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var ArtisanCmd = &cobra.Command{
	Use:   "artisan [command]",
	Short: "Executes a artisan command within the web container",
	Long: `Executes a artisan command at the project root in the web container. Generally,
any artisan command can be forwarded to the container context by prepending
the command with 'ddev'. For example:

ddev artisan key:generate
ddev artisan make:migration [options] [--] <name>`,
	Run: func(cmd *cobra.Command, args []string) {
		app, err := ddevapp.GetActiveApp("")
		if err != nil {
			util.Failed(err.Error())
		}

		if app.SiteStatus() != ddevapp.SiteRunning {
			if err = app.Start(); err != nil {
				util.Failed("Failed to start %s: %v", app.Name, err)
			}
		}

		stdout, stderr, err := app.Exec(&ddevapp.ExecOpts{
			Service: "web",
			Dir:     "/var/www/html",
			Cmd:     fmt.Sprintf("php artisan %s", strings.Join(args, " ")),
			Tty:     isatty.IsTerminal(os.Stdin.Fd()),
		})
		if err != nil {
			util.Failed("php artisan %v failed, %v. stderr=%v", args, err, stderr)
		}
		_, _ = fmt.Fprint(os.Stderr, stderr)
		_, _ = fmt.Fprint(os.Stdout, stdout)
	},
}

func init() {
	RootCmd.AddCommand(ArtisanCmd)
	ArtisanCmd.Flags().SetInterspersed(false)
}
