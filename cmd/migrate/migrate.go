package migrate

import (
	"github.com/spf13/cobra"
	"hub/cmd/migrate/down"
	"hub/cmd/migrate/drop"
	"hub/cmd/migrate/reset"
	"hub/cmd/migrate/up"
)

func NewMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migration commands",
		Long:  `Manage database migrations - up, down, reset, drop`,
	}

	// Add subcommands
	cmd.AddCommand(up.NewCommand())
	cmd.AddCommand(down.NewCommand())
	cmd.AddCommand(reset.NewCommand())
	cmd.AddCommand(drop.NewCommand())

	return cmd
}
