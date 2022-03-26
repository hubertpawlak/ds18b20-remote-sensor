/*
Copyright © 2022 Hubert Pawlak <hubertpawlak.dev>
*/
package helpers

import (
	"log"

	"github.com/spf13/viper"
	"github.com/stianeikeland/go-rpio/v4"
)

// TODO: GPIO LED output support
func initGpio() {
	verbose := viper.GetBool("verbose")
	if err := rpio.Open(); err != nil {

	} else if verbose {
		log.Printf("%v", err)
	}
}

func closeGpio() {
	// TODO
}

func turnOnLed() {
	// TODO
}

func turnOffLed() {
	// TODO
}
