package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
)

// ManageMembership オーガナイゼーション所属管理サービス。
type ManageMembership struct {
	DB      store.Store
	Clocker clock.Clock
}

// BelongMemberOnOrganization メンバーを組織に所属させる。
func (m *ManageMembership) BelongMemberOnOrganization(
	ctx context.Context,
	organizationID,
	ownerID uuid.UUID,
	memberIDs []uuid.UUID,
) (e int64, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	now := m.Clocker.Now()
	owner, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(OrganizationBelongingTargetOwner)
		}
		return 0, fmt.Errorf("failed to find member: %w", err)
	}
	org, err := m.DB.FindOrganizationByIDWithSd(ctx, sd, organizationID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(OrganizationBelongingTargetOrganization)
		}
		return 0, fmt.Errorf("failed to find organization: %w", err)
	}
	mm, err := m.DB.GetPluralMembersWithSd(
		ctx,
		sd,
		memberIDs,
		parameter.MemberOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get plural members: %w", err)
	}
	if len(mm.Data) != len(memberIDs) {
		return 0, errhandle.NewModelNotFoundError(OrganizationBelongingTargetMembers)
	}
	bcrp := make([]parameter.BelongOrganizationParam, 0, len(mm.Data))
	for _, m := range mm.Data {
		bcrp = append(bcrp, parameter.BelongOrganizationParam{
			OrganizationID: org.OrganizationID,
			MemberID:       m.MemberID,
			AddedAt:        now,
		})
	}
	if e, err = m.DB.BelongOrganizationsWithSd(
		ctx,
		sd,
		bcrp,
	); err != nil {
		var dpe errhandle.ModelDuplicatedError
		if errors.As(err, &dpe) {
			return 0, errhandle.NewModelDuplicatedError(OrganizationBelongingTargetOrganizationBelongings)
		}
		return 0, fmt.Errorf("failed to belong chat rooms: %w", err)
	}

	if org.ChatRoomID.Valid {
		room, err := m.DB.FindChatRoomByIDWithSd(
			ctx,
			sd,
			org.ChatRoomID.Bytes,
		)
		if err != nil {
			var nfe errhandle.ModelNotFoundError
			if errors.As(err, &nfe) {
				return 0, errhandle.NewModelNotFoundError(ChatRoomBelongingTargetChatRoom)
			}
			return 0, fmt.Errorf("failed to find chat room: %w", err)
		}
		if _, err = belongMembersOnChatRoom(
			ctx,
			sd,
			now,
			m.DB,
			room,
			owner,
			mm.Data,
			true,
		); err != nil {
			return 0, fmt.Errorf("failed to belong members on chat room: %w", err)
		}
	}

	return e, nil
}

// RemoveMembersFromOrganization 組織からメンバーを削除する。
func (m *ManageMembership) RemoveMembersFromOrganization(
	ctx context.Context,
	organizationID,
	ownerID uuid.UUID,
	memberIDs []uuid.UUID,
) (e int64, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	now := m.Clocker.Now()
	owner, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(OrganizationBelongingTargetOwner)
		}
		return 0, fmt.Errorf("failed to find member: %w", err)
	}
	org, err := m.DB.FindOrganizationByIDWithSd(ctx, sd, organizationID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(OrganizationBelongingTargetOrganization)
		}
		return 0, fmt.Errorf("failed to find organization: %w", err)
	}
	mm, err := m.DB.GetPluralMembersWithSd(
		ctx,
		sd,
		memberIDs,
		parameter.MemberOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get plural members: %w", err)
	}
	if len(mm.Data) != len(memberIDs) {
		return 0, errhandle.NewModelNotFoundError(OrganizationBelongingTargetMembers)
	}
	if e, err = m.DB.DisbelongPluralMembersOnOrganizationWithSd(
		ctx,
		sd,
		organizationID,
		memberIDs,
	); err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(OrganizationBelongingTargetOrganizationBelongings)
		}
		return 0, fmt.Errorf("failed to disbelong members on organization: %w", err)
	}

	if org.ChatRoomID.Valid {
		room, err := m.DB.FindChatRoomByIDWithSd(
			ctx,
			sd,
			org.ChatRoomID.Bytes,
		)
		if err != nil {
			var nfe errhandle.ModelNotFoundError
			if errors.As(err, &nfe) {
				return 0, errhandle.NewModelNotFoundError(ChatRoomBelongingTargetChatRoom)
			}
			return 0, fmt.Errorf("failed to find chat room: %w", err)
		}
		if _, err = removeMembersFromChatRoom(
			ctx,
			sd,
			now,
			m.DB,
			room,
			owner,
			mm.Data,
			true,
		); err != nil {
			return 0, fmt.Errorf("failed to remove members from chat room: %w", err)
		}
	}

	return e, nil
}

// WithdrawMemberFromOrganization 組織から退会する。
func (m *ManageMembership) WithdrawMemberFromOrganization(
	ctx context.Context,
	organizationID,
	memberID uuid.UUID,
) (e int64, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	now := m.Clocker.Now()
	member, err := m.DB.FindMemberByIDWithSd(ctx, sd, memberID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(OrganizationBelongingTargetMembers)
		}
		return 0, fmt.Errorf("failed to find member: %w", err)
	}
	org, err := m.DB.FindOrganizationByIDWithSd(ctx, sd, organizationID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(OrganizationBelongingTargetOrganization)
		}
		return 0, fmt.Errorf("failed to find organization: %w", err)
	}
	if e, err = m.DB.DisbelongOrganizationWithSd(
		ctx,
		sd,
		memberID,
		organizationID,
	); err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(OrganizationBelongingTargetOrganizationBelongings)
		}
		return 0, fmt.Errorf("failed to disbelong member on organization: %w", err)
	}

	if org.ChatRoomID.Valid {
		room, err := m.DB.FindChatRoomByIDWithSd(
			ctx,
			sd,
			org.ChatRoomID.Bytes,
		)
		if err != nil {
			var nfe errhandle.ModelNotFoundError
			if errors.As(err, &nfe) {
				return 0, errhandle.NewModelNotFoundError(ChatRoomBelongingTargetChatRoom)
			}
			return 0, fmt.Errorf("failed to find chat room: %w", err)
		}
		if _, err = withdrawMemberFromChatRoom(
			ctx,
			sd,
			now,
			m.DB,
			room,
			member,
			true,
		); err != nil {
			return 0, fmt.Errorf("failed to withdraw member from chat room: %w", err)
		}
	}

	return e, nil
}

// GetOrganizationsOnMember メンバーに関連付けられたチャットルームを取得する。
func (m *ManageMembership) GetOrganizationsOnMember(
	ctx context.Context,
	memberID uuid.UUID,
	whereSearchName string,
	order parameter.OrganizationOnMemberOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (es store.ListResult[entity.OrganizationOnMember], err error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereOrganizationOnMemberParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset), Valid: true},
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.NonePagination:
	}
	es, err = m.DB.GetOrganizationsOnMember(ctx, memberID, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.OrganizationOnMember]{}, fmt.Errorf("failed to get organizations on member: %w", err)
	}
	return es, nil
}

// GetOrganizationsOnMemberCount メンバーに関連付けられたチャットルームの数を取得する。
func (m *ManageMembership) GetOrganizationsOnMemberCount(
	ctx context.Context, memberID uuid.UUID,
	whereSearchName string,
) (es int64, err error) {
	where := parameter.WhereOrganizationOnMemberParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	es, err = m.DB.CountOrganizationsOnMember(ctx, memberID, where)
	if err != nil {
		return 0, fmt.Errorf("failed to get organizations on member count: %w", err)
	}
	return es, nil
}

// GetMembersOnOrganization チャットルームに関連付けられたメンバーを取得する。
func (m *ManageMembership) GetMembersOnOrganization(
	ctx context.Context,
	chatRoomID uuid.UUID,
	whereSearchName string,
	order parameter.MemberOnOrganizationOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (es store.ListResult[entity.MemberOnOrganization], err error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereMemberOnOrganizationParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset), Valid: true},
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.NonePagination:
	}
	es, err = m.DB.GetMembersOnOrganization(ctx, chatRoomID, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.MemberOnOrganization]{}, fmt.Errorf("failed to get members on organization: %w", err)
	}
	return es, nil
}

// GetMembersOnOrganizationCount チャットルームに関連付けられたメンバーの数を取得する。
func (m *ManageMembership) GetMembersOnOrganizationCount(
	ctx context.Context, chatRoomID uuid.UUID,
	whereSearchName string,
) (es int64, err error) {
	where := parameter.WhereMemberOnOrganizationParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	es, err = m.DB.CountMembersOnOrganization(ctx, chatRoomID, where)
	if err != nil {
		return 0, fmt.Errorf("failed to get members on organization count: %w", err)
	}
	return es, nil
}
