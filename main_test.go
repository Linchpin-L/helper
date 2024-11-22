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
