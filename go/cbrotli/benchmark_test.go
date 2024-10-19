package cbrotli

import (
	gbytes "bytes"
	"fmt"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/newacorn/goutils/bytes"
	"github.com/xyproto/randomstring"
	"io"
	"math/rand"
	"testing"
)

var bytesContainer [][]byte

func init() {
	for i := 10; i <= 20; i++ {
		dataStr := randomstring.HumanFriendlyString(2 << i)
		bytesContainer = append(bytesContainer, []byte(dataStr))
	}
}

func TestWWriterEqualWriter(t *testing.T) {
	for i := 0; i < 12; i++ {
		t.Run("level "+fmt.Sprint(i), func(t *testing.T) {
			w1 := NewWWriter(nil, WriterV2Options{Quality: 4})
			for i := 0; i < 10; i++ {
				srcBytes := bytesContainer[rand.Int31n(int32(len(bytesContainer)))]
				buf1 := gbytes.Buffer{}
				buf2 := gbytes.Buffer{}
				w1.Reset(&buf1)
				w2 := NewWriter(&buf2, WriterOptions{Quality: 4})
				//
				_, err := w1.Write(srcBytes)
				assert.NoErr(t, err)
				_, err = w2.Write(srcBytes)
				assert.NoErr(t, err)
				err = w1.Close()
				assert.NoErr(t, err)
				err = w2.Close()
				assert.NoErr(t, err)
				assert.Eq(t, buf1.Bytes(), buf2.Bytes())
				assert.NotEq(t, len(srcBytes), buf1.Len())
				assert.NotEq(t, 0, buf1.Len())
				buf1.Reset()
				buf2.Reset()
			}
			w1.Destroy()
		})
	}
}

func TestWWriterCompressCorrect(t *testing.T) {
	const loopCount = 20
	for i := 0; i < 12; i++ {
		level := i
		w := NewWWriter(nil, WriterV2Options{Quality: level})
		for i := 0; i < loopCount; i++ {
			index := rand.Int31n(int32(len(bytesContainer)))
			srcBytes := bytesContainer[index]
			buf := bytes.NewBufferSizeNoPtr(len(srcBytes))
			//
			w.Reset(&buf)
			_, err := w.Write(srcBytes)
			assert.NoErr(t, err)
			err = w.Close()
			//
			r := NewReader(&buf)
			rs, err := io.ReadAll(r)
			assert.NoErr(t, err)
			assert.Eq(t, srcBytes, rs)
			assert.NoErr(t, err)
			buf.Reset()
			buf.RecycleItems()
		}
		w.Destroy()
	}

}
