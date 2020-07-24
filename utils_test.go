package main

import (
	"testing"
)

func TestI2B(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"TestI2B 1", args{i: 1}, true},
		{"TestI2B 2", args{i: 2}, true},
		{"TestI2B 3", args{i: 0}, false},
		{"TestI2B 3", args{i: 0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := I2B(tt.args.i); got != tt.want {
				t.Errorf("I2B() = %v, want %v", got, tt.want)
			}
		})
	}
}
