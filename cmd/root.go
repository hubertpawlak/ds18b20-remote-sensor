/*
Copyright © 2022 Hubert Pawlak <hubertpawlak.dev>
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var customConfigFile string

var rootCmd = &cobra.Command{
	Use:   "ds18b20-remote-sensor",
	Short: "This program sends temperature readings",
	Long: `This program periodically sends temperature readings (from 1-Wire sensors)
to a remote server using HTTP POST requests.
It supports Bearer token authentication, self-signed SSL certificates and LED output.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&customConfigFile, "config", "c", "", "path to config file (default \"./rsConfig.yaml\")")

	// Global flags stored in a config file
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "show more messages")

	// Remember to "sync" flags to later Get
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if customConfigFile != "" {
		viper.SetConfigFile(customConfigFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("rsConfig.yaml")
	}

	viper.SetEnvPrefix("RS")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		if viper.GetBool("verbose") {
			log.Printf("Using config file: %v", viper.ConfigFileUsed())
		}
	}
}
