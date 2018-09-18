package main

import "testing"

func Test_initView(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"run"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initView()
		})
	}
}
