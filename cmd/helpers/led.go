/*
Copyright © 2022 Hubert Pawlak <hubertpawlak.dev>
*/
package helpers

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// TODO: GPIO LED output support
// TODO: use another method
func initGpio(pin int) error {
	err := rpio.Open()
	ledPin := rpio.Pin(pin)
	ledPin.Output()
	return err
}

func closeGpio() error {
	err := rpio.Close()
	return err
}

func turnOnLed(pinNumber int) {
	pin := rpio.Pin(pinNumber)
	rpio.WritePin(pin, rpio.High)
}

func turnOffLed(pinNumber int) {
	pin := rpio.Pin(pinNumber)
	rpio.WritePin(pin, rpio.Low)
}
