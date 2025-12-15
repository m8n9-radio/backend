package reset

import (
	"hub/internal/database"
	"hub/internal/wire"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset all migrations (down + up)",
		Long:  `Rollback all migrations and then run them again.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := wire.InitializeMigrateApp()
			if err != nil {
				return err
			}

			mg := database.NewMigrate(app.DSN)

			// Rollback all migrations
			if err := mg.Down(0); err != nil {
				return err
			}

			// Run all migrations
			return mg.Up(0)
		},
	}
}
