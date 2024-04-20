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

// AdminID context のセッション情報から管理者アカウント ID を取得する。
// 取得できない場合は 0 を返す。
func AdminID(ctx context.Context) string {
	sess := FromContext(ctx)
	if sess == nil {
		return ""
	}

	return sess.AdminID()
}

// UserID context のセッション情報からユーザー ID を取得する。
// 取得できない場合は 0 を返す。
func UserID(ctx context.Context) int32 {
	sess := FromContext(ctx)
	if sess == nil {
		return 0
	}

	return sess.UserID()
}
