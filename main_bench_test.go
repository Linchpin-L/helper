package helper

import (
	"encoding/json"
	"strconv"
	"testing"
)

func BenchmarkNumToAlphabet(t *testing.B) {
	for i := 0; i < t.N; i++ {
		_ = NumToAlphabet(int64(i))
	}
}

func BenchmarkJsonImplement(t *testing.B) {
	for i := 0; i < t.N; i++ {
		tt := &struct {
			A *UnstableInt
		}{}
		err := json.Unmarshal([]byte(`{"A":"`+strconv.Itoa(i)+`"}`), tt)
		if err != nil {
			t.Error(err)
		}
		_, err = json.Marshal(tt)
		if err != nil {
			t.Error(err)
		}
		// t.Log(string(rsp))
	}
}
