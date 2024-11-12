package service

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
)

// RandomColor generates a random color.
func RandomColor() string {
	return fmt.Sprintf("#%06X", rand.Intn(0xFFFFFF)) //nolint:gosec,gomnd
}

// ArrayShuffle shuffles an array.
func ArrayShuffle[T any](arr []T) {
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
}

const bufSize = 4 * 1024

func readAll(r io.Reader) ([]byte, int64, error) {
	var buf bytes.Buffer
	chunk := make([]byte, bufSize) // 4 KB のバッファ
	var total int64

	for {
		n, err := r.Read(chunk)
		if n > 0 {
			// バッファに読み取りデータを書き込む
			written, writeErr := buf.Write(chunk[:n])
			if writeErr != nil {
				return nil, total, fmt.Errorf("failed to write to buffer: %w", writeErr)
			}
			total += int64(written)
		}
		if err != nil {
			if err == io.EOF {
				break // ファイルの終わりに達した
			}
			return nil, total, fmt.Errorf("failed to read from reader: %w", err)
		}
	}

	return buf.Bytes(), total, nil
}
