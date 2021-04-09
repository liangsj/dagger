package cmd

import (
	"dagger.io/go/cmd/dagger/cmd/common"
	"dagger.io/go/cmd/dagger/logger"
	"dagger.io/go/dagger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Bring a deployment online with latest plan and inputs",
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

		state := common.GetCurrentDeploymentState(ctx, store)

		// TODO: Implement options: --no-cache
		result := common.DeploymentUp(ctx, state)
		state.Computed = result.Computed().JSON().String()
		if err := store.UpdateDeployment(ctx, state, nil); err != nil {
			lg.Fatal().Err(err).Msg("failed to update deployment")
		}
	},
}

func init() {
	newCmd.Flags().Bool("--no-cache", false, "Disable all run cache")

	if err := viper.BindPFlags(upCmd.Flags()); err != nil {
		panic(err)
	}
}