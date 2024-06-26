package cors

import (
	"strings"
	"testing"
)

func TestWildcard(t *testing.T) {
	w := wildcard{"foo", "bar"}
	if !w.match("foobar") {
		t.Error("foo*bar should match foobar")
	}
	if !w.match("foobazbar") {
		t.Error("foo*bar should match foobazbar")
	}
	if w.match("foobaz") {
		t.Error("foo*bar should not match foobaz")
	}

	w = wildcard{"foo", "oof"}
	if w.match("foof") {
		t.Error("foo*oof should not match foof")
	}
}

func TestConvert(t *testing.T) {
	s := convert([]string{"A", "b", "C"}, strings.ToLower)
	e := []string{"a", "b", "c"}
	if s[0] != e[0] || s[1] != e[1] || s[2] != e[2] {
		t.Errorf("%v != %v", s, e)
	}
}

func TestParseHeaderList(t *testing.T) {
	h := parseHeaderList(
		"header, second-header, THIRD-HEADER, Numb3r3d-H34d3r, Header_with_underscore Header.with.full.stop")
	e := []string{
		"Header", "Second-Header", "Third-Header", "Numb3r3d-H34d3r", "Header_with_underscore", "Header.with.full.stop",
	}
	if h[0] != e[0] || h[1] != e[1] || h[2] != e[2] || h[3] != e[3] || h[4] != e[4] || h[5] != e[5] {
		t.Errorf("%v != %v", h, e)
	}
}

func TestParseHeaderListEmpty(t *testing.T) {
	if len(parseHeaderList("")) != 0 {
		t.Error("should be empty slice")
	}
	if len(parseHeaderList(" , ")) != 0 {
		t.Error("should be empty slice")
	}
}

func BenchmarkParseHeaderList(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseHeaderList("header, second-header, THIRD-HEADER")
	}
}

func BenchmarkParseHeaderListSingle(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseHeaderList("header")
	}
}

func BenchmarkParseHeaderListNormalized(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseHeaderList("Header1, Header2, Third-Header")
	}
}

func BenchmarkWildcard(b *testing.B) {
	w := wildcard{"foo", "bar"}
	b.Run("match", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			w.match("foobazbar")
		}
	})
	b.Run("too short", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			w.match("fobar")
		}
	})
}
