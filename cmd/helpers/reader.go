package helpers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/spf13/afero"
)

var SensorsPath string = "/sys/bus/w1/devices"

func GetAllSensors(fs afero.Afero) ([]string, error) {
	sensors := make([]string, 0)
	// Check if dir exists
	exists, err := fs.DirExists(SensorsPath)
	if !exists {
		err := errors.New(SensorsPath + " does not exist or is not a directory")
		return sensors, err
	} else if err != nil {
		return sensors, err
	}
	// Get the directory content
	sensorsDir, _ := fs.ReadDir(SensorsPath)
	// Turn it into a list of sensors
	for _, sensor := range sensorsDir {
		if name := sensor.Name(); !strings.HasPrefix(name, "w1_bus_master") {
			sensors = append(sensors, name)
		}
	}
	return sensors, nil
}

type SensorReading struct {
	SensorHwId  string  `json:"hw_id"`
	Temperature float32 `json:"temperature"`
	Resolution  int     `json:"resolution"`
}

func GetReadings(fs afero.Afero, sensors []string) ([]SensorReading, error) {
	// Read values for each sensor
	readings := make([]SensorReading, 0)
	for _, sensorId := range sensors {
		// Paths
		sensorBasePath := SensorsPath + "/" + sensorId
		temperaturePath := sensorBasePath + "/temperature"
		resolutionPath := sensorBasePath + "/resolution"
		// File access
		temperatureFile, readTempErr := fs.ReadFile(temperaturePath)
		if readTempErr != nil {
			return nil, readTempErr
		}
		resolutionFile, readResErr := fs.ReadFile(resolutionPath)
		if readResErr != nil {
			return nil, readResErr
		}
		// Temperature, remove trailing \n and convert
		rawTemperature, convTempErr := strconv.Atoi(strings.TrimSpace(string(temperatureFile)))
		if convTempErr != nil {
			return nil, convTempErr
		}
		temperature := float32(rawTemperature) / 1000 // Temperature has to be divided by 1000
		// Resolution, must be positive
		resolution, convResErr := strconv.Atoi(strings.TrimSpace(string(resolutionFile)))
		if convResErr != nil {
			return nil, convResErr
		}
		if resolution <= 0 {
			return nil, errors.New("resolution must be positive")
		}
		// Add results to array
		readings = append(readings, SensorReading{SensorHwId: sensorId, Temperature: temperature, Resolution: resolution})
	}
	// Everything went fine, return readings
	return readings, nil
}
