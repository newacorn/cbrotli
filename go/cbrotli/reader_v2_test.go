package cbrotli

import (
	"bytes"
	"github.com/andybalholm/brotli"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/xyproto/randomstring"
	"io"
	"os"
	"testing"
)

func TestReaderV2Create(t *testing.T) {
	srcStr := randomstring.HumanFriendlyString(1301)
	srcBytes := []byte(srcStr)
	buf := bytes.Buffer{}
	w := NewWriterV2(&buf, WriterV2Options{Quality: 3})
	_, err := w.Write(srcBytes)
	assert.NoErr(t, err)
	err = w.Close()
	assert.NoErr(t, err)
	//
	r := NewReaderV2(&buf)
	rs, err := io.ReadAll(r)
	assert.NoErr(t, err)
	assert.Eq(t, srcBytes, rs)
	err = r.Close()
	assert.NoErr(t, err)
}

func BenchmarkUnCompress(b *testing.B) {
	curDir, err := os.Getwd()
	if err != nil {
		b.Fatal(err)
	}
	r, err := os.Open(curDir + "/testdata/jquery-3.7.1.js")
	if err != nil {
		b.Fatal(err)
	}
	srcBytesUn, err := io.ReadAll(r)
	buf := bytes.Buffer{}
	w := NewWriterV2(&buf, WriterV2Options{Quality: 3})
	_, err = w.Write(srcBytesUn)
	if err != nil {
		b.Fatal()
	}
	err = w.Close()
	if err != nil {
		b.Fatal(err)
	}
	srcBytes := buf.Bytes()
	//
	brotli.NewReader(bytes.NewReader(srcBytes))
	b.Run("CBrotli Reader", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		var un []byte
		for i := 0; i < b.N; i++ {
			r := NewReaderV2(bytes.NewReader(srcBytes))
			un, err = io.ReadAll(r)
			if err != nil {
				b.Fatal(err)
			}
			err = r.Close()
			if err != nil {
				b.Fatal(err)
			}
		}
		if !bytes.Equal(un, srcBytesUn) {
			b.Fatal("not equal")
		}
	})
	b.ResetTimer()
	b.ReportAllocs()
	b.Run("brotli.NewReader", func(b *testing.B) {
		var un []byte
		b.ResetTimer()
		b.ReportAllocs()
		r := brotli.NewReader(bytes.NewReader(srcBytes))
		for i := 0; i < b.N; i++ {
			un, err = io.ReadAll(r)
			if err != nil {
				b.Fatal(err)
			}
			err = r.Reset(bytes.NewReader(srcBytes))
			if err != nil {
				b.Fatal(err)
			}
		}
		if !bytes.Equal(un, srcBytesUn) {
			b.Fatal("not equal")
		}
		_ = r.Close()
	})
}
