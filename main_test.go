package helper

import (
	"reflect"
	"testing"
	"time"
)

func TestIsIDCard(t *testing.T) {
	type args struct {
		idcard string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "1", args: args{""}, want: false},
		{name: "2", args: args{"1"}, want: false},
		{name: "3", args: args{"440102198001021231"}, want: false},
		{name: "4", args: args{"440102198001021230"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIDCard(tt.args.idcard); got != tt.want {
				t.Errorf("IsIDCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeroClock(t *testing.T) {
	a := time.Date(2023, 1, 11, 12, 12, 12, 12, time.Local)
	b := time.Date(2024, 2, 29, 23, 59, 59, 59, time.Local)
	c := time.Date(2024, 2, 29, 0, 0, 0, 12, time.Local)
	e := time.Now()

	type args struct {
		t *time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{name: "1", args: args{&a}, want: time.Date(2023, 1, 11, 0, 0, 0, 0, time.Local)},
		{name: "2", args: args{&b}, want: time.Date(2024, 2, 29, 0, 0, 0, 0, time.Local)},
		{name: "3", args: args{&c}, want: time.Date(2024, 2, 29, 0, 0, 0, 0, time.Local)},
		{name: "4", args: args{nil}, want: time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, time.Local)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ZeroClock(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZeroClock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumToAlphabet(t *testing.T) {
	type args struct {
		num int64
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"1", args{0}, []byte("0")},
		{"2", args{-0}, []byte("0")},
		{"3", args{1}, []byte("1")},
		{"4", args{-1}, []byte("-1")},
		{"6", args{10}, []byte("a")},
		{"7", args{11}, []byte("b")},
		{"8", args{61}, []byte("1p")},
		{"9", args{62}, []byte("1q")},
		{"10", args{-653543}, []byte("-e09z")},
		{"5", args{9223372036854775807}, []byte("1y2p0ij32e8e7")},
		{"11", args{-9223372036854775808}, []byte("-1y2p0ij32e8e8")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumToAlphabet(tt.args.num); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NumToAlphabet3() = %v, want %v", got, tt.want)
			}
		})
	}
}
