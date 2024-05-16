// Package validation provides a request validation.
package validation

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	ja_translations "github.com/go-playground/validator/v10/translations/ja"
	"github.com/iancoleman/strcase"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
)

// DefaultTranslator is a default translator.
type DefaultTranslator func(v *validator.Validate, trans ut.Translator) (err error)

// Validator validates a request.
type Validator interface {
	Validate(ctx context.Context, req any) error
	ValidateWithLocale(ctx context.Context, req any, locale string) error
	GetDefaultLocale() string
	SetDefaultLocale(locale string)
}

// RequestValidator is a request validator.
type RequestValidator struct {
	mu       sync.RWMutex
	locale   string
	transMap map[string]*ut.Translator
	vMap     map[string]*validator.Validate
}

// NewRequestValidator creates a new request validator.
func NewRequestValidator() (*RequestValidator, error) {
	var rv RequestValidator
	rv.transMap = make(map[string]*ut.Translator, len(lang.Langs))
	rv.vMap = make(map[string]*validator.Validate, len(lang.Langs))
	en := en.New()
	uni := ut.New(en, en, ja.New())

	for _, lang := range lang.Langs {
		t, _ := uni.GetTranslator(lang)
		v := validator.New()
		err := findDefaultTranslator(lang)(v, t)
		if err != nil {
			return nil, err
		}

		rv.vMap[lang] = v
		rv.transMap[lang] = &t
	}

	for lang, v := range rv.vMap {
		func(vv *validator.Validate, l string) {
			vv.RegisterTagNameFunc(func(fld reflect.StructField) string {
				fieldName := fld.Tag.Get(l)
				if fieldName == "-" {
					return ""
				}
				return fieldName
			})
		}(v, lang)
	}

	return &rv, nil
}

// GetDefaultLocale returns the default locale.
func (v *RequestValidator) GetDefaultLocale() string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.locale
}

// SetDefaultLocale sets the default locale.
func (v *RequestValidator) SetDefaultLocale(locale string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.locale = locale
}

// Validate validates a request.
func (v *RequestValidator) Validate(ctx context.Context, req any) error {
	return v.ValidateWithLocale(ctx, req, v.GetDefaultLocale())
}

// ValidateWithLocale validates a request with a locale.
func (v *RequestValidator) ValidateWithLocale(ctx context.Context, req any, locale string) error {
	v.mu.RLock()
	defer v.mu.RUnlock()
	ev, ok := v.vMap[locale]
	if !ok {
		ev = v.vMap[lang.DefaultLang]
	}
	et, ok := v.transMap[locale]
	if !ok {
		et = v.transMap[lang.DefaultLang]
	}
	err := ev.StructCtx(ctx, req)

	errAttr := make(map[string][]string)

	if err != nil {
		var errs validator.ValidationErrors

		// translate all error at once
		if ok := errors.As(err, &errs); ok {
			for _, err := range errs {
				fieldName := strcase.ToSnake(err.StructField())
				errAttr[fieldName] = append(errAttr[fieldName], err.Translate(*et))
			}
		}

		return errhandle.NewValidationError(errAttr)
	}

	return nil
}

func findDefaultTranslator(locale string) DefaultTranslator {
	switch locale {
	case lang.LangJa:
		return ja_translations.RegisterDefaultTranslations
	case lang.LangEn:
		return en_translations.RegisterDefaultTranslations
	default:
		return en_translations.RegisterDefaultTranslations
	}
}

var _ Validator = (*RequestValidator)(nil)
