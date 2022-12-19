package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to config file")

	rootCmd.AddCommand(serveCmd)
}

func initConfig() {
	v := viper.New()
	if cfgFile != "" {
		v.SetConfigFile(cfgFile)

		if err := v.ReadInConfig(); err != nil {
			panic(fmt.Errorf("reading config %s error: %w", v.ConfigFileUsed(), err))
		}
	}

	v.AutomaticEnv()

	envKeysMap := map[string]interface{}{}
	if err := mapstructure.Decode(*configs.Configuration, &envKeysMap); err != nil {
		panic(fmt.Errorf("decoding config to mapstructure error: %w", err))
	}

	for k := range envKeysMap {
		if bindErr := v.BindEnv(k); bindErr != nil {
			panic(fmt.Errorf("binding viper env variable '%s' error: %w", k, bindErr))
		}
	}

	if err := v.Unmarshal(configs.Configuration); err != nil {
		panic(fmt.Errorf("decoding env configuration error: %w", err))
	}
}
