package cbrotli

/*
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <brotli/encode.h>
*/
import "C"
import (
	"bytes"
	"io"
	"sync/atomic"
	"unsafe"
)

var GetPutInfo atomic.Int64

const DefaultQuality = 3

// WWriter for support pool.
type WWriter struct {
	WriterV2
	Quality int
}
type WriterV2Options struct {
	Quality     int
	LGWin       int
	LGBlock     int
	LargeWindow bool
}

// NewWriterV2 only CGO Call one time.
func NewWriterV2(dst io.Writer, options WriterV2Options) *WriterV2 {
	largeWindow := 0
	if options.LargeWindow {
		largeWindow = 1
	}
	state := C.BrotliEncoderStateCreate(C.int(options.Quality), C.int(options.LGBlock), C.int(options.LGWin), C.int(largeWindow))
	return &WriterV2{
		dst:   dst,
		state: state,
	}
}

func NewWWriter(dst io.Writer, options WriterV2Options) *WWriter {
	return &WWriter{
		WriterV2: *NewWriterV2(dst, options),
		Quality:  options.Quality,
	}
}
func (w *WWriter) Write(p []byte) (n int, err error) {
	return w.WriterV2.Write(p)
}

func (w *WWriter) Close() (err error) {
	_, err = w.writeChunk(nil, C.BROTLI_OPERATION_FINISH)
	C.BrotliEncoderResetState(w.state)
	w.dst = nil
	return
}
func (w *WWriter) Reset(dst io.Writer) {
	w.dst = dst
}

func (w *WWriter) Flush() error {
	return w.WriterV2.Flush()
}
func (w *WWriter) Destroy() {
	C.BrotliEncoderDestroyInstance(w.state)
	w.state = nil
}

// WriterV2 original Writer every write chunk with another two c call.
// Now only one. has_more and bytes_consumed return with compressed result.
type WriterV2 struct {
	dst   io.Writer
	state *C.BrotliEncoderState
}

// after every write chunk, should consume state's internal out_buf, because we don't
// provide buf. if we don't consume it's internal buffer, and it is full, input consume nothing.
func (w *WriterV2) writeChunk(p []byte, op C.BrotliEncoderOperation) (n int, err error) {
	if w.state == nil {
		return 0, errWriterClosed
	}

	for {
		var data *C.uint8_t
		if len(p) != 0 {
			data = (*C.uint8_t)(&p[0])
		}
		result := C.CompressStreamV2(w.state, op, data, C.size_t(len(p)))

		if result.success == 0 {
			return n, errEncode
		}
		p = p[int(result.bytes_consumed):]
		n += int(result.bytes_consumed)

		length := int(result.output_data_size)
		if length != 0 {
			// It is a workaround for non-copying-wrapping of native memory.
			// C-encoder never pushes output block longer than ((2 << 25) + 502).
			// TODO(eustas): use natural wrapper, when it becomes available, see
			//               https://golang.org/issue/13656.
			output := (*[1 << 30]byte)(unsafe.Pointer(result.output_data))[:length:length]
			_, err = w.dst.Write(output)
			if err != nil {
				return n, err
			}
		}
		if len(p) == 0 && result.has_more == 0 {
			return n, nil
		}
	}
}
func (w *WriterV2) Write(p []byte) (n int, err error) {
	return w.writeChunk(p, C.BROTLI_OPERATION_PROCESS)
}

// Flush outputs encoded data for all input provided to Write. The resulting
// output can be decoded to match all input before Flush, but the stream is
// not yet complete until after Close.
// Flush has a negative impact on compression.
func (w *WriterV2) Flush() error {
	_, err := w.writeChunk(nil, C.BROTLI_OPERATION_FLUSH)
	return err
}

// Close flushes remaining data to the decorated writer and frees C resources.
func (w *WriterV2) Close() error {
	// If stream is already closed, it is reported by `writeChunk`.
	_, err := w.writeChunk(nil, C.BROTLI_OPERATION_FINISH)
	// C-Brotli tolerates `nil` pointer here.
	C.BrotliEncoderDestroyInstance(w.state)
	w.state = nil
	return err
}

// CloseWithContent for reduce CGO call counts.
func (w *WriterV2) CloseWithContent(content []byte) error {
	// If stream is already closed, it is reported by `writeChunk`.
	_, err := w.writeChunk(content, C.BROTLI_OPERATION_FINISH)
	// C-Brotli tolerates `nil` pointer here.
	C.BrotliEncoderDestroyInstance(w.state)
	w.state = nil
	return err
}

// EncodeV2 returns content encoded with Brotli.
func EncodeV2(content []byte, options WriterV2Options) ([]byte, error) {
	var buf bytes.Buffer
	writer := NewWriterV2(&buf, options)
	_, err := writer.Write(content)
	if closeErr := writer.Close(); err == nil {
		err = closeErr
	}
	return buf.Bytes(), err
}
