package main

import (
	"testing"
)

func Test_formatNumber(t *testing.T) {
	type args struct {
		unformattedNum string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"1234567890", args{unformattedNum: "1234567890"}, "1234567890",
		},
		{
			"123 456 7891", args{unformattedNum: "123 456 7891"}, "1234567891",
		},
		{
			"(123) 456 7892", args{unformattedNum: "(123) 456 7892"}, "1234567892",
		},
		{
			"(123) 456-7893", args{unformattedNum: "(123) 456-7893"}, "1234567893",
		},
		{
			"123-456-7894", args{unformattedNum: "123-456-7894"}, "1234567894",
		},
		{
			"123-456-7890", args{unformattedNum: "123-456-7890"}, "1234567890",
		},
		{
			"1234567892", args{unformattedNum: "1234567892"}, "1234567892",
		},
		{
			"(123)456-7892", args{unformattedNum: "(123)456-7892"}, "1234567892",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatNumber(tt.args.unformattedNum); got != tt.want {
				t.Errorf("formatNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
