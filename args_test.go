package main

import "testing"

func TestArgs_setDefault(t *testing.T) {
	type fields struct {
		FileName   string
		Sep        string
		SkipSymbol string
		SkipNum    int
		Header     int
		Transpose  bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"set default argument", fields{"", "", "", 0, 0, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := &Args{
				FileName:   tt.fields.FileName,
				Sep:        tt.fields.Sep,
				SkipSymbol: tt.fields.SkipSymbol,
				SkipNum:    tt.fields.SkipNum,
				Header:     tt.fields.Header,
				Transpose:  tt.fields.Transpose,
			}
			args.setDefault()
		})
	}
}
