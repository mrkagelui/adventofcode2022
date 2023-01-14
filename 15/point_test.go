package main

import (
	"reflect"
	"testing"
)

type check func(*testing.T, error)

func ok(t *testing.T, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func errorSays(s string) check {
	return func(t *testing.T, err error) {
		if err == nil {
			t.Error("expecting error but got nil")
			return
		}
		if eStr := err.Error(); eStr != s {
			t.Errorf("expect error [%v] but got [%v]", s, eStr)
		}
	}
}

func Test_parsePoint(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		want   point
		assert check
	}{
		{
			name:   "sanity",
			s:      "x=213, y=432",
			want:   point{432, 213},
			assert: ok,
		},
		{
			name:   "malformed",
			s:      "x=213,y=432",
			want:   point{},
			assert: errorSays("invalid str [x=213,y=432]"),
		},
		{
			name:   "nan",
			s:      "x=death, y=432",
			want:   point{},
			assert: errorSays(`parsing col in [x=death, y=432]: strconv.Atoi: parsing "death": invalid syntax`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePoint(tt.s)
			tt.assert(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePoint() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parse(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		want   pair
		assert check
	}{
		{
			name: "sanity",
			s:    "Sensor at x=3556832, y=3209801: closest beacon is at x=3520475, y=3164417",
			want: pair{
				sensor: point{3209801, 3556832},
				beacon: point{3164417, 3520475},
			},
			assert: ok,
		},
		{
			name:   "malformed",
			s:      "Sensor at x=3556832, y=3209801: farthest beacon is at x=3520475, y=3164417",
			want:   pair{},
			assert: errorSays("invalid str [Sensor at x=3556832, y=3209801: farthest beacon is at x=3520475, y=3164417]"),
		},
		{
			name:   "nan",
			s:      "Sensor at x=death, y=3209801: closest beacon is at x=3520475, y=3164417",
			want:   pair{},
			assert: errorSays(`parsing sensor in [Sensor at x=death, y=3209801: closest beacon is at x=3520475, y=3164417]: parsing col in [x=death, y=3209801]: strconv.Atoi: parsing "death": invalid syntax`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.s)
			tt.assert(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
