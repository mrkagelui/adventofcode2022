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

func Test_parse(t *testing.T) {
	tests := []struct {
		name   string
		str    string
		want   node
		assert check
	}{
		{
			name:   "malformed",
			str:    "[10,[]",
			want:   node{},
			assert: errorSays("malformed part: {[10,[]}"),
		},
		{
			name: "simple empty",
			str:  "[]",
			want: node{
				subs: []node{},
			},
			assert: ok,
		},
		{
			name: "simple",
			str:  "[10,11]",
			want: node{
				subs: []node{
					{value: 10},
					{value: 11},
				},
			},
			assert: ok,
		},
		{
			name: "deep",
			str:  "[[[21]]]",
			want: node{
				subs: []node{
					{
						subs: []node{
							{
								subs: []node{
									{
										value: 21,
									},
								},
							},
						},
					},
				},
			},
			assert: ok,
		},
		{
			name: "empty",
			str:  "[]",
			want: node{
				subs: []node{},
			},
			assert: ok,
		},
		{
			name: "multiple",
			str:  "[33,[1,2],[9,8],90]",
			want: node{
				subs: []node{
					{value: 33},
					{subs: []node{{value: 1}, {value: 2}}},
					{subs: []node{{value: 9}, {value: 8}}},
					{value: 90},
				},
			},
			assert: ok,
		},
		{
			name: "long",
			str:  "[[],[[[0],[8,10,2,8],[4]],[[7,7,2,2],10,1,2,[]]],[[[6,1,6,8,10],[8,6,4],[],[],2]]]",
			want: node{
				subs: []node{
					{
						subs: []node{},
					},
					{
						subs: []node{
							{
								subs: []node{
									{
										subs: []node{
											{value: 0},
										},
									},
									{
										subs: []node{
											{value: 8},
											{value: 10},
											{value: 2},
											{value: 8},
										},
									},
									{
										subs: []node{
											{value: 4},
										},
									},
								},
							},
							{
								subs: []node{
									{
										subs: []node{
											{value: 7},
											{value: 7},
											{value: 2},
											{value: 2},
										},
									},
									{value: 10},
									{value: 1},
									{value: 2},
									{subs: []node{}},
								},
							},
						},
					},
					{
						subs: []node{
							{
								subs: []node{
									{subs: []node{{value: 6}, {value: 1}, {value: 6}, {value: 8}, {value: 10}}},
									{subs: []node{{value: 8}, {value: 6}, {value: 4}}},
									{subs: []node{}},
									{subs: []node{}},
									{value: 2},
								},
							},
						},
					},
				},
			},
			assert: ok,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parse(tt.str)
			tt.assert(t, err)
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() got: \n%+v\n, want: \n%+v\n", got, tt.want)
			}
			if err == nil && got1 != len(tt.str) {
				t.Errorf("parse() got1 = %v, want %v", got1, len(tt.str))
			}
		})
	}
}

func Test_smallerThan(t *testing.T) {
	type args struct {
		a node
		b node
	}
	tests := []struct {
		name string
		args args
		want *bool
	}{
		{
			name: "both empty",
			args: args{
				a: node{},
				b: node{},
			},
			want: nil,
		},
		{
			name: "left empty",
			args: args{
				a: node{subs: []node{}},
				b: node{subs: []node{{value: 8}, {value: 7}, {value: 6}}},
			},
			want: ptrOf(true),
		},
		{
			name: "right empty",
			args: args{
				a: node{subs: []node{{value: 8}, {value: 7}, {value: 6}}},
				b: node{subs: []node{}},
			},
			want: ptrOf(false),
		},
		{
			name: "left deeper",
			args: args{
				a: node{subs: []node{{subs: []node{}}}},
				b: node{subs: []node{}},
			},
			want: ptrOf(false),
		},
		{
			name: "both values, left smaller",
			args: args{
				a: node{value: 100},
				b: node{value: 101},
			},
			want: ptrOf(true),
		},
		{
			name: "both values, left bigger",
			args: args{
				a: node{value: 100},
				b: node{value: 20},
			},
			want: ptrOf(false),
		},
		{
			name: "left value, right list",
			args: args{
				a: node{value: 9},
				b: node{subs: []node{{value: 8}, {value: 7}, {value: 6}}},
			},
			want: ptrOf(false),
		},
		{
			name: "left list, right value",
			args: args{
				a: node{subs: []node{{value: 8}, {value: 7}, {value: 6}}},
				b: node{value: 9},
			},
			want: ptrOf(true),
		},
		{
			name: "both list - same length and value",
			args: args{
				a: node{subs: []node{{value: 8}, {value: 7}, {value: 6}}},
				b: node{subs: []node{{value: 8}, {value: 7}, {value: 6}}},
			},
			want: nil,
		},
		{
			name: "both list - left smaller value",
			args: args{
				a: node{subs: []node{{value: 8}, {value: 7}, {value: 6}}},
				b: node{subs: []node{{value: 8}, {value: 9}, {value: 6}}},
			},
			want: ptrOf(true),
		},
		{
			name: "both list - right smaller value",
			args: args{
				a: node{subs: []node{{value: 8}, {value: 7}, {value: 6}}},
				b: node{subs: []node{{value: 8}, {value: 7}, {value: 2}}},
			},
			want: ptrOf(false),
		},
		{
			name: "both list - left shorter",
			args: args{
				a: node{subs: []node{{value: 8}, {value: 7}}},
				b: node{subs: []node{{value: 8}, {value: 7}, {value: 2}}},
			},
			want: ptrOf(true),
		},
		{
			name: "both list - right shorter",
			args: args{
				a: node{subs: []node{{value: 8}, {value: 7}, {value: 6}}},
				b: node{subs: []node{{value: 8}, {value: 7}}},
			},
			want: ptrOf(false),
		},
		{
			name: "both list - right shorter but value bigger",
			args: args{
				a: node{subs: []node{{value: 8}, {value: 7}, {value: 6}}},
				b: node{subs: []node{{value: 9}, {value: 7}}},
			},
			want: ptrOf(true),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := smallerThan(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("smallerThan() = %v, want %v", got, tt.want)
			}
		})
	}
}
