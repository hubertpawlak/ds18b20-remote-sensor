/*
Copyright © 2022 Hubert Pawlak <hubertpawlak.dev>
*/
package helpers

import (
	"log"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var sensorsPath string

// Make it easy to swap for a fake FS
var fs = afero.NewOsFs()
var afs = afero.Afero{Fs: fs}

func switchToFakeFs() {
	fs = afero.NewMemMapFs()
	afs.Fs = fs
	afs.MkdirAll(sensorsPath, 0750)
	afs.Mkdir(sensorsPath+"/w1_bus_master1", 0750)
	afs.WriteFile(sensorsPath+"/28-000000000001/temperature", []byte("12345\n"), 0740)
	afs.WriteFile(sensorsPath+"/28-000000000001/resolution", []byte("8\n"), 0740)
	afs.WriteFile(sensorsPath+"/28-000000000002/temperature", []byte("102000\n"), 0740)
	afs.WriteFile(sensorsPath+"/28-000000000002/resolution", []byte("12\n"), 0740)
}

func AutoFakeSensors() {
	// Check if faking MUST be enabled
	if fakingEnabled, os := viper.GetBool("fakeData"), runtime.GOOS; !fakingEnabled && os != "linux" {
		viper.Set("fakeData", true)
		log.Printf("Forced fake readings due to OS incompatibility")
	}
	// Setup fake sensors
	if fakingEnabled := viper.GetBool("fakeData"); fakingEnabled {
		log.Printf("Fake sensor readings are ENABLED")
		sensorsPath = "/fake/devices"
		switchToFakeFs()
		return
	}
	// Use real sensors
	sensorsPath = "/sys/bus/w1/devices"
}

func GetAllSensors() []string {
	// Check if dir exists
	if exists, _ := afs.DirExists(sensorsPath); !exists {
		log.Fatalf("%v does not exist or is not a directory", sensorsPath)
	}
	// Get the directory content
	sensorsDir, _ := afs.ReadDir(sensorsPath)
	// Turn it into a list of sensors
	sensors := make([]string, 0)
	for _, sensor := range sensorsDir {
		if name := sensor.Name(); !strings.HasPrefix(name, "w1_bus_master") {
			sensors = append(sensors, name)
		}
	}
	return sensors
}

type SensorReading struct {
	SensorHwId  string  `json:"sensorHwId"`
	Temperature float32 `json:"temperature"`
	Resolution  int     `json:"resolution"`
}

func GetReadings(sensors []string) []SensorReading {
	// verbose := viper.GetBool("verbose")
	// Read values for each sensor
	readings := make([]SensorReading, 0)
	for _, sensorId := range sensors {
		// Paths
		sensorBasePath := sensorsPath + "/" + sensorId
		temperaturePath := sensorBasePath + "/temperature"
		resolutionPath := sensorBasePath + "/resolution"
		// File access
		temperatureFile, errT := afs.ReadFile(temperaturePath)
		resolutionFile, errR := afs.ReadFile(resolutionPath)
		if errT == nil && errR == nil {
			// Remove trailing \n and convert
			rawTemperature, errT := strconv.Atoi(strings.TrimSpace(string(temperatureFile)))
			temperature := float32(rawTemperature) / 1000
			resolution, errR := strconv.Atoi(strings.TrimSpace(string(resolutionFile)))
			if errT == nil && errR == nil {
				// Temperature has to be divided 1000
				readings = append(readings, SensorReading{SensorHwId: sensorId, Temperature: temperature, Resolution: resolution})
			} else {
				log.Printf("R1: %v: %v", sensorBasePath, errT)
				log.Printf("R2: %v: %v", sensorBasePath, errR)
			}
		} else {
			log.Printf("R3: %v: %v", sensorBasePath, errT)
			log.Printf("R4: %v: %v", sensorBasePath, errR)
		}
	}
	return readings
}
