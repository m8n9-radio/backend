package down

import (
	"hub/internal/database"
	"hub/internal/wire"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var steps int

	cmd := &cobra.Command{
		Use:   "down",
		Short: "Rollback migrations",
		Long:  `Rollback all migrations or rollback specific number of steps.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := wire.InitializeMigrateApp()
			if err != nil {
				return err
			}

			mg := database.NewMigrate(app.DSN)
			return mg.Down(steps)
		},
	}

	cmd.Flags().IntVarP(&steps, "steps", "s", 0, "Number of migrations to rollback (0 = all)")

	return cmd
}
