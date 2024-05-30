package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/hasher"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/config"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
)

// ManageAuth 認証管理サービス。
type ManageAuth struct {
	DB             store.Store
	Hash           hasher.Hash
	Auth           auth.Auth
	SessionManager session.Manager
	Clocker        clock.Clock
	Config         config.Config
}

// Login ログインする。
func (m *ManageAuth) Login(
	ctx context.Context, loginID, password string,
) (entity.AuthJwt, error) {
	member, err := m.DB.FindMemberCredentialsByLoginID(ctx, loginID)
	if err != nil {
		return entity.AuthJwt{}, fmt.Errorf("failed to find member credentials by login ID: %w", err)
	}
	ok, err := m.Hash.Compare(password, member.Password)
	if err != nil || !ok {
		return entity.AuthJwt{}, errhandle.NewCommonError(response.InvalidLoginIDOrPassword, nil)
	}

	e, err := m.DB.FindMemberWithDetail(ctx, member.MemberID)
	if err != nil {
		return entity.AuthJwt{}, fmt.Errorf("failed to find member with detail: %w", err)
	}

	var mType session.MemberType

	if e.Student.Valid {
		mType = session.MemberTypeStudent
	} else if e.Professor.Valid {
		mType = session.MemberTypeProfessor
	} else {
		mType = session.MemberTypeInvalid
	}

	// access token
	sessID, token, err := m.Auth.NewSessionToken(mType, e.MemberID, m.Clocker.Now())
	if err != nil {
		return entity.AuthJwt{}, fmt.Errorf("failed to create session token: %w", err)
	}
	if err := m.SessionManager.UpdateSession(ctx, e.MemberID, sessID); err != nil {
		return entity.AuthJwt{}, fmt.Errorf("failed to update session: %w", err)
	}

	// refresh token
	refreshSessionID, refreshToken, err := m.Auth.NewRefreshToken(mType, e.MemberID, m.Clocker.Now())
	if err != nil {
		return entity.AuthJwt{}, fmt.Errorf("failed to create refresh token: %w", err)
	}
	if err := m.SessionManager.UpdateRefreshSession(ctx, e.MemberID, refreshSessionID); err != nil {
		return entity.AuthJwt{}, fmt.Errorf("failed to update refresh session: %w", err)
	}

	return entity.AuthJwt{
		AccessToken:  token,
		RefreshToken: refreshToken,
		SessionID:    sessID,
		ExpiresIn:    int(m.Config.AuthAccessTokenExpiresIn),
	}, nil
}

// RefreshToken トークンをリフレッシュする。
func (m *ManageAuth) RefreshToken(
	ctx context.Context, refreshToken string,
) (entity.AuthJwt, error) {
	token, err := m.Auth.ParseRefreshToken(refreshToken, m.Clocker.Now())
	if err != nil {
		return entity.AuthJwt{}, errhandle.NewCommonError(response.InvalidRefreshToken, nil)
	}

	ok, err := m.SessionManager.CheckRefreshSession(ctx, token.MemberID, token.ID)
	if err != nil {
		return entity.AuthJwt{}, fmt.Errorf("failed to check refresh session: %w", err)
	}
	if !ok {
		return entity.AuthJwt{}, errhandle.NewCommonError(response.InvalidRefreshToken, nil)
	}

	e, err := m.DB.FindMemberWithDetail(ctx, token.MemberID)
	if err != nil {
		return entity.AuthJwt{}, fmt.Errorf("failed to find member with detail: %w", err)
	}

	var mType session.MemberType

	if e.Student.Valid {
		mType = session.MemberTypeStudent
	} else if e.Professor.Valid {
		mType = session.MemberTypeProfessor
	} else {
		mType = session.MemberTypeInvalid
	}

	// access token
	sessID, newToken, err := m.Auth.NewSessionToken(mType, token.MemberID, m.Clocker.Now())
	if err != nil {
		return entity.AuthJwt{}, fmt.Errorf("failed to create session token: %w", err)
	}
	if err := m.SessionManager.UpdateSession(ctx, e.MemberID, sessID); err != nil {
		return entity.AuthJwt{}, fmt.Errorf("failed to update session: %w", err)
	}

	// refresh token
	refreshSessionID, refreshToken, err := m.Auth.NewRefreshToken(mType, e.MemberID, m.Clocker.Now())
	if err != nil {
		return entity.AuthJwt{}, fmt.Errorf("failed to create refresh token: %w", err)
	}
	if err := m.SessionManager.UpdateRefreshSession(ctx, e.MemberID, refreshSessionID); err != nil {
		return entity.AuthJwt{}, fmt.Errorf("failed to update refresh session: %w", err)
	}

	return entity.AuthJwt{
		AccessToken:  newToken,
		RefreshToken: refreshToken,
		SessionID:    sessID,
		ExpiresIn:    int(m.Config.AuthAccessTokenExpiresIn),
	}, nil
}

// Logout ログアウトする。
func (m *ManageAuth) Logout(ctx context.Context, memberID uuid.UUID) error {
	if err := m.SessionManager.DeleteSession(ctx, memberID); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
