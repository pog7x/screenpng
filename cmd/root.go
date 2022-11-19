package cmd

import (
	"fmt"
	"os"

	"github.com/pog7x/screenpng/configs"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "screenpng",
	Short: "Screenshot Factory",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)

	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"./configs/.screenpng-config.dev.yml",
		"path to config file",
	)

	initConfig()
}

func initConfig() {
	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("reading config %s error: %w", viper.ConfigFileUsed(), err))
	}

	if err := viper.Unmarshal(configs.Configuration); err != nil {
		panic(fmt.Errorf("decoding config %s error: %w", viper.ConfigFileUsed(), err))
	}
}
