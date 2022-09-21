package helpers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendReadings(t *testing.T) {
	type args struct {
		readings []SensorReading
		endpoint string // used only if mockResponseCode == 0
		token    string
		verbose  bool
	}
	tests := []struct {
		name             string
		args             args
		wantErr          bool
		mockResponseCode int // 0 by default, use real endpoint if not 0
	}{
		{
			name: "invalid endpoint",
			args: args{
				readings: nil,
				endpoint: "1234567890endpoint",
			},
			wantErr: true,
		},
		{
			name: "unreachable endpoint",
			args: args{
				readings: nil,
				endpoint: "http://127.0.0.1:1",
			},
			wantErr: true,
		},
		{
			name: "valid negative reading",
			args: args{
				readings: []SensorReading{
					{SensorHwId: "s1", Temperature: -12.345, Resolution: 12},
				},
			},
			mockResponseCode: http.StatusOK,
			wantErr:          false,
		},
		{
			name: "valid negative reading with token",
			args: args{
				readings: []SensorReading{
					{SensorHwId: "s1", Temperature: -12.345, Resolution: 12},
				},
				token: "1234567890token",
			},
			mockResponseCode: http.StatusOK,
			wantErr:          false,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var server *httptest.Server
			if testCase.mockResponseCode != 0 {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(testCase.mockResponseCode)
				}))
				testCase.args.endpoint = server.URL // Replace endpoint with mock server
				defer server.Close()
			}
			if err := SendReadings(testCase.args.readings, testCase.args.endpoint, testCase.args.token, testCase.args.verbose); (err != nil) != testCase.wantErr {
				t.Errorf("SendReadings() error = %v, wantErr %v", err, testCase.wantErr)
			}
		})
	}
}
