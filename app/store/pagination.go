package store

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

const (
	// CursorID カーソルID。
	CursorID = "id"
	// CursorPointsNext 次のカーソル。
	CursorPointsNext = "points_next"
)

// WithCountParam カウントパラメータを表す構造体。
type WithCountParam struct {
	Valid bool
}

// WithCountAttribute カウント属性を表す構造体。
type WithCountAttribute struct {
	Count int64 `json:"count"`
}

// NumberedPaginationParam ページネーションのパラメータを表す構造体。
type NumberedPaginationParam struct {
	Valid  bool
	Offset entity.Int
	Limit  entity.Int
}

// CursorPaginationParam カーソルページネーションのパラメータを表す構造体。
type CursorPaginationParam struct {
	Valid  bool
	Cursor string
	Limit  entity.Int
}

// CursorPaginationAttribute カーソルページネーションのレスポンスを表す構造体。
type CursorPaginationAttribute struct {
	NextCursor string `json:"next_cursor"`
	PrevCursor string `json:"prev_cursor"`
}

// Cursor カーソルを表す構造体。
type Cursor map[string]any

// CreateCursor カーソルを生成する。
func CreateCursor(id int64, pointsNext bool, name string, value any) Cursor {
	c := Cursor{
		CursorID:         id,
		CursorPointsNext: pointsNext,
	}
	c[name] = value
	return c
}

// GeneratePager ページネーション情報を生成する。
func GeneratePager(next, prev Cursor) CursorPaginationAttribute {
	return CursorPaginationAttribute{
		NextCursor: encodeCursor(next),
		PrevCursor: encodeCursor(prev),
	}
}

// encodeCursor カーソルをエンコードする。
func encodeCursor(cursor Cursor) string {
	if len(cursor) == 0 {
		return ""
	}
	serializedCursor, err := json.Marshal(cursor)
	if err != nil {
		return ""
	}
	encodedCursor := base64.StdEncoding.EncodeToString(serializedCursor)
	return encodedCursor
}

// DecodeCursor カーソルをデコードする。
func DecodeCursor(cursor string) (Cursor, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cursor: %w", err)
	}

	var cur Cursor
	if err := json.Unmarshal(decodedCursor, &cur); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cursor: %w", err)
	}
	return cur, nil
}

// CursorData カーソルデータを表す構造体。
type CursorData struct {
	ID    entity.Int
	Name  string
	Value any
}

// CalculatePagination ページネーション情報を計算する。
func CalculatePagination(
	isFirstPage, hasPagination, pointsNext bool, firstData, lastData CursorData,
) CursorPaginationAttribute {
	pagination := CursorPaginationAttribute{}
	nextCur := Cursor{}
	prevCur := Cursor{}
	if isFirstPage {
		if hasPagination {
			nextCur := CreateCursor(lastData.ID.Int64, true, lastData.Name, lastData.Value)
			pagination = GeneratePager(nextCur, nil)
		}
	} else {
		if pointsNext {
			// if pointing next, it always has prev but it might not have next
			if hasPagination {
				nextCur = CreateCursor(lastData.ID.Int64, true, lastData.Name, lastData.Value)
			}
			prevCur = CreateCursor(firstData.ID.Int64, false, firstData.Name, firstData.Value)
			pagination = GeneratePager(nextCur, prevCur)
		} else {
			// this is case of prev, there will always be nest, but prev needs to be calculated
			nextCur = CreateCursor(lastData.ID.Int64, true, lastData.Name, lastData.Value)
			if hasPagination {
				prevCur = CreateCursor(firstData.ID.Int64, false, firstData.Name, firstData.Value)
			}
			pagination = GeneratePager(nextCur, prevCur)
		}
	}
	return pagination
}
