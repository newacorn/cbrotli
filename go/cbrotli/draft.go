//go:build draft
// +build draft

package cbrotli

/*
// BrotliEncoderCompressStream() function internal dereferences available_out parameter,so must handle a valid value.
// Because we dont want this function copy data to next_out, so we set available_out = 0.
static BrotliEncodeResult CompressStreamWithResult(
    BrotliEncoderState* s, BrotliEncoderOperation op,
    const uint8_t* data, size_t data_size) {

  size_t available_in = data_size;
  const uint8_t* next_in = data;
  size_t available_out = 0;
  return BrotliEncoderCompressStreamWithResult(s, op,
      &available_in, &next_in, &available_out, 0, 0);
}
*/
import "C"
