package cmd

import (
	"fmt"
	"log"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "alerts",
		Short: "alerts is a alerts aggregator",
		Long:  `alerts is a alerts aggregator that will collect alerts from different sources and aggregate them`,
	}
)

// Execute runs the rootCmd command.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("rootCmd.Execute, Err: %w", err)
	}

	return nil
}

// init is cmd's package init function.
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile,
		"config", "", "config file (default is $HOME/.alerts.yaml)")

	rootCmd.AddCommand(MigrateCommand())
}

// initConfig configures cobra with an already configured config file.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalf("homdir.Dir, Err: %v", err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".alerts")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
