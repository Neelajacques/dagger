// nolint:dupl
package cmd

import (
	"dagger.io/go/cmd/dagger/logger"
	"dagger.io/go/dagger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Bring a route online with latest layout and inputs",
	Args:  cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Fix Viper bug for duplicate flags:
		// https://github.com/spf13/viper/issues/233
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			panic(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		lg := logger.New()
		ctx := lg.WithContext(cmd.Context())
		store, err := dagger.DefaultStore()
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to load store")
		}

		routeName := getRouteName(ctx)
		st, err := store.LookupRouteByName(ctx, routeName)
		if err != nil {
			lg.
				Fatal().
				Err(err).
				Str("routeName", routeName).
				Msg("failed to lookup route")
		}

		route, err := dagger.NewRoute(st)
		if err != nil {
			lg.
				Fatal().
				Err(err).
				Msg("failed to initialize route")
		}

		// TODO: Implement options: --no-cache
		routeUp(ctx, route)
	},
}

func init() {
	newCmd.Flags().Bool("--no-cache", false, "Disable all run cache")

	if err := viper.BindPFlags(upCmd.Flags()); err != nil {
		panic(err)
	}
}
