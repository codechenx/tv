package main

import "testing"

func Test_render(t *testing.T) {
	type args struct {
		b   *Buffer
		trs bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty buffer", args{&Buffer{name: "", sep: "", header: 0, vHRN: 0, vHCN: 0, rowNum: 0, colNum: 0, cont: [][]string{}}, false}, true},
		{"render", args{&Buffer{name: "", sep: "", header: 2, vHCN: 0, vHRN: 0, rowNum: 2, colNum: 4, cont: [][]string{{"some", "thing", "other", "thing"}, {"some", "thing", "other", "thing"}}}, false}, false},
		{"transpose render", args{&Buffer{name: "", sep: "", header: 2, vHCN: 0, vHRN: 0, rowNum: 2, colNum: 4, cont: [][]string{{"some", "thing", "other", "thing"}, {"some", "thing", "other", "thing"}}}, true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := render(tt.args.b, tt.args.trs); (err != nil) != tt.wantErr {
				t.Errorf("render() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
