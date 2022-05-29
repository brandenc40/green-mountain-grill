package gmg

import (
	"reflect"
	"testing"
)

func TestGetGrillInfoResponseToGrillInfo(t *testing.T) {
	type args struct {
		response []byte
	}
	tests := []struct {
		name        string
		args        args
		expectError bool
		want        *State
	}{
		{
			name: "power off",
			args: args{[]byte{0x55, 0x52, 0x66, 0x0, 0x59, 0x2, 0x96, 0x0, 0x5, 0xb, 0x14, 0x32, 0x19, 0x19, 0x19, 0x19, 0x59, 0x2, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x3}},
			want: &State{
				CurrentTemperature:      102,
				TargetTemperature:       150,
				Probe1Temperature:       0,
				Probe1TargetTemperature: 0,
				Probe2Temperature:       0,
				Probe2TargetTemperature: 0,
				WarnCode:                WarnCodeNone,
				PowerState:              PowerStateOff,
				FireState:               FireStateOff,
			},
		},
		{
			name: "power on cold smoke",
			args: args{[]byte{0x55, 0x52, 0x66, 0x0, 0x59, 0x2, 0x1e, 0x0, 0x5, 0xb, 0x14, 0x32, 0x19, 0x19, 0x19, 0x19, 0x59, 0x2, 0xfa, 0x0, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0, 0x96, 0x0, 0x3, 0x0, 0xc6, 0x0, 0x0, 0x3}},
			want: &State{
				CurrentTemperature:      102,
				TargetTemperature:       30,
				Probe1Temperature:       0,
				Probe1TargetTemperature: 150,
				Probe2Temperature:       0,
				Probe2TargetTemperature: 250,
				WarnCode:                WarnCodeNone,
				PowerState:              PowerStateColdSmoke,
				FireState:               FireStateColdSmoke,
			},
		},
		{
			name:        "error: unexpected EOF",
			expectError: true,
			args:        args{[]byte{0x55, 0x52}},
			want:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BytesToState(tt.args.response)
			if err != nil && !tt.expectError {
				t.Errorf("BytesToState() returned an error: %v", err)
			}
			if !tt.expectError && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BytesToState() = %v, want %v", got, tt.want)
			}
		})
	}
}
