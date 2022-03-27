/*
Copyright © 2022 Hubert Pawlak <hubertpawlak.dev>
*/
package helpers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func SendReadings(readings []SensorReading, endpoint string, token string, verbose bool) error {
	// Encode readings as JSON
	jsonBody, err := json.Marshal(readings)
	if err != nil {
		return err
	}
	// Print body
	if verbose {
		log.Printf("readingsJson: %v", string(jsonBody))
	}
	// Prepare POST request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if len(token) > 0 {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := http.DefaultClient.Do(req) // Send it
	if err != nil {
		return err
	}
	res.Body.Close() // Prevent resource leak
	if verbose && res != nil {
		log.Printf("Endpoint returned %v", res.Status)
	}
	return nil
}
