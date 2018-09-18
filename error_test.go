package main

import (
	"errors"
	"testing"
)

func Test_fatalError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{"err", args{errors.New("some fatalError")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fatalError(tt.args.err)
		})
	}
}

func Test_warningError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{"err", args{errors.New("some warningError")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warningError(tt.args.err)
		})
	}
}
