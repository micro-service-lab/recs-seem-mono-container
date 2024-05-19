// Package entity エンティティ関連の構造体を定義するパッケージ。
package entity

import "time"

// Entity エンティティのインターフェース。
type Entity any

// NullableEntity NULL許容エンティティ。
type NullableEntity[T Entity] struct {
	Entity T
	Valid  bool
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
	// Infinity 無限大。
	Infinity InfinityModifier = 1
	// None 無限大なし。
	None InfinityModifier = 0
	// NegativeInfinity 負の無限大。
	NegativeInfinity InfinityModifier = -Infinity
)

// Int 整数型。
type Int struct {
	Int64 int64
	Valid bool
}

// String 文字列型。
type String struct {
	String string
	Valid  bool
}

// Float 浮動小数点数型。
type Float struct {
	Float64 float64
	Valid   bool
}

// UUID UUID型。
type UUID struct {
	Bytes [16]byte
	Valid bool
}

// UUIDs UUIDのスライス型。
type UUIDs struct {
	UUIDs  []UUID
	Status Status
}

// Date 日付型。
type Date struct {
	Time             time.Time
	Status           Status
	InfinityModifier InfinityModifier
}

// Timestamptz タイムスタンプ型。
type Timestamptz struct {
	Time             time.Time
	InfinityModifier InfinityModifier
	Valid            bool
}

// OID オブジェクトID型。
type OID uint32
