package handler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
)

// DownloadObject is a handler for downloading object.
type DownloadObject struct {
	Service service.ManagerInterface
}

func (h *DownloadObject) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	id := uuid.MustParse(chi.URLParam(r, "attachable_item_id"))

	reader, alias, err := h.Service.DownloadAttachableItem(ctx, authUser.MemberID, id)
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("failed to close reader: %v", err)
		}
	}()
	data, _, err := readAll(reader)
	if err != nil {
		log.Printf("failed to read all data: %v", err)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", alias))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))

	if _, err := w.Write(data); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func readAll(r io.Reader) ([]byte, int64, error) {
	const bufSize = 4 * 1024

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
