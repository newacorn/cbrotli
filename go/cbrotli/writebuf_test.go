package cbrotli

import (
	"github.com/andybalholm/brotli"
	pbytes "github.com/newacorn/bytes-pool/bytes"
	"github.com/xyproto/randomstring"
	"sync"
	"testing"
)

func BenchmarkWrite(b *testing.B) {
	const dataLen = 1200
	dataStr := randomstring.HumanFriendlyString(dataLen)
	dataBytes := []byte(dataStr)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := pbytes.NewBufferWithSize(dataLen)
		w := NewWriter(buf, WriterOptions{Quality: 4})
		_, _ = w.Write(dataBytes)
		_ = w.Close()
		buf.RecycleItems()
	}
}

var p = sync.Pool{
	New: func() interface{} {
		return brotli.NewWriterLevel(nil, 4)
	},
}
