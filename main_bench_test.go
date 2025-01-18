package helper

import "testing"

func BenchmarkNumToAlphabet(t *testing.B) {
	for i:=0; i<t.N; i++ {
		_ = NumToAlphabet(int64(i))
	}
}