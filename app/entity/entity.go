// Package entity エンティティ関連の構造体を定義するパッケージ。
package entity

import "time"

// Entity エンティティのインターフェース。
type Entity any

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

// String 文字列型。
type String struct {
	String string
	Valid  bool
}

// Float 浮動小数点数型。
type Float struct {
	Float64 float64
	Status  Status
}

// UUID UUID型。
type UUID struct {
	Bytes  [16]byte
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
	Status           Status
	InfinityModifier InfinityModifier
}
