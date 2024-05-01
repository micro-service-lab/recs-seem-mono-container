// Package queryparam provides a query parameter parser.
package queryparam

import (
	"encoding"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/stoewer/go-strcase"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/agerror"
)

var defaultBuiltInParsers = map[reflect.Kind]ParserFunc{
	reflect.Bool: func(v string) (any, error) {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false, nil
		}
		return b, nil
	},
	reflect.String: func(v string) (any, error) {
		return v, nil
	},
	reflect.Int: func(v string) (any, error) {
		i, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return 0, nil
		}
		return int(i), nil
	},
	reflect.Int16: func(v string) (any, error) {
		i, err := strconv.ParseInt(v, 10, 16)
		if err != nil {
			return 0, nil
		}
		return int16(i), nil
	},
	reflect.Int32: func(v string) (any, error) {
		i, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return 0, nil
		}
		return int32(i), nil
	},
	reflect.Int64: func(v string) (any, error) {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, nil
		}
		return i, nil
	},
	reflect.Int8: func(v string) (any, error) {
		i, err := strconv.ParseInt(v, 10, 8)
		if err != nil {
			return 0, nil
		}
		return int8(i), nil
	},
	reflect.Uint: func(v string) (any, error) {
		i, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0, nil
		}
		return uint(i), nil
	},
	reflect.Uint16: func(v string) (any, error) {
		i, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return 0, nil
		}
		return uint16(i), nil
	},
	reflect.Uint32: func(v string) (any, error) {
		i, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0, nil
		}
		return uint32(i), nil
	},
	reflect.Uint64: func(v string) (any, error) {
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, nil
		}
		return i, nil
	},
	reflect.Uint8: func(v string) (any, error) {
		i, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return 0, nil
		}
		return uint8(i), nil
	},
	reflect.Float64: func(v string) (any, error) {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, nil
		}
		return f, nil
	},
	reflect.Float32: func(v string) (any, error) {
		f, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return 0, nil
		}
		return float32(f), nil
	},
}

func defaultTypeParsers() map[reflect.Type]ParserFunc {
	return map[reflect.Type]ParserFunc{
		reflect.TypeOf(time.Nanosecond): func(v string) (any, error) {
			s, err := time.ParseDuration(v)
			if err != nil {
				return time.Duration(0), nil
			}
			return s, nil
		},
	}
}

// ParserFunc defines the signature of a function that can be used within `CustomParsers`.
type ParserFunc func(v string) (any, error)

// OnSetFn is a hook that can be run when a value is set.
type OnSetFn func(tag string, value any, isDefault bool)

// processFieldFn is a function which takes all information about a field and processes it.
type processFieldFn func(
	refField reflect.Value, refTypeField reflect.StructField, opts Options, fieldParams FieldParams) error

// Options for the parser.
type Options struct {
	// Param keys and values that will be accessible for the service.
	Param map[string][]string

	// DefaultValueSeparator for slice values.
	DefaultValueSeparator string

	// TagName specifies another tagname to use rather than the default queryParam.
	TagName string

	// RequiredIfNoDef automatically sets all param as required if they do not
	// declare 'paramDefault'.
	RequiredIfNoDef bool

	// OnSet allows to run a function when a value is set.
	OnSet OnSetFn

	// Prefix define a prefix for each key.
	Prefix string

	// UseFieldNameByDefault defines whether or not param should use the field
	// name by default if the `queryParam` key is missing.
	// Note that the field name will be "converted" to conform with param
	// names conventions.
	UseFieldNameByDefault bool

	// Custom parse functions for different types.
	FuncMap map[reflect.Type]ParserFunc

	// Used internally. maps the param key to its resolved string value.
	rawParams map[string][]string
}

func defaultOptions(param url.Values) Options {
	return Options{
		TagName:               "queryParam",
		DefaultValueSeparator: "&",
		Param:                 param,
		FuncMap:               defaultTypeParsers(),
		rawParams:             make(map[string][]string),
	}
}

func customOptions(param url.Values, opt Options) Options {
	defOpts := defaultOptions(param)
	if opt.TagName == "" {
		opt.TagName = defOpts.TagName
	}
	if opt.Param == nil {
		opt.Param = defOpts.Param
	}
	if opt.FuncMap == nil {
		opt.FuncMap = map[reflect.Type]ParserFunc{}
	}
	if opt.DefaultValueSeparator == "" {
		opt.DefaultValueSeparator = defOpts.DefaultValueSeparator
	}
	if opt.rawParams == nil {
		opt.rawParams = defOpts.rawParams
	}
	for k, v := range defOpts.FuncMap {
		if _, exists := opt.FuncMap[k]; !exists {
			opt.FuncMap[k] = v
		}
	}
	return opt
}

func optionsWithParamPrefix(field reflect.StructField, opts Options) Options {
	return Options{
		Param:                 opts.Param,
		TagName:               opts.TagName,
		RequiredIfNoDef:       opts.RequiredIfNoDef,
		OnSet:                 opts.OnSet,
		Prefix:                opts.Prefix + field.Tag.Get("paramPrefix"),
		UseFieldNameByDefault: opts.UseFieldNameByDefault,
		FuncMap:               opts.FuncMap,
		DefaultValueSeparator: opts.DefaultValueSeparator,
		rawParams:             opts.rawParams,
	}
}

// Parser is a struct that can be used to parse query params.
type Parser struct {
	Param url.Values
}

// NewParser creates a new parser with the given query params.
func NewParser(param url.Values) *Parser {
	return &Parser{Param: param}
}

// Parse parses a struct containing `queryParam` tags and loads its values from
// query params.
func (p *Parser) Parse(v any) error {
	return parseInternal(v, setField, defaultOptions(p.Param))
}

// ParseWithOptions parses a struct containing `queryParam` tags and loads its values from
// query params.
func (p *Parser) ParseWithOptions(v any, opts Options) error {
	return parseInternal(v, setField, customOptions(p.Param, opts))
}

// GetFieldParams parses a struct containing `queryParam` tags and returns information about
// tags it found.
func (p *Parser) GetFieldParams(v any) ([]FieldParams, error) {
	return p.GetFieldParamsWithOptions(v, defaultOptions(p.Param))
}

// GetFieldParamsWithOptions parses a struct containing `queryParam` tags and returns information about
// tags it found.
func (p *Parser) GetFieldParamsWithOptions(v any, opts Options) ([]FieldParams, error) {
	// フィールドのメタデータ
	var result []FieldParams
	err := parseInternal(
		v,
		func(_ reflect.Value, _ reflect.StructField, _ Options, fieldParams FieldParams) error {
			if fieldParams.OwnKey != "" { // キーが空でない場合
				result = append(result, fieldParams)
			}
			return nil
		},
		customOptions(p.Param, opts),
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func parseInternal(v any, processField processFieldFn, opts Options) error {
	// 値がポインタ型であることを確認する
	ptrRef := reflect.ValueOf(v)
	if ptrRef.Kind() != reflect.Ptr {
		return agerror.NewAggregateErrorWithErr(NotStructPtrError{}, "query.Parse")
	}
	// ポインタの指す先が構造体であることを確認する
	ref := ptrRef.Elem()
	if ref.Kind() != reflect.Struct {
		return agerror.NewAggregateErrorWithErr(NotStructPtrError{}, "query.Parse")
	}

	return doParse(ref, processField, opts)
}

func doParse(ref reflect.Value, processField processFieldFn, opts Options) error {
	refType := ref.Type()

	agrErr := agerror.NewAggregateError("query.Parse")

	// 構造体のフィールドを反復処理する
	for i := 0; i < refType.NumField(); i++ {
		// 構造体のフィールドを取得する
		refField := ref.Field(i)
		// 構造体のフィールドのメタデータを取得する
		refTypeField := refType.Field(i)

		if err := doParseField(refField, refTypeField, processField, opts); err != nil {
			//nolint: errorlint
			if val, ok := err.(agerror.AggregateError); ok {
				agrErr.Errors = append(agrErr.Errors, val.Errors...)
			} else {
				agrErr.Errors = append(agrErr.Errors, err)
			}
		}
	}

	if len(agrErr.Errors) == 0 {
		return nil
	}

	return agrErr
}

// フィールドごとの処理を行う
func doParseField(
	refField reflect.Value, refTypeField reflect.StructField, processField processFieldFn, opts Options,
) error {
	if !refField.CanSet() { // フィールドが設定可能でない場合
		return nil
	}
	if reflect.Ptr == refField.Kind() && !refField.IsNil() { // フィールドがポインタ型であり、nil でない場合
		// ポインタの指す先の処理を再帰的に行う
		return parseInternal(refField.Interface(), processField, optionsWithParamPrefix(refTypeField, opts))
	}
	if reflect.Struct == refField.Kind() &&
		refField.CanAddr() && refField.Type().Name() == "" { // フィールドが構造体であり、アドレスを取得可能であり、名前が空の場合
		// 無名構造体の中の処理を再帰的に行う
		err := parseInternal(
			refField.Addr().Interface(), processField, optionsWithParamPrefix(refTypeField, opts))
		if err != nil {
			return fmt.Errorf("failed to parse field %s: %w", refTypeField.Name, err)
		}
	}

	params, err := parseFieldParams(refTypeField, opts) // フィールドのメタデータを解析する
	if err != nil {
		return err
	}

	// processFieldFnを実行する
	// Parse系ならrefFieldに値を設定する
	// GetFieldParams系ならFieldParamsを内部に保持する
	if err := processField(refField, refTypeField, opts, params); err != nil {
		return err
	}

	if reflect.Struct == refField.Kind() { // フィールドが構造体である場合
		// 構造体の中の処理を再帰的に行う(定義済み構造体)
		return doParse(refField, processField, optionsWithParamPrefix(refTypeField, opts))
	}

	return nil
}

// Parseする際のprocessFieldFnの実装
func setField(refField reflect.Value, refTypeField reflect.StructField, opts Options, fieldParams FieldParams) error {
	// default値なども考慮しながら値を取得する
	value := get(fieldParams, opts)

	if len(value) == 0 {
		value = []string{""}
	}

	return set(refField, refTypeField, value, opts.FuncMap)
}

// フィールド名からキー(タグの名前)を生成する
func toParamName(input string) string {
	return strcase.SnakeCase(input)
}

// FieldParams contains information about parsed field tags.
type FieldParams struct {
	OwnKey          string
	Key             string
	DefaultValue    []string
	HasDefaultValue bool
}

// フィールドのメタデータを解析する
func parseFieldParams(field reflect.StructField, opts Options) (FieldParams, error) {
	// opts.TagNameの値が(ownKey,tags(,区切りの配列))に分割される
	ownKey, tags := parseKeyForOption(field.Tag.Get(opts.TagName))
	if ownKey == "" && opts.UseFieldNameByDefault { // queryParamタグがなく(もしくは空)、かつUseFieldNameByDefaultがtrueの場合
		ownKey = toParamName(field.Name) // フィールド名からキーを生成する
	}

	// デフォルト設定の取得
	defaultValue, hasDefaultValue := field.Tag.Lookup("paramDefault")

	defaultValueSeparator, hasDefaultValueSeparator := field.Tag.Lookup("paramDefaultSeparator")
	if !hasDefaultValueSeparator {
		defaultValueSeparator = opts.DefaultValueSeparator
	}
	defaultValues := strings.Split(defaultValue, defaultValueSeparator)

	result := FieldParams{
		OwnKey:          ownKey,
		Key:             opts.Prefix + ownKey,
		DefaultValue:    defaultValues,
		HasDefaultValue: hasDefaultValue,
	}

	for _, tag := range tags {
		switch tag {
		case "":
			continue
		default:
			return FieldParams{}, newNoSupportedTagOptionError(tag)
		}
	}

	return result, nil
}

func get(fieldParams FieldParams, opts Options) (val []string) {
	var isDefault bool

	val, isDefault = getOr(fieldParams.Key, fieldParams.DefaultValue, fieldParams.HasDefaultValue, opts.Param)

	opts.rawParams[fieldParams.OwnKey] = val

	if opts.OnSet != nil { // 値が設定されたときに実行される関数がある場合
		if fieldParams.OwnKey != "" { // queryParamタグがある場合
			opts.OnSet(fieldParams.Key, val, isDefault) // 設定されたフックを実行
		}
	}
	return val
}

// split the env tag's key into the expected key and desired option, if any.
func parseKeyForOption(key string) (string, []string) {
	opts := strings.Split(key, ",")
	return opts[0], opts[1:]
}

func getOr(
	key string, defaultValue []string, defExists bool, param map[string][]string,
) (val []string, isDefault bool) {
	// 値を取得する
	value, exists := param[key]
	switch {
	case (!exists || key == "") && defExists: // paramにはなく、default値がある場合
		return defaultValue, true
	case exists && len(value) == 0 && defExists: // paramにはあるが、値がなく、default値がある場合
		return defaultValue, true
	case !exists: // paramになく、default値もない場合
		return []string{}, false
	}

	return value, true
}

func set(field reflect.Value, sf reflect.StructField, value []string, funcMap map[reflect.Type]ParserFunc) error {
	if tm := asTextUnmarshaler(field); tm != nil { // フィールド値がencoding.TextUnmarshalerを実装している場合
		if err := tm.UnmarshalText([]byte(value[0])); err != nil { // フィールド値に対してUnmarshalTextを実行する
			return newParseError(sf, err)
		}
		return nil
	}

	typee := sf.Type                 // フィールドの型
	fieldee := field                 // フィールドの値
	if typee.Kind() == reflect.Ptr { // フィールドの型がポインタ型である場合
		typee = typee.Elem()   // ポインタが示す型を取得する
		fieldee = field.Elem() // ポインタが示す値を取得する
	}

	parserFunc, ok := funcMap[typee] // フィールドの型に対応するカスタムパーサー関数を取得する
	if ok {
		val, err := parserFunc(value[0]) // 型ごとに設定されたパーサー関数を実行する
		if err != nil {
			return newParseError(sf, err)
		}

		fieldee.Set(reflect.ValueOf(val)) // フィールドに値を設定する
		return nil
	}

	parserFunc, ok = defaultBuiltInParsers[typee.Kind()] // フィールドの型に対応するビルトインのパーサー関数を取得する
	if ok {
		val, err := parserFunc(value[0])
		if err != nil {
			return newParseError(sf, err)
		}

		fieldee.Set(reflect.ValueOf(val).Convert(typee))
		return nil
	}

	//nolint: exhaustive
	switch field.Kind() {
	case reflect.Slice: // フィールドがスライス型である場合
		return handleSlice(field, value, sf, funcMap)
	}

	return newNoParserError(sf)
}

func handleSlice(
	field reflect.Value, value []string, sf reflect.StructField, funcMap map[reflect.Type]ParserFunc,
) error {
	typee := sf.Type.Elem()          // 配列の要素の型を取得する
	if typee.Kind() == reflect.Ptr { // 配列の要素の型がポインタ型である場合
		typee = typee.Elem() // ポインタが示す型を取得する
	}

	if _, ok := reflect.New(typee).Interface().(encoding.TextUnmarshaler); ok { // 配列の要素がencoding.TextUnmarshalerを実装している場合
		return parseTextUnmarshalers(field, value, sf) // 配列の要素ずつのTextUnmarshalerであるデータを解析する
	}

	// parseFuncを取得する
	parserFunc, ok := funcMap[typee]
	if !ok {
		parserFunc, ok = defaultBuiltInParsers[typee.Kind()]
		if !ok { // パーサー関数がない場合
			return newNoParserError(sf)
		}
	}

	result := reflect.MakeSlice(sf.Type, 0, len(value)) // 要素の型とデータ数を指定してスライスを作成する
	for _, v := range value {
		r, err := parserFunc(v) // パーサー関数を実行する
		if err != nil {
			return newParseError(sf, err)
		}
		v := reflect.ValueOf(r).Convert(typee)    // パーサー関数の結果を指定された型に変換する
		if sf.Type.Elem().Kind() == reflect.Ptr { // 配列の要素がポインタ型である場合
			v = reflect.New(typee)                          // ポインタが示す型のゼロ値を取得する
			v.Elem().Set(reflect.ValueOf(r).Convert(typee)) // ポインタが示す型に変換した値を設定する
		}
		result = reflect.Append(result, v) // スライスに値を追加する
	}
	field.Set(result) // フィールドにスライスを設定する
	return nil
}

func asTextUnmarshaler(field reflect.Value) encoding.TextUnmarshaler {
	if reflect.Ptr == field.Kind() { // フィールドがポインタ型である場合
		if field.IsNil() { // フィールドがnilである場合
			field.Set(reflect.New(field.Type().Elem())) // フィールドに新しい値を設定(設定値はポインタが示す型のゼロ値)
		}
	} else if field.CanAddr() { // フィールドがアドレスを取得可能である場合
		field = field.Addr() // フィールドのアドレスを取得する
	}

	tm, ok := field.Interface().(encoding.TextUnmarshaler) // フィールドがencoding.TextUnmarshalerを実装しているか確認する
	if !ok {                                               // 実装していない場合
		return nil
	}
	return tm
}

func parseTextUnmarshalers(field reflect.Value, data []string, sf reflect.StructField) error {
	s := len(data)                                              // データ数
	elemType := field.Type().Elem()                             // 配列の要素の型
	slice := reflect.MakeSlice(reflect.SliceOf(elemType), s, s) // データ数と要素の型を指定してスライスを作成する
	for i, v := range data {
		sv := slice.Index(i)
		kind := sv.Kind()
		if kind == reflect.Ptr { // ポインタ型である場合
			sv = reflect.New(elemType.Elem()) // ポインタが示す型のゼロ値を取得する
		} else { // ポインタ型でない場合
			sv = sv.Addr() // アドレスを取得する
		}
		tm, ok := sv.Interface().(encoding.TextUnmarshaler)
		if !ok { // TextUnmarshalerでない場合
			return newParseError(sf, fmt.Errorf("field %s is not a TextUnmarshaler", sf.Name))
		}
		if err := tm.UnmarshalText([]byte(v)); err != nil { // TextUnmarshalerであるデータを解析する
			return newParseError(sf, err)
		}
		if kind == reflect.Ptr { // ポインタ型である場合
			slice.Index(i).Set(sv) // スライスに値を設定する
		}
	}

	field.Set(slice) // フィールドにスライスを設定する

	return nil
}
