package store_test

import (
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/testutils"
)

type CurDir string

const (
	Next CurDir = "next"
	Prev CurDir = "prev"
)

func TestGetCursorData(t *testing.T) {
	type want struct {
		rspFile string
	}
	tests := map[string]struct {
		reqTurn       int
		order         DummyStatusOrderMethod
		limit         int32
		cursorRewrite []bool   // 次回のカーソルの誤りの書き換えを行うか(reqTurn - 1個)
		dirs          []CurDir // カーソルの方向(reqTurn - 1個)
		want          want
	}{
		"simple chunk": {
			reqTurn:       5,
			order:         DummyStatusOrderMethodDefault,
			limit:         2,
			cursorRewrite: []bool{false, false, false, false},
			dirs:          []CurDir{Next, Next, Next, Next},
			want: want{
				rspFile: "testdata/pagination/simple_chunk.json.golden",
			},
		},
		"simple chunk with remain": {
			reqTurn:       4,
			order:         DummyStatusOrderMethodDefault,
			limit:         3,
			cursorRewrite: []bool{false, false, false},
			dirs:          []CurDir{Next, Next, Next, Next},
			want: want{
				rspFile: "testdata/pagination/simple_chunk_with_remain.json.golden",
			},
		},
		"over access(invalid)": {
			reqTurn:       5,
			order:         DummyStatusOrderMethodDefault,
			limit:         3,
			cursorRewrite: []bool{false, false, false, false},
			dirs:          []CurDir{Next, Next, Next, Next},
			want: want{
				rspFile: "testdata/pagination/over_access.json.golden",
			},
		},
		"over page(Over Access)": {
			reqTurn:       2,
			order:         DummyStatusOrderMethodDefault,
			limit:         100,
			cursorRewrite: []bool{false, false},
			dirs:          []CurDir{Next, Next},
			want: want{
				rspFile: "testdata/pagination/over_page.json.golden",
			},
		},
		"num ordered": {
			reqTurn:       10,
			order:         DummyStatusOrderMethodAge,
			limit:         1,
			cursorRewrite: []bool{false, false, false, false, false, false, false, false, false},
			dirs:          []CurDir{Next, Next, Next, Next, Next, Next, Next, Next, Next, Next},
			want: want{
				rspFile: "testdata/pagination/num_order.json.golden",
			},
		},
		"reverse num ordered": {
			reqTurn:       10,
			order:         DummyStatusOrderMethodReverseAge,
			limit:         1,
			cursorRewrite: []bool{false, false, false, false, false, false, false, false, false},
			dirs:          []CurDir{Next, Next, Next, Next, Next, Next, Next, Next, Next, Next},
			want: want{
				rspFile: "testdata/pagination/reverse_num_order.json.golden",
			},
		},
		"strings ordered": {
			reqTurn:       10,
			order:         DummyStatusOrderMethodName,
			limit:         1,
			cursorRewrite: []bool{false, false, false, false, false, false, false, false, false},
			dirs:          []CurDir{Next, Next, Next, Next, Next, Next, Next, Next, Next, Next},
			want: want{
				rspFile: "testdata/pagination/string_order.json.golden",
			},
		},
		"reverse strings ordered": {
			reqTurn:       10,
			order:         DummyStatusOrderMethodReverseName,
			limit:         1,
			cursorRewrite: []bool{false, false, false, false, false, false, false, false, false},
			dirs:          []CurDir{Next, Next, Next, Next, Next, Next, Next, Next, Next, Next},
			want: want{
				rspFile: "testdata/pagination/reverse_string_order.json.golden",
			},
		},
		"time ordered": {
			reqTurn:       10,
			order:         DummyStatusOrderMethodLastLogin,
			limit:         1,
			cursorRewrite: []bool{false, false, false, false, false, false, false, false, false},
			dirs:          []CurDir{Next, Next, Next, Next, Next, Next, Next, Next, Next, Next},
			want: want{
				rspFile: "testdata/pagination/time_order.json.golden",
			},
		},
		"reverse time ordered": {
			reqTurn:       10,
			order:         DummyStatusOrderMethodReverseLastLogin,
			limit:         1,
			cursorRewrite: []bool{false, false, false, false, false, false, false, false, false},
			dirs:          []CurDir{Next, Next, Next, Next, Next, Next, Next, Next, Next, Next},
			want: want{
				rspFile: "testdata/pagination/reverse_time_order.json.golden",
			},
		},
		"rewrite cursor": {
			reqTurn:       5,
			order:         DummyStatusOrderMethodDefault,
			limit:         2,
			cursorRewrite: []bool{false, false, true, true},
			dirs:          []CurDir{Next, Next, Next, Next},
			want: want{
				rspFile: "testdata/pagination/rewrite_cursor.json.golden",
			},
		},
		"random access 1": {
			reqTurn:       5,
			order:         DummyStatusOrderMethodDefault,
			limit:         2,
			cursorRewrite: []bool{false, false, false, false},
			dirs:          []CurDir{Next, Next, Prev, Prev},
			want: want{
				rspFile: "testdata/pagination/random_1.json.golden",
			},
		},
		"random access 2": {
			reqTurn:       5,
			order:         DummyStatusOrderMethodDefault,
			limit:         2,
			cursorRewrite: []bool{false, false, false, false},
			dirs:          []CurDir{Next, Next, Prev, Next},
			want: want{
				rspFile: "testdata/pagination/random_2.json.golden",
			},
		},
		"random access 3(Over Access)": {
			reqTurn:       5,
			order:         DummyStatusOrderMethodDefault,
			limit:         2,
			cursorRewrite: []bool{true, false, false, false},
			dirs:          []CurDir{Next, Next, Prev, Prev},
			want: want{
				rspFile: "testdata/pagination/random_3.json.golden",
			},
		},
		"random access 4": {
			reqTurn:       5,
			order:         DummyStatusOrderMethodLastLogin,
			limit:         2,
			cursorRewrite: []bool{false, false, false, false},
			dirs:          []CurDir{Next, Next, Prev, Prev},
			want: want{
				rspFile: "testdata/pagination/random_4.json.golden",
			},
		},
		"random access 5": {
			reqTurn:       5,
			order:         DummyStatusOrderMethodLastLogin,
			limit:         2,
			cursorRewrite: []bool{false, false, false, false},
			dirs:          []CurDir{Next, Next, Prev, Next},
			want: want{
				rspFile: "testdata/pagination/random_5.json.golden",
			},
		},
		"random access 6(Over Access)": {
			reqTurn:       5,
			order:         DummyStatusOrderMethodLastLogin,
			limit:         2,
			cursorRewrite: []bool{true, false, false, false},
			dirs:          []CurDir{Next, Next, Prev, Prev},
			want: want{
				rspFile: "testdata/pagination/random_6.json.golden",
			},
		},
		"random access 7": {
			reqTurn:       5,
			order:         DummyStatusOrderMethodReverseLastLogin,
			limit:         2,
			cursorRewrite: []bool{false, false, false, false},
			dirs:          []CurDir{Next, Next, Prev, Prev},
			want: want{
				rspFile: "testdata/pagination/random_7.json.golden",
			},
		},
		"random access 8": {
			reqTurn:       5,
			order:         DummyStatusOrderMethodReverseLastLogin,
			limit:         2,
			cursorRewrite: []bool{false, false, false, false},
			dirs:          []CurDir{Next, Next, Prev, Next},
			want: want{
				rspFile: "testdata/pagination/random_8.json.golden",
			},
		},
		"random access 9(Over Access)": {
			reqTurn:       5,
			order:         DummyStatusOrderMethodReverseLastLogin,
			limit:         2,
			cursorRewrite: []bool{true, false, false, false},
			dirs:          []CurDir{Next, Next, Prev, Prev},
			want: want{
				rspFile: "testdata/pagination/random_9.json.golden",
			},
		},
		"random access 10": {
			reqTurn:       5,
			order:         DummyStatusOrderMethodDefault,
			limit:         3,
			cursorRewrite: []bool{false, false, false, false},
			dirs:          []CurDir{Next, Next, Prev, Prev},
			want: want{
				rspFile: "testdata/pagination/random_10.json.golden",
			},
		},
		"random access 11": {
			reqTurn:       5,
			order:         DummyStatusOrderMethodDefault,
			limit:         3,
			cursorRewrite: []bool{false, false, false, false},
			dirs:          []CurDir{Next, Next, Prev, Next},
			want: want{
				rspFile: "testdata/pagination/random_11.json.golden",
			},
		},
		"random access 12(Over Access)": {
			reqTurn:       5,
			order:         DummyStatusOrderMethodDefault,
			limit:         3,
			cursorRewrite: []bool{true, false, false, false},
			dirs:          []CurDir{Next, Next, Prev, Prev},
			want: want{
				rspFile: "testdata/pagination/random_12.json.golden",
			},
		},
		"random access 13": {
			reqTurn:       5,
			order:         DummyStatusOrderMethodDefault,
			limit:         3,
			cursorRewrite: []bool{false, false, false, false},
			dirs:          []CurDir{Next, Next, Next, Prev},
			want: want{
				rspFile: "testdata/pagination/random_13.json.golden",
			},
		},
	}
	reqData := testutils.LoadFile(t, "testdata/pagination/in.json.golden")
	for n, tt := range tests {
		ttt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			result := make([]DummyStatuses, ttt.reqTurn)
			// DummyStatusのデータを作成
			var dummyStatuses DummyStatuses
			err := json.Unmarshal(reqData, &dummyStatuses)
			if err != nil {
				t.Fatalf("failed to unmarshal request data: %v", err)
			}
			runQCPFunc := func(subCursor, orderMethod string,
				limit int32, cursorDir string, cursor int32, subCursorValue any,
			) ([]DummyStatus, error) {
				var nameCursor string
				var ageCursor int
				var lastLoginCursor time.Time
				var ok bool
				var err error
				switch subCursor {
				case DummyStatusNameCursorKey:
					nameCursor, ok = subCursorValue.(string)
					if !ok {
						nameCursor = ""
					}
				case DummyStatusAgeCursorKey:
					ac, ok := subCursorValue.(float64)
					ageCursor = int(ac)
					if !ok {
						ageCursor = 0
					}
				case DummyStatusLastLoginCursorKey:
					cv, ok := subCursorValue.(string)
					lastLoginCursor, err = time.Parse(time.RFC3339, cv)
					if !ok || err != nil {
						lastLoginCursor = time.Time{}
					}
				}
				r := dummyStatuses.RetrieveWithCursor(
					t,
					DummyStatusOrderMethod(orderMethod),
					cursor,
					limit,
					cursorDir,
					nameCursor,
					ageCursor,
					lastLoginCursor,
				)
				return r, nil
			}

			runQNPFunc := func(orderMethod string, limit, offset int32) ([]DummyStatus, error) {
				r := dummyStatuses.RetrieveWithNumbered(
					t,
					DummyStatusOrderMethod(orderMethod),
					limit,
					offset,
				)
				return r, nil
			}
			selector := func(subCursor string, e DummyStatus) (entity.Int, any) {
				switch subCursor {
				case DummyStatusDefaultCursorKey:
					return entity.Int{Int64: int64(e.Pkey)}, nil
				case DummyStatusNameCursorKey:
					return entity.Int{Int64: int64(e.Pkey)}, e.Name
				case DummyStatusAgeCursorKey:
					return entity.Int{Int64: int64(e.Pkey)}, e.Age
				case DummyStatusLastLoginCursorKey:
					return entity.Int{Int64: int64(e.Pkey)}, e.LastLogin
				}
				return entity.Int{Int64: int64(e.Pkey)}, nil
			}

			cursor := ""

			for i := 0; i < ttt.reqTurn; i++ {
				res, pi, err := store.GetCursorData[DummyStatus](
					cursor,
					ttt.order,
					ttt.limit,
					runQCPFunc,
					runQNPFunc,
					selector,
				)
				if err != nil {
					t.Errorf("failed to get cursor data: %v", err)
				}
				result[i] = res
				if i < len(ttt.dirs) {
					if ttt.dirs[i] == Next {
						cursor = pi.NextCursor
					} else {
						cursor = pi.PrevCursor
					}
				}
				if i < len(ttt.cursorRewrite) && ttt.cursorRewrite[i] {
					cursor = cursor + "invalid"
				}
			}
			got, err := json.Marshal(result)
			if err != nil {
				t.Fatalf("failed to marshal response: %v", err)
			}

			testutils.AssertJSON(t, testutils.LoadFile(t, ttt.want.rspFile), got)
		})
	}
}

// DummyStatus ダミーステータス。
type DummyStatus struct {
	Pkey      int32     `json:"pkey"`
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	LastLogin time.Time `json:"lastLoginAt"`
	IsActive  bool      `json:"isActive"`
}

// DummyStatuses ダミーステータスのスライス。
type DummyStatuses []DummyStatus

// RetrieveWithNumbered 番号つきページネーションでデータを取得する。
func (d DummyStatuses) RetrieveWithNumbered(
	t *testing.T,
	method DummyStatusOrderMethod,
	limit int32,
	offset int32,
) DummyStatuses {
	t.Helper()

	var result DummyStatuses
	result = append(result, d...)
	switch method {
	case DummyStatusOrderMethodName:
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Name < result[j].Name
		})
	case DummyStatusOrderMethodReverseName:
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Name > result[j].Name
		})
	case DummyStatusOrderMethodLastLogin:
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].LastLogin.Before(result[j].LastLogin)
		})
	case DummyStatusOrderMethodReverseLastLogin:
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].LastLogin.After(result[j].LastLogin)
		})
	case DummyStatusOrderMethodAge:
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Age < result[j].Age
		})
	case DummyStatusOrderMethodReverseAge:
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Age > result[j].Age
		})
	case DummyStatusOrderMethodDefault:
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Pkey < result[j].Pkey
		})
	}
	if int(offset) < len(result) {
		result = result[offset:]
	}
	if int(limit) < len(result) {
		result = result[:limit]
	}
	return result
}

// RetrieveWithCursor カーソルページネーションでデータを取得する。
func (d DummyStatuses) RetrieveWithCursor(
	t *testing.T,
	method DummyStatusOrderMethod,
	cursor, limit int32,
	curDir, nameCur string,
	ageCur int,
	lastLoginCur time.Time,
) DummyStatuses {
	t.Helper()
	var result DummyStatuses
	result = append(result, d...)

	// where句の条件を満たすデータを取得
	var filtered DummyStatuses
	for _, v := range result {
		switch curDir {
		case "next":
			switch method {
			case DummyStatusOrderMethodName:
				if v.Name > nameCur || (v.Name == nameCur && v.Pkey > cursor) {
					filtered = append(filtered, v)
				}
			case DummyStatusOrderMethodReverseName:
				if v.Name < nameCur || (v.Name == nameCur && v.Pkey > cursor) {
					filtered = append(filtered, v)
				}
			case DummyStatusOrderMethodLastLogin:
				if v.LastLogin.After(lastLoginCur) || (v.LastLogin.Equal(lastLoginCur) && v.Pkey > cursor) {
					filtered = append(filtered, v)
				}
			case DummyStatusOrderMethodReverseLastLogin:
				if v.LastLogin.Before(lastLoginCur) || (v.LastLogin.Equal(lastLoginCur) && v.Pkey > cursor) {
					filtered = append(filtered, v)
				}
			case DummyStatusOrderMethodAge:
				if v.Age > ageCur || (v.Age == ageCur && v.Pkey > cursor) {
					filtered = append(filtered, v)
				}
			case DummyStatusOrderMethodReverseAge:
				if v.Age < ageCur || (v.Age == ageCur && v.Pkey > cursor) {
					filtered = append(filtered, v)
				}
			case DummyStatusOrderMethodDefault:
				if v.Pkey > cursor {
					filtered = append(filtered, v)
				}
			}
		case "prev":
			switch method {
			case DummyStatusOrderMethodName:
				if v.Name < nameCur || (v.Name == nameCur && v.Pkey < cursor) {
					filtered = append(filtered, v)
				}
			case DummyStatusOrderMethodReverseName:
				if v.Name > nameCur || (v.Name == nameCur && v.Pkey < cursor) {
					filtered = append(filtered, v)
				}
			case DummyStatusOrderMethodLastLogin:
				if v.LastLogin.Before(lastLoginCur) || (v.LastLogin.Equal(lastLoginCur) && v.Pkey < cursor) {
					filtered = append(filtered, v)
				}
			case DummyStatusOrderMethodReverseLastLogin:
				if v.LastLogin.After(lastLoginCur) || (v.LastLogin.Equal(lastLoginCur) && v.Pkey < cursor) {
					filtered = append(filtered, v)
				}
			case DummyStatusOrderMethodAge:
				if v.Age < ageCur || (v.Age == ageCur && v.Pkey < cursor) {
					filtered = append(filtered, v)
				}
			case DummyStatusOrderMethodReverseAge:
				if v.Age > ageCur || (v.Age == ageCur && v.Pkey < cursor) {
					filtered = append(filtered, v)
				}
			case DummyStatusOrderMethodDefault:
				if v.Pkey < cursor {
					filtered = append(filtered, v)
				}
			}
		}
	}

	// ORDER BY
	switch curDir {
	case "next":
		sort.SliceStable(filtered, func(i, j int) bool {
			return filtered[i].Pkey < filtered[j].Pkey
		})
		switch method {
		case DummyStatusOrderMethodName:
			sort.SliceStable(filtered, func(i, j int) bool {
				return filtered[i].Name < filtered[j].Name
			})
		case DummyStatusOrderMethodReverseName:
			sort.SliceStable(filtered, func(i, j int) bool {
				return filtered[i].Name > filtered[j].Name
			})
		case DummyStatusOrderMethodLastLogin:
			sort.SliceStable(filtered, func(i, j int) bool {
				return filtered[i].LastLogin.Before(filtered[j].LastLogin)
			})
		case DummyStatusOrderMethodReverseLastLogin:
			sort.SliceStable(filtered, func(i, j int) bool {
				return filtered[i].LastLogin.After(filtered[j].LastLogin)
			})
		case DummyStatusOrderMethodAge:
			sort.SliceStable(filtered, func(i, j int) bool {
				return filtered[i].Age < filtered[j].Age
			})
		case DummyStatusOrderMethodReverseAge:
			sort.SliceStable(filtered, func(i, j int) bool {
				return filtered[i].Age > filtered[j].Age
			})
		case DummyStatusOrderMethodDefault:
		}
	case "prev":
		sort.SliceStable(filtered, func(i, j int) bool {
			return filtered[i].Pkey > filtered[j].Pkey
		})
		switch method {
		case DummyStatusOrderMethodName:
			sort.SliceStable(filtered, func(i, j int) bool {
				return filtered[i].Name > filtered[j].Name
			})
		case DummyStatusOrderMethodReverseName:
			sort.SliceStable(filtered, func(i, j int) bool {
				return filtered[i].Name < filtered[j].Name
			})
		case DummyStatusOrderMethodLastLogin:
			sort.SliceStable(filtered, func(i, j int) bool {
				return filtered[i].LastLogin.After(filtered[j].LastLogin)
			})
		case DummyStatusOrderMethodReverseLastLogin:
			sort.SliceStable(filtered, func(i, j int) bool {
				return filtered[i].LastLogin.Before(filtered[j].LastLogin)
			})
		case DummyStatusOrderMethodAge:
			sort.SliceStable(filtered, func(i, j int) bool {
				return filtered[i].Age > filtered[j].Age
			})
		case DummyStatusOrderMethodReverseAge:
			sort.SliceStable(filtered, func(i, j int) bool {
				return filtered[i].Age < filtered[j].Age
			})
		case DummyStatusOrderMethodDefault:
		}
	}

	// LIMIT
	if int(limit) < len(filtered) {
		filtered = filtered[:limit]
	}
	return filtered
}

// DummyStatusOrderMethod 出席ステータスの並び替え方法。
type DummyStatusOrderMethod string

// ParseDummyStatusOrderMethod は並び替え方法をパースする。
func ParseDummyStatusOrderMethod(v string) (any, error) {
	if v == "" {
		return DummyStatusOrderMethodDefault, nil
	}
	switch v {
	case string(DummyStatusOrderMethodDefault):
		return DummyStatusOrderMethodDefault, nil
	case string(DummyStatusOrderMethodName):
		return DummyStatusOrderMethodName, nil
	case string(DummyStatusOrderMethodReverseName):
		return DummyStatusOrderMethodReverseName, nil
	case string(DummyStatusOrderMethodLastLogin):
		return DummyStatusOrderMethodLastLogin, nil
	case string(DummyStatusOrderMethodReverseLastLogin):
		return DummyStatusOrderMethodReverseLastLogin, nil
	case string(DummyStatusOrderMethodAge):
		return DummyStatusOrderMethodAge, nil
	case string(DummyStatusOrderMethodReverseAge):
		return DummyStatusOrderMethodReverseAge, nil
	default:
		return DummyStatusOrderMethodDefault, nil
	}
}

const (
	// DummyStatusDefaultCursorKey はデフォルトカーソルキー。
	DummyStatusDefaultCursorKey = "default"
	// DummyStatusNameCursorKey は名前カーソルキー。
	DummyStatusNameCursorKey = "name"
	// DummyStatusLastLoginCursorKey は最終ログインカーソルキー。
	DummyStatusLastLoginCursorKey = "last_login"
	// DummyStatusAgeCursorKey は年齢カーソルキー。
	DummyStatusAgeCursorKey = "age"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m DummyStatusOrderMethod) GetCursorKeyName() string {
	switch m {
	case DummyStatusOrderMethodDefault:
		return DummyStatusDefaultCursorKey
	case DummyStatusOrderMethodName:
		return DummyStatusNameCursorKey
	case DummyStatusOrderMethodReverseName:
		return DummyStatusNameCursorKey
	case DummyStatusOrderMethodLastLogin:
		return DummyStatusLastLoginCursorKey
	case DummyStatusOrderMethodReverseLastLogin:
		return DummyStatusLastLoginCursorKey
	case DummyStatusOrderMethodAge:
		return DummyStatusAgeCursorKey
	case DummyStatusOrderMethodReverseAge:
		return DummyStatusAgeCursorKey
	default:
		return DummyStatusDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m DummyStatusOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// DummyStatusOrderMethodDefault はデフォルト。
	DummyStatusOrderMethodDefault DummyStatusOrderMethod = "default"
	// DummyStatusOrderMethodName は名前順。
	DummyStatusOrderMethodName DummyStatusOrderMethod = "name"
	// DummyStatusOrderMethodReverseName は名前逆順。
	DummyStatusOrderMethodReverseName DummyStatusOrderMethod = "r_name"
	// DummyStatusOrderMethodLastLogin は最終ログイン順。
	DummyStatusOrderMethodLastLogin DummyStatusOrderMethod = "last_login"
	// DummyStatusOrderMethodReverseLastLogin は最終ログイン逆順。
	DummyStatusOrderMethodReverseLastLogin DummyStatusOrderMethod = "r_last_login"
	// DummyStatusOrderMethodAge は年齢順。
	DummyStatusOrderMethodAge DummyStatusOrderMethod = "age"
	// DummyStatusOrderMethodReverseAge は年齢逆順。
	DummyStatusOrderMethodReverseAge DummyStatusOrderMethod = "r_age"
)
