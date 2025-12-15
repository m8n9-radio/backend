package drop

import (
	"hub/internal/database"
	"hub/internal/wire"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "drop",
		Short: "Drop all database tables",
		Long:  `Drop all tables in the database schema. WARNING: This is destructive!`,
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := wire.InitializeMigrateApp()
			if err != nil {
				return err
			}

			mg := database.NewMigrate(app.DSN)
			return mg.Drop()
		},
	}
}
