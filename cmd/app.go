package cmd

import (
	"context"
	"fmt"

	"Edot/config"
	"Edot/modules"
	"Edot/packages/logger"
	"Edot/packages/postgres"
	"Edot/routers"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var app = &cobra.Command{
	Use:   "start",
	Short: "Running service",
	Run: func(cmd *cobra.Command, args []string) {
		fx.New(
			fx.Provide(routers.NewRouter),
			fx.Provide(postgres.NewPostgres),
			fx.Provide(logger.NewLogger),
			modules.AppRepository,
			modules.AppController,
			modules.AppRoute,
			fx.Invoke(registerHooks),
		).Run()
	},
}

func init() {
	rootCmd.AddCommand(app)
}

func registerHooks(lifecycle fx.Lifecycle, echoRoute *routers.Router, psql *postgres.DB, logger *logger.Logger) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go echoRoute.Start(fmt.Sprintf(":%d", config.Get().Port))
				return nil
			},
			OnStop: func(ctx context.Context) error {
				if err := echoRoute.Shutdown(ctx); err != nil {
					logger.Fatal(err.Error())
					return err
				}
				defer func() {
					if err := psql.Sql.Close(); err != nil {
						logger.Fatal(err.Error())
					}
				}()
				return nil
			},
		},
	)
}
