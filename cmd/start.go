/*
Copyright © 2022 Hubert Pawlak <hubertpawlak.dev>
*/
package cmd

import (
	"errors"
	"time"

	"github.com/hubertpawlak/ds18b20-remote-sensor/cmd/helpers"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start sending temperature readings",
	Long: `Start sending temperature readings from sensors
connected using 1-Wire to your specified endpoint.`,
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if all flags required for sending data are provided
		if !viper.GetBool("readOnly") {
			switch {
			case viper.GetString("endpoint") == "":
				return errors.New("endpoint is required")
			case viper.GetString("token") == "":
				return errors.New("token is required")
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Setup ticker
		ticker := time.NewTicker(time.Duration(viper.GetUint("interval")) * time.Millisecond)
		defer ticker.Stop()
		tick := ticker.C
		// Setup FS
		fs := afero.Afero{Fs: afero.NewOsFs()}
		// Run the first tick instantly
		for ; true; <-tick {
			// Sensors may be connected/disconnected at any time
			if sensors, err := helpers.GetAllSensors(fs); err == nil {
				if readings, err := helpers.GetReadings(fs, sensors); err == nil {
					helpers.SendReadings(readings)
				} else {
					panic(err)
				}
			} else {
				panic(err)
			}
			// TODO: measure measurement time and increase interval if greater than current interval
			// ticker.Reset() + viper.set
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// TODO: multiple endpoints
	// startCmd.Flags().StringArrayP("endpoint", "e", "destination(s) of your readings")
	startCmd.Flags().StringP("endpoint", "e", "", "destination of your readings (API)")
	startCmd.Flags().UintP("interval", "i", 5000, "time between readings (in ms)")
	startCmd.Flags().StringP("token", "t", "", "secret key used to authorize requests")
	startCmd.Flags().Int("led", 0, "output GPIO PIN number")
	startCmd.Flags().Bool("readOnly", false, "print readings instead of sending (useful for testing)")
	startCmd.Flags().MarkDeprecated("readOnly", "it will be removed in the next version")

	viper.BindPFlag("endpoint", startCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("interval", startCmd.Flags().Lookup("interval"))
	viper.BindPFlag("token", startCmd.Flags().Lookup("token"))
	viper.BindPFlag("led", startCmd.Flags().Lookup("led"))
	viper.BindPFlag("readOnly", startCmd.Flags().Lookup("readOnly"))
}
