package helper

import (
	"bytes"
	"encoding/json"
	"math"
	"strconv"
	"testing"
)

func Test_unstableInt(t *testing.T) {
	tests := []struct {
		name    string
		in      []byte
		out     int
		wantErr bool
	}{
		{"1", []byte(`"1"`), 1, false},
		{"2", []byte(`""`), 0, false},
		{"3", []byte(`"9223372036854775807"`), math.MaxInt, false},
		{"4", []byte(`"-9223372036854775808"`), math.MinInt, false},
		{"5", []byte(`"0"`), 0, false},
		{"6", []byte(`"9.9"`), 0, true},
		{"6", []byte(`"`), 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := new(UnstableInt)
			if err := temp.UnmarshalJSON(tt.in); (err != nil) != tt.wantErr {
				t.Errorf("unstableInt.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.out != int(*temp) {
				t.Errorf("unstableInt.UnmarshalJSON() = %v, want %v", *temp, tt.out)
			}
			out, err := temp.MarshalJSON()
			if err != nil {
				t.Errorf("unstableInt.MarshalJSON() = %v, want %v", out, tt.out)
			}
			if !bytes.Equal(out, []byte(strconv.Itoa(tt.out))) {
				t.Errorf("unstableInt.MarshalJSON() = %v, want %v", out, tt.in)
			}
		})
	}

	type f struct {
		in  []byte
		out []byte
	}
	tests2 := []f{
		{[]byte(`{"A":"1"}`), []byte(`{"A":1}`)},
		{[]byte(`{"A":""}`), []byte(`{"A":0}`)},
		{[]byte(`{}`), []byte(`{"A":null}`)},
		{[]byte(`{"A":"9223372036854775807"}`), []byte(`{"A":9223372036854775807}`)},
		{[]byte(`{"A":"-1"}`), []byte(`{"A":-1}`)},
	}
	for _, tt := range tests2 {
		t.Run(string(tt.in), func(t *testing.T) {
			plain, plains := tt.in, tt.out
			var s struct {
				A *UnstableInt
			}
			// fmt.Println(1, s.A)
			err := json.Unmarshal(plain, &s)
			if err != nil {
				t.Errorf("unstableInt.UnmarshalJSON() error = %v", err)
			}
			// fmt.Println(2, s.A)

			b, err := json.Marshal(s)
			if err != nil {
				t.Errorf("unstableInt.MarshalJSON() error = %v", err)
			}
			if !bytes.Equal(b, plains) {
				t.Errorf("unstableInt.MarshalJSON() = %v, want %v", b, plains)
			}
		})
	}
}

func Test_UnstableFloat(t *testing.T) {
	tests := []struct {
		name    string
		in      []byte
		out     float64
		wantErr bool
	}{
		{"1", []byte(`"1.0"`), 1, false},
		{"2", []byte(`""`), 0, false},
		{"3", []byte(`"1.79769e+308"`), 1.79769e+308, false},
		{"5", []byte(`"0.0"`), 0, false},
		{"6", []byte(`"9.99"`), 9.99, false},
		{"7", []byte(`"`), 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := new(UnstableFloat)
			if err := temp.UnmarshalJSON(tt.in); (err != nil) != tt.wantErr {
				t.Errorf("UnstableFloat.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.out != float64(*temp) {
				t.Errorf("UnstableFloat.UnmarshalJSON() = %v, want %v", *temp, tt.out)
			}
			out, err := temp.MarshalJSON()
			if err != nil {
				t.Errorf("UnstableFloat.MarshalJSON() = %v, want %v", out, tt.out)
			}
			if !bytes.Equal(out, []byte(strconv.FormatFloat(tt.out, 'f', -1, 64))) {
				t.Errorf("UnstableFloat.MarshalJSON() = %v, want %v", out, tt.in)
			}
		})
	}

	type f struct {
		in  []byte
		out []byte
	}
	tests2 := []f{
		{[]byte(`{"A":"1.1"}`), []byte(`{"A":1.1}`)},
		{[]byte(`{"A":""}`), []byte(`{"A":0}`)},
		{[]byte(`{}`), []byte(`{"A":null}`)},
		{[]byte(`{"A":"922.22"}`), []byte(`{"A":922.22}`)},
		{[]byte(`{"A":"-10.0"}`), []byte(`{"A":-10}`)},
	}
	for _, tt := range tests2 {
		t.Run(string(tt.in), func(t *testing.T) {
			plain, plains := tt.in, tt.out
			var s struct {
				A *UnstableFloat
			}
			// fmt.Println(1, s.A)
			err := json.Unmarshal(plain, &s)
			if err != nil {
				t.Errorf("UnstableFloat.UnmarshalJSON() error = %v", err)
			}
			// fmt.Println(2, s.A)

			b, err := json.Marshal(s)
			if err != nil {
				t.Errorf("UnstableFloat.MarshalJSON() error = %v", err)
			}
			if !bytes.Equal(b, plains) {
				t.Errorf("UnstableFloat.MarshalJSON() = %s, want %s", string(b), string(plains))
			}
		})
	}
}
