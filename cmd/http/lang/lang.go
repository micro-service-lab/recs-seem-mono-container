// Package lang is defined middleware that sets the Accept-Language header.
package lang

import (
	"context"
	"net/http"
)

type localeKey struct{}

const (
	// LangEn is the English language.
	LangEn = "en"
	// LangJa is the Japanese language.
	LangJa = "ja"
)

// Langs is a list of supported languages.
var Langs = []string{LangEn, LangJa}

// DefaultLang is the default language.
const DefaultLang = LangEn

// Handler is a middleware that sets the Accept-Language header.
func Handler(lang string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lng := r.Header.Get("Accept-Language")
			if lng == "" {
				lng = lang
			}
			r = r.WithContext(setLocale(r.Context(), lng))

			next.ServeHTTP(w, r)
		})
	}
}

// setLocale sets the locale to the context.
func setLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, localeKey{}, locale)
}

// GetLocale gets the locale from the context.
func GetLocale(ctx context.Context) string {
	locale, ok := ctx.Value(localeKey{}).(string)
	if !ok {
		return DefaultLang
	}
	return locale
}