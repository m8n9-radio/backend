package serve

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"hub/internal/wire"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewServeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP server",
		Long:  `Start the HTTP server and listen for incoming requests.`,
		Run: func(cmd *cobra.Command, args []string) {
			app, cleanup, err := wire.InitializeApp()
			if err != nil {
				panic(err)
			}
			defer cleanup()
			defer app.Database.Pool().Close()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

			g, gCtx := errgroup.WithContext(ctx)

			g.Go(func() error {
				app.Logger.Infof("starting HTTP server on port %d", app.Config.Port())
				return app.Server.Listen(app.Config.Port())
			})

			g.Go(func() error {
				return app.Scheduler.Start(gCtx)
			})

			g.Go(func() error {
				select {
				case sig := <-sigChan:
					app.Logger.Infof("received signal: %v, initiating graceful shutdown", sig)
				case <-gCtx.Done():
					app.Logger.Info("context cancelled")
				}
				cancel()
				return app.Server.Shutdown(ctx)
			})

			if err := g.Wait(); err != nil {
				if err != context.Canceled {
					app.Logger.Errorf("application error: %v", err)
				}
			}

			app.Logger.Info("application shutdown complete")
		},
	}
}
