/*
Copyright © 2022 Hubert Pawlak <hubertpawlak.dev>
*/
package helpers

import "testing"

func Test_initGpio(t *testing.T) {
	t.Skip("TODO")
	type args struct {
		pin int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "simple",
			args:    args{pin: 11},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initGpio(tt.args.pin); (err != nil) != tt.wantErr {
				t.Errorf("initGpio() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
