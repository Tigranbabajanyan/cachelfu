package cachelfu

import (
	"strconv"
	"testing"
)

func BenchmarkLfuAdd(b *testing.B) {

	b.ResetTimer()

	lfu := New(10000)

	for i := 0; i < b.N; i++ {
		err := lfu.Add(strconv.FormatInt(int64(i), 10), i)
		if err != nil {
			b.Fatal(err)
		}
	}
}
