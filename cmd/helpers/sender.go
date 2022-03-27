/*
Copyright © 2022 Hubert Pawlak <hubertpawlak.dev>
*/
package helpers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func SendReadings(readings []SensorReading) {
	verbose := viper.GetBool("verbose")
	// Encode readings as JSON
	if jsonBody, err := json.Marshal(readings); err == nil {
		// Don't send data - flag
		if viper.GetBool("readOnly") {
			log.Printf("%v", string(jsonBody))
			return
		}
		if verbose {
			log.Printf("Sending: %v", string(jsonBody))
		}
		// Prepare request
		endpoint := viper.GetString("endpoint")
		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
		if err != nil {
			log.Printf("%v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+viper.GetString("token"))
		res, err := http.DefaultClient.Do(req) // Send it
		// FIXME: close connection (or somehow cache/reuse)
		if err != nil {
			log.Printf("%v", err)
		}
		defer res.Body.Close()
		if verbose && res != nil {
			log.Printf("Endpoint returned %v", res.Status)
		}
	} else {
		log.Printf("%v", err)
	}
}
