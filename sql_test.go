package helper

import (
	"database/sql/driver"
	"reflect"
	"testing"
)

func TestMutipleUint64_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		j       *MutipleUint64
		args    args
		result  []uint64
		wantErr bool
	}{
		{"1", new(MutipleUint64), args{[]byte("1,2,3")}, []uint64{1, 2, 3}, false},
		{"empty bytes", new(MutipleUint64), args{[]byte("")}, []uint64{}, false},
		{"single value", new(MutipleUint64), args{[]byte("42")}, []uint64{42}, false},
		{"multiple values with spaces", new(MutipleUint64), args{[]byte("10, 20,30")}, []uint64{10, 0, 30}, false},
		{"invalid type", new(MutipleUint64), args{123}, nil, true},
		{"contains invalid uint", new(MutipleUint64), args{[]byte("1,abc,3")}, []uint64{1, 0, 3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.j.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("MutipleUint64.Scan() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if !reflect.DeepEqual([]uint64(*tt.j), tt.result) {
					t.Errorf("MutipleUint64.Scan() = %v, want %v", *tt.j, tt.result)
				}
			}
		})
	}
}

func TestMutipleUint64_Value(t *testing.T) {
	tests := []struct {
		name    string
		j       MutipleUint64
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "empty slice",
			j:       MutipleUint64{},
			want:    "",
			wantErr: false,
		},
		{
			name:    "single value",
			j:       MutipleUint64{42},
			want:    "42",
			wantErr: false,
		},
		{
			name:    "multiple values",
			j:       MutipleUint64{1, 2, 3},
			want:    "1,2,3",
			wantErr: false,
		},
		{
			name:    "zero values",
			j:       MutipleUint64{0, 0, 0},
			want:    "0,0,0",
			wantErr: false,
		},
		{
			name:    "large numbers",
			j:       MutipleUint64{1234567890123456789, 9876543210987654321},
			want:    "1234567890123456789,9876543210987654321",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("MutipleUint64.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MutipleUint64.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
