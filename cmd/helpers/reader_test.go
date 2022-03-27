/*
Copyright © 2022 Hubert Pawlak <hubertpawlak.dev>
*/
package helpers

import (
	"testing"

	"github.com/spf13/afero"
)

func TestGetAllSensors(t *testing.T) {
	t.Parallel()
	t.Run("no dir", func(t *testing.T) {
		t.Parallel()
		fs := afero.Afero{Fs: afero.NewMemMapFs()} // Fresh empty fake FS before running test
		if _, err := GetAllSensors(fs); err == nil {
			t.Errorf("GetAllSensors should fail if %v does not exist", SensorsPath)
		}
	})
	t.Run("empty dir", func(t *testing.T) {
		t.Parallel()
		fs := afero.Afero{Fs: afero.NewMemMapFs()} // Fresh empty fake FS before running test
		fs.MkdirAll(SensorsPath, 0750)             // Create empty sensors directory
		sensors, err := GetAllSensors(fs)
		if err != nil {
			t.Error(err)
		}
		if len(sensors) != 0 {
			t.Error("GetAllSensors should return an empty array")
		}
	})
	t.Run("master without sensors", func(t *testing.T) {
		t.Parallel()
		fs := afero.Afero{Fs: afero.NewMemMapFs()} // Fresh empty fake FS before running test
		fs.MkdirAll(SensorsPath, 0750)
		fs.Mkdir(SensorsPath+"/w1_bus_master1", 0750) // This should be filtered out by GetAllSensors
		sensors, err := GetAllSensors(fs)
		if err != nil {
			t.Error(err)
		}
		if len(sensors) != 0 {
			t.Error("GetAllSensors should return an empty array")
		}
	})
	t.Run("2 sensors", func(t *testing.T) {
		t.Parallel()
		fs := afero.Afero{Fs: afero.NewMemMapFs()} // Fresh empty fake FS before running test
		// Fake dir with sensors
		fs.MkdirAll(SensorsPath, 0750)
		fs.Mkdir(SensorsPath+"/w1_bus_master1", 0750)
		fs.WriteFile(SensorsPath+"/id1/temperature", []byte("12345\n"), 0740)
		fs.WriteFile(SensorsPath+"/id1/resolution", []byte("8\n"), 0740)
		fs.WriteFile(SensorsPath+"/id2/temperature", []byte("102000\n"), 0740)
		fs.WriteFile(SensorsPath+"/id2/resolution", []byte("12\n"), 0740)
		// Test
		sensors, _ := GetAllSensors(fs)
		if sensorsLen := len(sensors); sensorsLen != 2 {
			t.Errorf("sensors length should be 2, was %v: %v", sensorsLen, sensors)
		}
	})
	t.Run("unusual sensor id", func(t *testing.T) {
		t.Parallel()
		fs := afero.Afero{Fs: afero.NewMemMapFs()} // Fresh empty fake FS before running test
		fs.MkdirAll(SensorsPath, 0750)
		fs.Mkdir(SensorsPath+"/w1_bus_master1", 0750)
		fs.WriteFile(SensorsPath+"/unusualId/temperature", []byte("12345\n"), 0740)
		fs.WriteFile(SensorsPath+"/unusualId/resolution", []byte("12\n"), 0740)
		// Test
		sensors, _ := GetAllSensors(fs)
		if sensorsLen := len(sensors); sensorsLen != 1 {
			t.Errorf("sensors length should be 1, was %v: %v", sensorsLen, sensors)
		}
	})
}

func TestGetReadings(t *testing.T) {
	t.Parallel()
	t.Run("no sensor dir", func(t *testing.T) {
		t.Parallel()
		fs := afero.Afero{Fs: afero.NewMemMapFs()} // Fresh empty fake FS before running test
		sensors := []string{"id1", "id2"}          // Simulate GetAllSensors
		readings, err := GetReadings(fs, sensors)
		if readingsLen := len(readings); readingsLen != 0 {
			t.Errorf("GetReadings should not return any readings, got %v", readingsLen)
		}
		if err == nil {
			t.Errorf("GetReadings should fail if there is no %v", SensorsPath)
		}
	})
	t.Run("no temperature file", func(t *testing.T) {
		t.Parallel()
		fs := afero.Afero{Fs: afero.NewMemMapFs()} // Fresh empty fake FS before running test
		fs.MkdirAll(SensorsPath+"/w1_bus_master1", 0750)
		fs.WriteFile(SensorsPath+"/id1/resolution", []byte("8\n"), 0740)
		sensors := []string{"id1"}
		readings, err := GetReadings(fs, sensors)
		if readingsLen := len(readings); readingsLen != 0 {
			t.Errorf("GetReadings should not return any readings, got %v", readingsLen)
		}
		if err == nil {
			t.Error("GetReadings should fail if there is no temperature file")
		}
	})
	t.Run("no resolution file", func(t *testing.T) {
		t.Parallel()
		fs := afero.Afero{Fs: afero.NewMemMapFs()} // Fresh empty fake FS before running test
		fs.MkdirAll(SensorsPath+"/w1_bus_master1", 0750)
		fs.WriteFile(SensorsPath+"/id1/temperature", []byte("12345\n"), 0740)
		sensors := []string{"id1"}
		readings, err := GetReadings(fs, sensors)
		if readingsLen := len(readings); readingsLen != 0 {
			t.Errorf("GetReadings should not return any readings, got %v", readingsLen)
		}
		if err == nil {
			t.Error("GetReadings should fail if there is no resolution file")
		}
	})
	t.Run("invalid temperature", func(t *testing.T) {
		t.Parallel()
		fs := afero.Afero{Fs: afero.NewMemMapFs()} // Fresh empty fake FS before running test
		fs.MkdirAll(SensorsPath+"/w1_bus_master1", 0750)
		fs.WriteFile(SensorsPath+"/id1/temperature", []byte("12345\n"), 0740) // This one is fine
		fs.WriteFile(SensorsPath+"/id1/resolution", []byte("12\n"), 0740)
		fs.WriteFile(SensorsPath+"/id2/temperature", []byte("102.000\n"), 0740) // This one is invalid
		fs.WriteFile(SensorsPath+"/id2/resolution", []byte("12\n"), 0740)
		sensors := []string{"id1", "id2"}
		readings, err := GetReadings(fs, sensors)
		if readingsLen := len(readings); readingsLen != 0 {
			t.Errorf("GetReadings should not return any readings, got %v", readingsLen)
		}
		if err == nil {
			t.Error("GetReadings should fail if temperature is invalid")
		}
	})
	t.Run("invalid resolution", func(t *testing.T) {
		t.Parallel()
		fs := afero.Afero{Fs: afero.NewMemMapFs()} // Fresh empty fake FS before running test
		fs.MkdirAll(SensorsPath+"/w1_bus_master1", 0750)
		fs.WriteFile(SensorsPath+"/id1/temperature", []byte("12345\n"), 0740)
		fs.WriteFile(SensorsPath+"/id1/resolution", []byte("1\n"), 0740) // This one is fine
		fs.WriteFile(SensorsPath+"/id2/temperature", []byte("102000\n"), 0740)
		fs.WriteFile(SensorsPath+"/id2/resolution", []byte("0\n"), 0740) // This one is invalid
		sensors := []string{"id1", "id2"}
		readings, err := GetReadings(fs, sensors)
		if readingsLen := len(readings); readingsLen != 0 {
			t.Errorf("GetReadings should not return any readings, got %v", readingsLen)
		}
		if err == nil {
			t.Error("GetReadings should fail if temperature is invalid")
		}
	})
	t.Run("valid positive reading", func(t *testing.T) {
		t.Parallel()
		fs := afero.Afero{Fs: afero.NewMemMapFs()} // Fresh empty fake FS before running test
		fs.MkdirAll(SensorsPath+"/w1_bus_master1", 0750)
		fs.WriteFile(SensorsPath+"/id1/temperature", []byte("12345\n"), 0740)
		fs.WriteFile(SensorsPath+"/id1/resolution", []byte("12\n"), 0740)
		sensors := []string{"id1"}
		readings, err := GetReadings(fs, sensors)
		if readingsLen := len(readings); readingsLen != 1 {
			t.Errorf("GetReadings should return 1 reading, got %v", readingsLen)
		}
		if err != nil {
			t.Error("GetReadings should not fail")
		}
		if got := readings[0].SensorHwId; got != "id1" {
			t.Errorf("SensorHwId changed, got %v", got)
		}
		if got := readings[0].Temperature; got != 12.345 {
			t.Errorf("Returned temperature does not match, got %v", got)
		}
		if got := readings[0].Resolution; got != 12 {
			t.Errorf("Returned resolution does not match, got %v", got)
		}
	})
	t.Run("valid negative reading", func(t *testing.T) {
		t.Parallel()
		fs := afero.Afero{Fs: afero.NewMemMapFs()} // Fresh empty fake FS before running test
		fs.MkdirAll(SensorsPath+"/w1_bus_master1", 0750)
		fs.WriteFile(SensorsPath+"/id1/temperature", []byte("-12345\n"), 0740)
		fs.WriteFile(SensorsPath+"/id1/resolution", []byte("12\n"), 0740)
		sensors := []string{"id1"}
		readings, err := GetReadings(fs, sensors)
		if readingsLen := len(readings); readingsLen != 1 {
			t.Errorf("GetReadings should return 1 reading, got %v", readingsLen)
		}
		if err != nil {
			t.Error("GetReadings should not fail")
		}
		if got := readings[0].SensorHwId; got != "id1" {
			t.Errorf("SensorHwId changed, got %v", got)
		}
		if got := readings[0].Temperature; got != -12.345 {
			t.Errorf("Returned temperature does not match, got %v", got)
		}
		if got := readings[0].Resolution; got != 12 {
			t.Errorf("Returned resolution does not match, got %v", got)
		}
	})
}
