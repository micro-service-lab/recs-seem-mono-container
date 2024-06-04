package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
)

// AuthMiddleware API の認証・認可を行うミドルウェアを返す。
// base には API のベースパスを与える。
func AuthMiddleware(
	now func() time.Time,
	auth auth.Auth,
	service service.ManagerInterface,
	manager session.Manager,
) func(next http.Handler) http.Handler {
	mw := &authMiddleware{
		now:     now,
		auth:    auth,
		service: service,
		manager: manager,
	}
	return mw.handler
}

// authMiddleware API の認証・認可を行うミドルウェアを表す構造体。
type authMiddleware struct {
	now     func() time.Time
	auth    auth.Auth
	service service.ManagerInterface
	manager session.Manager
}

// handler API の認証・認可を行うミドルウェアの実装。
func (s *authMiddleware) handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ok, sess, err := s.authenticate(r)
		if err != nil {
			handled, err := errhandle.ErrorHandle(ctx, w, err)
			if !handled || err != nil {
				log.Printf("failed to handle error: %v", err)
			}
			return
		}
		if !ok {
			if err := response.JSONResponseWriter(ctx, w, response.Unauthorized, nil, nil); err != nil {
				log.Printf("failed to write response: %+v", err)
			}
			return
		}

		if sess != nil {
			m, err := s.service.FindAuthMemberByID(ctx, sess.MemberID)
			if err != nil {
				var me errhandle.ModelNotFoundError
				if errors.As(err, &me) {
					if err := response.JSONResponseWriter(ctx, w, response.Unauthorized, nil, nil); err != nil {
						log.Printf("failed to write response: %+v", err)
					}
					return
				}
				log.Printf("find member failed: %+v", err)
				handled, err := errhandle.ErrorHandle(ctx, w, err)
				if !handled || err != nil {
					log.Printf("failed to handle error: %v", err)
				}
				return
			}
			r = r.WithContext(session.NewContext(ctx, sess))
			r = r.WithContext(auth.NewContext(ctx, &m))
		}

		next.ServeHTTP(w, r)
	})
}

// authenticate メンバーの認証を行う。
func (s *authMiddleware) authenticate(r *http.Request) (bool, *session.Session, error) {
	token := extractToken(r)
	if token == "" {
		return false, nil, nil
	}

	sess, err := s.auth.ParseAccessToken(token, s.now())
	if err != nil {
		return false, nil, errhandle.NewCommonError(response.Unauthorized, nil)
	}

	ok, err := s.checkSession(r.Context(), sess)
	if err != nil {
		return false, nil, fmt.Errorf("check session: %w", err)
	}
	if !ok {
		return false, nil, nil
	}

	return true, sess, nil
}

// extractToken リクエストから Bearer トークンを抽出する。
func extractToken(r *http.Request) string {
	var token string
	// まず Authorization ヘッダーから Bearer トークンを取得する。
	v := strings.TrimSpace(r.Header.Get("Authorization"))
	token = strings.TrimSpace(strings.TrimPrefix(v, "Bearer"))

	// 次にクッキーからトークンを取得する。
	if token == "" {
		cookie, err := r.Cookie(auth.AccessTokenCookieKey)
		if err == nil {
			token = cookie.Value
		}
	}

	return token
}

// checkSession セッション情報とデータベースのメンバー情報とを照らし合わせ、
// 正常なセッションであることを確認する。
func (s *authMiddleware) checkSession(ctx context.Context, sess *session.Session) (bool, error) {
	switch sess.MemberType {
	case session.MemberTypeInvalid:
		// nop
	case session.MemberTypeStudent:
		return s.checkMember(ctx, sess)
	case session.MemberTypeProfessor:
		return s.checkMember(ctx, sess)
	}

	return false, nil
}

// checkMember セッション情報とデータベース情報とを照らし合わせ、
// 正常なセッションであることを確認する。
func (s *authMiddleware) checkMember(
	ctx context.Context, sess *session.Session,
) (bool, error) {
	_, err := s.service.FindMemberByID(ctx, sess.MemberID)
	if err != nil {
		var me errhandle.ModelNotFoundError
		if errors.As(err, &me) {
			return false, nil
		}
		return false, fmt.Errorf("find admin with detail: %w", err)
	}
	ok, err := s.manager.CheckSession(ctx, sess.MemberID, sess.ID)
	if err != nil {
		return false, fmt.Errorf("check session: %w", err)
	}

	return ok, nil
}
