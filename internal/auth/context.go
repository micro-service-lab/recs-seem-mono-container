package auth

import (
	"context"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

type ctxKeyAuth int

const authKey ctxKeyAuth = 0

// NewContext 認証情報を含めた新しい context を返す。
func NewContext(ctx context.Context, member *entity.AuthMember) context.Context {
	return context.WithValue(ctx, authKey, member)
}

// FromContext context からセッション情報を取得する。
// 取得できない場合は nil を返す。
func FromContext(ctx context.Context) *entity.AuthMember {
	if sess, ok := ctx.Value(authKey).(*entity.AuthMember); ok {
		return sess
	}
	return nil
}
