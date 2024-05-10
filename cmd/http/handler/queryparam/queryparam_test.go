package queryparam_test

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/queryparam"
)

func TestParser_Parse(t *testing.T) {
	t.Parallel()
	type args struct {
		query string
	}
	type wants struct {
		ts testStruct
	}
	//nolint: lll
	cases := map[string]struct {
		args args
		want wants
	}{
		"valid query": {
			args: args{
				query: "bool=true&int=20&int8=30&int16=40&int32=50&int64=60&uint=70&uint8=80&uint16=90&uint32=100&uint64=110&str=foo&float32=1.1&float64=2.2",
			},
			want: wants{
				ts: testStruct{
					Bool:                   true,
					Int:                    20,
					Int8:                   30,
					Int16:                  40,
					Int32:                  50,
					Int64:                  60,
					Uint:                   70,
					Uint8:                  80,
					Uint16:                 90,
					Uint32:                 100,
					Uint64:                 110,
					Str:                    "foo",
					Float32:                1.1,
					Float64:                2.2,
					DefaultInt:             10,
					SliceInt:               []int{},
					DefaultSliceInt:        []int{1, 2, 3},
					SliceSeparatorOverride: []int{10, 20, 30},
					Custom:                 CustomType("no"),
					SliceCustom:            []CustomType{"a b", "b c", "c"},
					SliceCustomNoDefault:   []CustomType{},
					CustomStruct:           CustomStruct{First: "no", Last: "no"},
					Dive: Dive{
						SliceInt: []int{},
					},
				},
			},
		},
		"zero set value": {
			args: args{
				query: "",
			},
			want: wants{
				ts: testStruct{
					Bool:                   false,
					Int:                    0,
					Int8:                   0,
					Int16:                  0,
					Int32:                  0,
					Int64:                  0,
					Uint:                   0,
					Uint8:                  0,
					Uint16:                 0,
					Uint32:                 0,
					Uint64:                 0,
					Str:                    "",
					Float32:                0,
					Float64:                0,
					DefaultInt:             10,
					SliceInt:               []int{},
					DefaultSliceInt:        []int{1, 2, 3},
					SliceSeparatorOverride: []int{10, 20, 30},
					Custom:                 CustomType("no"),
					SliceCustom:            []CustomType{"a b", "b c", "c"},
					SliceCustomNoDefault:   []CustomType{},
					CustomStruct:           CustomStruct{First: "no", Last: "no"},
					Dive: Dive{
						SliceInt: []int{},
					},
				},
			},
		},
		"default Parser": {
			args: args{
				query: "timeDuration=1h",
			},
			want: wants{
				ts: testStruct{
					TimeDuration:           time.Hour,
					DefaultInt:             10,
					SliceInt:               []int{},
					DefaultSliceInt:        []int{1, 2, 3},
					SliceSeparatorOverride: []int{10, 20, 30},
					Custom:                 CustomType("no"),
					SliceCustom:            []CustomType{"a b", "b c", "c"},
					SliceCustomNoDefault:   []CustomType{},
					CustomStruct:           CustomStruct{First: "no", Last: "no"},
					Dive: Dive{
						SliceInt: []int{},
					},
				},
			},
		},
		"not work parser": {
			args: args{
				query: "int=foo",
			},
			want: wants{
				ts: testStruct{
					DefaultInt:             10,
					SliceInt:               []int{},
					DefaultSliceInt:        []int{1, 2, 3},
					SliceSeparatorOverride: []int{10, 20, 30},
					Custom:                 CustomType("no"),
					SliceCustom:            []CustomType{"a b", "b c", "c"},
					SliceCustomNoDefault:   []CustomType{},
					CustomStruct:           CustomStruct{First: "no", Last: "no"},
					Dive: Dive{
						SliceInt: []int{},
					},
				},
			},
		},
		"slice": {
			args: args{
				query: "sliceInt=1&sliceInt=2&sliceInt=3",
			},
			want: wants{
				ts: testStruct{
					SliceInt:               []int{1, 2, 3},
					DefaultInt:             10,
					DefaultSliceInt:        []int{1, 2, 3},
					SliceSeparatorOverride: []int{10, 20, 30},
					Custom:                 CustomType("no"),
					SliceCustom:            []CustomType{"a b", "b c", "c"},
					SliceCustomNoDefault:   []CustomType{},
					CustomStruct:           CustomStruct{First: "no", Last: "no"},
					Dive: Dive{
						SliceInt: []int{},
					},
				},
			},
		},
		"default override": {
			args: args{
				query: "defaultInt=20&defaultSliceInt=4&defaultSliceInt=5&defaultSliceInt=6",
			},
			want: wants{
				ts: testStruct{
					SliceInt:               []int{},
					DefaultInt:             20,
					DefaultSliceInt:        []int{4, 5, 6},
					SliceSeparatorOverride: []int{10, 20, 30},
					Custom:                 CustomType("no"),
					SliceCustom:            []CustomType{"a b", "b c", "c"},
					SliceCustomNoDefault:   []CustomType{},
					CustomStruct:           CustomStruct{First: "no", Last: "no"},
					Dive: Dive{
						SliceInt: []int{},
					},
				},
			},
		},
		"custom": {
			args: args{
				query: "custom=first+last&sliceCustom=a+d&sliceCustom=a+ee+r&sliceCustom=fff&sliceCustomNoDefault=a+b&sliceCustomNoDefault=c+d",
			},
			want: wants{
				ts: testStruct{
					SliceInt:               []int{},
					DefaultInt:             10,
					DefaultSliceInt:        []int{1, 2, 3},
					SliceSeparatorOverride: []int{10, 20, 30},
					Custom:                 CustomType("first last"),
					SliceCustom:            []CustomType{"a d", "a ee r", "fff"},
					SliceCustomNoDefault:   []CustomType{"a b", "c d"},
					CustomStruct:           CustomStruct{First: "no", Last: "no"},
					Dive: Dive{
						SliceInt: []int{},
					},
				},
			},
		},
		"dive": {
			args: args{
				query: "dive_int=10&dive_sliceInt=1&dive_sliceInt=2&dive_sliceInt=3",
			},
			want: wants{
				ts: testStruct{
					SliceInt:               []int{},
					DefaultInt:             10,
					DefaultSliceInt:        []int{1, 2, 3},
					SliceSeparatorOverride: []int{10, 20, 30},
					Custom:                 CustomType("no"),
					SliceCustom:            []CustomType{"a b", "b c", "c"},
					SliceCustomNoDefault:   []CustomType{},
					CustomStruct:           CustomStruct{First: "no", Last: "no"},
					Dive: Dive{
						Int:      10,
						SliceInt: []int{1, 2, 3},
					},
				},
			},
		},
	}
	parseFuncMap := map[reflect.Type]queryparam.ParserFunc{
		reflect.TypeOf(CustomType("")): parseCustomType,
		reflect.TypeOf(CustomStruct{}): parseCustomStruct,
		reflect.TypeOf(Dive{}):         parseDive,
	}
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			v, err := url.ParseQuery(tc.args.query)
			assert.NoError(t, err)
			p := queryparam.NewParser(v)
			var ts testStruct
			err = p.ParseWithOptions(&ts, queryparam.Options{
				DefaultValueSeparator: ",",
				TagName:               "queryParam",
				FuncMap:               parseFuncMap,
			})
			assert.NoError(t, err)
			assert.Equal(t, tc.want.ts, ts)
		})
	}
}

type testStruct struct {
	Bool    bool    `queryParam:"bool"`
	Int     int     `queryParam:"int"`
	Int8    int8    `queryParam:"int8"`
	Int16   int16   `queryParam:"int16"`
	Int32   int32   `queryParam:"int32"`
	Int64   int64   `queryParam:"int64"`
	Uint    uint    `queryParam:"uint"`
	Uint8   uint8   `queryParam:"uint8"`
	Uint16  uint16  `queryParam:"uint16"`
	Uint32  uint32  `queryParam:"uint32"`
	Uint64  uint64  `queryParam:"uint64"`
	Str     string  `queryParam:"str"`
	Float32 float32 `queryParam:"float32"`
	Float64 float64 `queryParam:"float64"`

	TimeDuration time.Duration `queryParam:"timeDuration"`

	DefaultInt int `queryParam:"defaultInt" paramDefault:"10"`

	SliceInt               []int `queryParam:"sliceInt"`
	DefaultSliceInt        []int `queryParam:"defaultSliceInt" paramDefault:"1,2,3"`
	SliceSeparatorOverride []int `queryParam:"sliceSeparatorOverride" paramDefault:"10|20|30" paramSeparator:"|"`

	Custom               CustomType   `queryParam:"custom"`
	SliceCustom          []CustomType `queryParam:"sliceCustom" paramSeparator:"@" paramDefault:"a+b@b+c@c"`
	SliceCustomNoDefault []CustomType `queryParam:"sliceCustomNoDefault"`

	CustomStruct CustomStruct `queryParam:"customStruct"`

	Dive Dive `queryParam:",dive"`
}

type CustomType string

func parseCustomType(v string) (any, error) {
	if v == "" {
		return CustomType("no"), nil
	}
	return CustomType(strings.ReplaceAll(v, "+", " ")), nil
}

type CustomStruct struct {
	First string
	Last  string
}

func parseCustomStruct(v string) (any, error) {
	if v == "" {
		return CustomStruct{
			First: "no",
			Last:  "no",
		}, nil
	}
	if !strings.Contains(v, " ") {
		return CustomStruct{
			First: v,
			Last:  "no",
		}, nil
	}
	s := strings.Split(v, " ")
	return CustomStruct{
		First: s[0],
		Last:  s[1],
	}, nil
}

type Dive struct {
	Int      int   `queryParam:"dive_int"`
	SliceInt []int `queryParam:"dive_sliceInt"`
}

func parseDive(_ string) (any, error) {
	return Dive{}, nil
}
