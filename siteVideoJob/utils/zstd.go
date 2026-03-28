package utils

import "github.com/klauspost/compress/zstd"

var (
	enc, _ = zstd.NewWriter(nil, zstd.WithEncoderConcurrency(3), zstd.WithEncoderLevel(zstd.SpeedDefault))
	dec, _ = zstd.NewReader(nil, zstd.WithDecoderConcurrency(3))
)

func Compress(src []byte) []byte {
	return enc.EncodeAll(src, make([]byte, 0, len(src)))
}

func Decompress(dist []byte) ([]byte, error) {
	return dec.DecodeAll(dist, nil)
}
