// Package entity エンティティ関連の構造体を定義するパッケージ。
//
//nolint:gomnd,stylecheck,revive,wrapcheck
package entity

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Entity エンティティのインターフェース。
type Entity any

// NullableEntity NULL許容エンティティ。
type NullableEntity[T Entity] struct {
	Entity T
	Valid  bool
}

func (ne NullableEntity[T]) MarshalJSON() ([]byte, error) {
	if !ne.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(ne.Entity)
}

func (ne *NullableEntity[T]) UnmarshalJSON(b []byte) error {
	var n *T
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}

	if n == nil {
		*ne = NullableEntity[T]{Valid: false}
	} else {
		*ne = NullableEntity[T]{Entity: *n, Valid: true}
	}

	return nil
}

// Status ステータスを表す型。
type Status byte

const (
	// Undefined 未定義。
	Undefined Status = iota
	// Null 無効。
	Null
	// Present 有効。
	Present
)

// InfinityModifier 無限大の修飾子。
type InfinityModifier int8

const (
	Infinity         InfinityModifier = 1
	Finite           InfinityModifier = 0
	NegativeInfinity InfinityModifier = -Infinity
)

// Int 整数型。
type Int struct {
	Int64 int64
	Valid bool
}

// MarshalJSON 整数をJSON形式に変換する。
func (src Int) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(int64(src.Int64), 10)), nil
}

// UnmarshalJSON JSON形式を整数に変換する。
func (dst *Int) UnmarshalJSON(b []byte) error {
	var n *int64
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}

	if n == nil {
		*dst = Int{}
	} else {
		*dst = Int{Int64: *n, Valid: true}
	}

	return nil
}

// String 文字列型。
type String struct {
	String string
	Valid  bool
}

// MarshalJSON 文字列をJSON形式に変換する。
func (src String) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(src.String)
}

// UnmarshalJSON JSON形式を文字列に変換する。
func (dst *String) UnmarshalJSON(b []byte) error {
	var s *string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	if s == nil {
		*dst = String{}
	} else {
		*dst = String{String: *s, Valid: true}
	}

	return nil
}

// Float 浮動小数点数型。
type Float struct {
	Float64 float64
	Valid   bool
}

// MarshalJSON 浮動小数点数をJSON形式に変換する。
func (f Float) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(f.Float64)
}

// UnmarshalJSON JSON形式を浮動小数点数に変換する。
func (f *Float) UnmarshalJSON(b []byte) error {
	var n *float64
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}

	if n == nil {
		*f = Float{}
	} else {
		*f = Float{Float64: *n, Valid: true}
	}

	return nil
}

// UUID UUID型。
type UUID struct {
	Bytes [16]byte
	Valid bool
}

// parseUUID converts a string UUID in standard form to a byte array.
func parseUUID(src string) (dst [16]byte, err error) {
	switch len(src) {
	case 36:
		src = src[0:8] + src[9:13] + src[14:18] + src[19:23] + src[24:]
	case 32:
		// dashes already stripped, assume valid
	default:
		// assume invalid.
		return dst, fmt.Errorf("cannot parse UUID %v", src)
	}

	buf, err := hex.DecodeString(src)
	if err != nil {
		return dst, err
	}

	copy(dst[:], buf)
	return dst, err
}

// encodeUUID converts a uuid byte array to UUID standard string form.
func encodeUUID(src [16]byte) string {
	var buf [36]byte

	hex.Encode(buf[0:8], src[:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], src[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], src[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], src[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], src[10:])

	return string(buf[:])
}

// MarshalJSON UUIDをJSON形式に変換する。
func (src UUID) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}

	var buff bytes.Buffer
	buff.WriteByte('"')
	buff.WriteString(encodeUUID(src.Bytes))
	buff.WriteByte('"')
	return buff.Bytes(), nil
}

// UnmarshalJSON JSON形式をUUIDに変換する。
func (dst *UUID) UnmarshalJSON(src []byte) error {
	if bytes.Equal(src, []byte("null")) {
		*dst = UUID{}
		return nil
	}
	if len(src) != 38 {
		return fmt.Errorf("invalid length for UUID: %v", len(src))
	}
	buf, err := parseUUID(string(src[1 : len(src)-1]))
	if err != nil {
		return err
	}
	*dst = UUID{Bytes: buf, Valid: true}
	return nil
}

// Date 日付型。
type Date struct {
	Time             time.Time
	InfinityModifier InfinityModifier
	Valid            bool
}

func (src Date) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}

	var s string

	switch src.InfinityModifier {
	case Finite:
		s = src.Time.Format("2006-01-02")
	case Infinity:
		s = "infinity"
	case NegativeInfinity:
		s = "-infinity"
	}

	return json.Marshal(s)
}

func (dst *Date) UnmarshalJSON(b []byte) error {
	var s *string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	if s == nil {
		*dst = Date{}
		return nil
	}

	switch *s {
	case "infinity":
		*dst = Date{Valid: true, InfinityModifier: Infinity}
	case "-infinity":
		*dst = Date{Valid: true, InfinityModifier: -Infinity}
	default:
		t, err := time.ParseInLocation("2006-01-02", *s, time.UTC)
		if err != nil {
			return err
		}

		*dst = Date{Time: t, Valid: true}
	}

	return nil
}

// Timestamptz タイムスタンプ型。
type Timestamptz struct {
	Time             time.Time
	InfinityModifier InfinityModifier
	Valid            bool
}

func (tstz Timestamptz) MarshalJSON() ([]byte, error) {
	if !tstz.Valid {
		return []byte("null"), nil
	}

	var s string

	switch tstz.InfinityModifier {
	case Finite:
		s = tstz.Time.Format(time.RFC3339Nano)
	case Infinity:
		s = "infinity"
	case NegativeInfinity:
		s = "-infinity"
	}

	return json.Marshal(s)
}

func (tstz *Timestamptz) UnmarshalJSON(b []byte) error {
	var s *string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	if s == nil {
		*tstz = Timestamptz{}
		return nil
	}

	switch *s {
	case "infinity":
		*tstz = Timestamptz{Valid: true, InfinityModifier: Infinity}
	case "-infinity":
		*tstz = Timestamptz{Valid: true, InfinityModifier: -Infinity}
	default:
		// PostgreSQL uses ISO 8601 for to_json function and casting from a string to timestamptz
		tim, err := time.Parse(time.RFC3339Nano, *s)
		if err != nil {
			return err
		}

		*tstz = Timestamptz{Time: tim, Valid: true}
	}

	return nil
}

// OID オブジェクトID型。
type OID uint32
