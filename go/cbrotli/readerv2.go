package cbrotli

/*
#include <stddef.h>
#include <stdint.h>

#include <brotli/decode.h>
*/
import "C"
import (
	pool "github.com/newacorn/simple-bytes-pool"
	"io"
)

// ReaderV2 implements io.ReadCloser by reading Brotli-encoded data from an
// underlying Reader.
type ReaderV2 struct {
	Reader
	recycleItems []*pool.Bytes
	hasMore      bool
}

func (r *ReaderV2) RecycleItems() {
	for _, item := range r.recycleItems {
		item.RecycleToPool00()
	}
	r.recycleItems = r.recycleItems[:0]
}

func (r *ReaderV2) Close() error {
	err := r.Reader.Close()
	r.RecycleItems()
	return err
}

func NewReader2WithContent(src []byte) *ReaderV2 {
	r := Reader{
		state: C.BrotliDecoderCreateInstance(nil, nil, nil),
		buf:   src,
	}
	return &ReaderV2{
		Reader: r,
	}
}
func NewReaderV2(src io.Reader) *ReaderV2 {
	py := pool.Get(readBufSize)
	py.B = py.B[:readBufSize]
	r := Reader{
		src:   src,
		state: C.BrotliDecoderCreateInstance(nil, nil, nil),
		buf:   py.B,
	}
	return &ReaderV2{
		Reader:       r,
		recycleItems: []*pool.Bytes{py},
	}
}

func (r *ReaderV2) Reset(src io.Reader) error {
	r.Reader.src = src
	return nil
}

func (r *ReaderV2) Read(p []byte) (n int, err error) {
	if r.state == nil {
		return 0, errReaderClosed
	}
	if !r.hasMore && len(r.in) == 0 {
		m, readErr := r.src.Read(r.buf)
		if m == 0 {
			// If readErr is `nil`, we just proxy underlying stream behavior.
			return 0, readErr
		}
		r.in = r.buf[:m]
	}

	if len(p) == 0 {
		return 0, nil
	}

	for {
		var written, consumed, hasMore C.size_t
		var data *C.uint8_t
		if len(r.in) != 0 {
			data = (*C.uint8_t)(&r.in[0])
		}
		result := C.DecompressStreamV2(r.state,
			(*C.uint8_t)(&p[0]), C.size_t(len(p)),
			data, C.size_t(len(r.in)),
			&written, &consumed, &hasMore)
		r.in = r.in[int(consumed):]
		n = int(written)
		r.hasMore = false
		if hasMore != 0 {
			r.hasMore = true
		}
		switch result {
		case C.BROTLI_DECODER_RESULT_SUCCESS:
			if len(r.in) > 0 {
				return n, errExcessiveInput
			}
			return n, nil
		case C.BROTLI_DECODER_RESULT_ERROR:
			return n, decodeError(C.BrotliDecoderGetErrorCode(r.state))
		case C.BROTLI_DECODER_RESULT_NEEDS_MORE_OUTPUT:
			if n == 0 {
				return 0, io.ErrShortBuffer
			}
			return n, nil
		case C.BROTLI_DECODER_NEEDS_MORE_INPUT:
		}

		if len(r.in) != 0 {
			return 0, errInvalidState
		}

		// Calling r.src.Read may block. Don't block if we have data to return.
		if n > 0 {
			return n, nil
		}

		// Top off the buffer.
		encN, err := r.src.Read(r.buf)
		if encN == 0 {
			// Not enough data to complete decoding.
			if err == io.EOF {
				return 0, io.ErrUnexpectedEOF
			}
			return 0, err
		}
		r.in = r.buf[:encN]
	}

}
