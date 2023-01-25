package main

import (
	"reflect"
	"testing"
)

func Test_parseValve(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		wantName string
		want     valve
	}{
		{
			name:     "1",
			s:        "Valve IO has flow rate=20; tunnels lead to valves YT, TX",
			wantName: "IO",
			want: valve{
				flow: 20,
				next: []string{"YT", "TX"},
			},
		},
		{
			name:     "2",
			s:        "Valve AA has flow rate=0; tunnels lead to valves NY, IA, WK, FU, NU",
			wantName: "AA",
			want: valve{
				flow: 0,
				next: []string{"NY", "IA", "WK", "FU", "NU"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, got, err := parseValve(tt.s)
			if err != nil {
				t.Errorf("parseValve() error = %v", err)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("parseValve() gotName = %v, wantName %v", gotName, tt.wantName)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseValve() got = %v, want %v", got, tt.want)
			}
		})
	}
}
