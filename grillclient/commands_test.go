package grillclient

import (
	"reflect"
	"testing"
)

func TestCommand_Build(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		c    Command
		args args
		want []byte
	}{
		{
			name: "Set temp formats properly",
			c:    CommandSetGrillTemp,
			args: args{[]interface{}{150}},
			want: []byte("UT150!"),
		},
		{
			name: "Int pads with 0",
			c:    CommandSetGrillTemp,
			args: args{[]interface{}{20}},
			want: []byte("UT020!"),
		},
		{
			name: "Command without args",
			c:    CommandPowerOn,
			args: args{[]interface{}{}},
			want: []byte("UK001!"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Build(tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
