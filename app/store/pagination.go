package store

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
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
	Valid bool  `json:"valid"`
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
type Cursor struct {
	Valid            bool   `json:"valid"`
	CursorID         int64  `json:"id"`
	CursorPointsNext bool   `json:"points_next"`
	SubCursorName    string `json:"sub_cursor_name"`
	SubCursor        any    `json:"sub_cursor"`
}

// CreateCursor カーソルを生成する。
func CreateCursor(id int64, pointsNext bool, name string, value any) Cursor {
	c := Cursor{
		Valid:            true,
		CursorID:         id,
		CursorPointsNext: pointsNext,
		SubCursorName:    name,
		SubCursor:        value,
	}
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
	if !cursor.Valid {
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
		return Cursor{}, fmt.Errorf("failed to decode cursor: %w", err)
	}

	var cur Cursor
	if err := json.Unmarshal(decodedCursor, &cur); err != nil {
		return Cursor{}, fmt.Errorf("failed to unmarshal cursor: %w", err)
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
			pagination = GeneratePager(nextCur, Cursor{})
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

// GetCursorData カーソルデータを取得する。
func GetCursorData[T any](
	cursor string,
	order parameter.OrderMethod,
	limit int32,
	runQWithCursor RunQueryWithCursorParamsFunc[T],
	runQWithNumbered RunQueryWithNumberedParamsFunc[T],
	selector CursorIDAndValueSelector[T],
) ([]T, CursorPaginationAttribute, error) {
	var err error
	isFirst := cursor == "" // 初回のリクエストかどうか
	pointsNext := false     // ページネーションの方向(true: 次データ, false: 前データ)
	subCursor := order.GetCursorKeyName()
	// カーソルのデコード+チェック
	var decodedCursor Cursor
	var cursorData any
	var data []T
	if !isFirst {
		cursorCheck := func(cur string) bool {
			decodedCursor, err = DecodeCursor(cur)
			if err != nil {
				return false
			}
			if decodedCursor.SubCursorName != subCursor {
				return false
			}
			cursorData = decodedCursor.SubCursor
			return true
		}
		if !cursorCheck(cursor) {
			isFirst = true
		}
	}

	if !isFirst {
		// 今回の指定カーソルの方向を引き継ぐ
		pointsNext = decodedCursor.CursorPointsNext
		var cursorDir string
		if pointsNext {
			cursorDir = "next"
		} else {
			cursorDir = "prev"
		}
		ID := decodedCursor.CursorID
		data, err = runQWithCursor(subCursor, order.GetStringValue(), limit+1, cursorDir, int32(ID), cursorData)
		if err != nil {
			return nil, CursorPaginationAttribute{}, fmt.Errorf("failed to run query with cursor params: %w", err)
		}
	} else {
		data, err = runQWithNumbered(order.GetStringValue(), limit+1, 0)
		if err != nil {
			return nil, CursorPaginationAttribute{}, fmt.Errorf("failed to run query with numbered params: %w", err)
		}
	}

	// dataの要素数が0である場合
	if len(data) == 0 {
		return nil, CursorPaginationAttribute{}, ErrDataNoRecord
	}
	hasPagination := len(data) > int(limit)
	if hasPagination {
		data = data[:limit]
	}
	eLen := len(data)

	var firstValue, lastValue any
	lastIndex := eLen - 1
	if lastIndex < 0 {
		lastIndex = 0
	}
	firstID, firstValue := selector(subCursor, data[0])
	lastID, lastValue := selector(subCursor, data[lastIndex])

	firstData := CursorData{
		ID:    firstID,
		Name:  subCursor,
		Value: firstValue,
	}
	lastData := CursorData{
		ID:    lastID,
		Name:  subCursor,
		Value: lastValue,
	}
	var pageInfo CursorPaginationAttribute
	if pointsNext || isFirst {
		pageInfo = CalculatePagination(isFirst, hasPagination, pointsNext, firstData, lastData)
	} else {
		pageInfo = CalculatePagination(isFirst, hasPagination, pointsNext, lastData, firstData)
	}

	return data, pageInfo, nil
}
