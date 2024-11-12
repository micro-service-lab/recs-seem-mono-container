package session

import (
	"context"
)

type ctxKeySession int

const sessionKey ctxKeySession = 0

// NewContext セッション情報を含めた新しい context を返す。
func NewContext(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, sessionKey, sess)
}

// FromContext context からセッション情報を取得する。
// 取得できない場合は nil を返す。
func FromContext(ctx context.Context) *Session {
	if sess, ok := ctx.Value(sessionKey).(*Session); ok {
		return sess
	}
	return nil
}
