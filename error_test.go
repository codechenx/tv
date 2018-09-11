package main

import (
	"errors"
	"testing"
)

func Test_warningError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{"warningError", args{errors.New("warning error")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warningError(tt.args.err)
		})
	}
}
