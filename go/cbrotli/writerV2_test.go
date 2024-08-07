package cbrotli

import (
	"github.com/newacorn/bytes-pool/bytes"
	"io"
	"log"
	"os"
	"testing"
)

func BenchmarkCBrotli(b *testing.B) {
	curDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var filepath = curDir + "/testdata/jquery-3.7.1.js"
	f, err := os.Open(filepath)
	defer func() { _ = f.Close() }()
	if err != nil {
		b.Fatal(err)
	}
	srcBytes, err := io.ReadAll(f)
	if err != nil {
		b.Fatal(err)
	}
	buf := bytes.NewBufferSizeNoPtr(len(srcBytes))
	buf.Reset()
	var compressLen int
	b.Logf("compress data len:%d", len(srcBytes))
	b.Run("CBrotli", func(b *testing.B) {
		buf.Reset()
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			w := NewWriter(&buf, WriterOptions{Quality: 3})
			_, err = io.Copy(w, bytes.NewBuffer(srcBytes))
			if err != nil {
				b.Fatal(err)
			}
			err = w.Close()
			compressLen = buf.Len()
			buf.Reset()
			if err != nil {
				b.Fatal(err)
			}
		}
		b.Logf("CBrotli level-3 compressed ratio:%.3f; %d", float64(compressLen)/float64(len(srcBytes)), compressLen)
	})
	b.Run("CBrotli V2", func(b *testing.B) {
		buf.Reset()
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			w := NewWriterV2(&buf, WriterV2Options{Quality: 3})
			_, err = io.Copy(w, bytes.NewBuffer(srcBytes))
			if err != nil {
				b.Fatal(err)
			}
			err = w.Close()
			compressLen = buf.Len()
			buf.Reset()
			if err != nil {
				b.Fatal(err)
			}
		}
		b.Logf("CBrotli level-3 compressed ratio:%.3f; %d", float64(compressLen)/float64(len(srcBytes)), compressLen)
	})
	b.Run("CBrotli V2 Reuse", func(b *testing.B) {
		buf.Reset()
		b.ResetTimer()
		b.ReportAllocs()
		w := NewWWriter(nil, WriterV2Options{Quality: 3})
		for i := 0; i < b.N; i++ {
			w.Reset(&buf)
			_, err = io.Copy(w, bytes.NewBuffer(srcBytes))
			if err != nil {
				b.Fatal(err)
			}
			err = w.Close()
			compressLen = buf.Len()
			buf.Reset()
			if err != nil {
				b.Fatal(err)
			}
		}
		w.Destroy()
		b.Logf("CBrotli reuse level-3 compressed ratio:%.3f; %d", float64(compressLen)/float64(len(srcBytes)), compressLen)
	})

}
