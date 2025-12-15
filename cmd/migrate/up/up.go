package up

import (
	"hub/internal/database"
	"hub/internal/wire"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var steps int

	cmd := &cobra.Command{
		Use:   "up",
		Short: "Run pending migrations",
		Long:  `Apply all pending database migrations or run specific number of steps.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := wire.InitializeMigrateApp()
			if err != nil {
				return err
			}

			mg := database.NewMigrate(app.DSN)
			return mg.Up(steps)
		},
	}

	cmd.Flags().IntVarP(&steps, "steps", "s", 0, "Number of migrations to run (0 = all)")

	return cmd
}
