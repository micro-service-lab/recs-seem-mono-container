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
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/ws"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
)

// ManageMembership オーガナイゼーション所属管理サービス。
type ManageMembership struct {
	DB      store.Store
	Clocker clock.Clock
	WsHub   ws.HubInterface
}

// BelongMembersOnOrganization メンバーを組織に所属させる。
func (m *ManageMembership) BelongMembersOnOrganization(
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
	org, err := m.DB.FindOrganizationWithDetailWithSd(ctx, sd, organizationID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(OrganizationBelongingTargetOrganization)
		}
		return 0, fmt.Errorf("failed to find organization: %w", err)
	}
	if org.Grade.Valid {
		return 0, errhandle.NewCommonError(response.AttemptOperateGradeOrganization, nil)
	}
	if org.Group.Valid {
		return 0, errhandle.NewCommonError(response.AttemptOperateGroupOrganization, nil)
	}
	if org.IsPersonal {
		return 0, errhandle.NewCommonError(response.AttemptOperatePersonalOrganization, nil)
	}
	if org.IsWhole {
		return 0, errhandle.NewCommonError(response.AttemptOperateWholeOrganization, nil)
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

	belongingMembers, err := m.DB.GetMembersOnOrganizationWithSd(
		ctx,
		sd,
		organizationID,
		parameter.WhereMemberOnOrganizationParam{},
		parameter.MemberOnOrganizationOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get members on organization: %w", err)
	}
	alreadyMemberIDs := make([]uuid.UUID, 0, len(belongingMembers.Data))
	for _, v := range belongingMembers.Data {
		alreadyMemberIDs = append(alreadyMemberIDs, v.Member.MemberID)
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
		var action entity.ChatRoomAddMemberActionWithAddedByAndAddMembers
		var actAttr entity.ChatRoomAction
		if _, action, actAttr, err = belongMembersOnChatRoom(
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
		defer func(
			room entity.ChatRoom, membersIDs, alreadyMemberIDs []uuid.UUID,
			action entity.ChatRoomAddMemberActionWithAddedByAndAddMembers,
			actAttr entity.ChatRoomAction,
		) {
			if err == nil {
				m.WsHub.Dispatch(ws.EventTypeChatRoomAddedMe, ws.Targets{
					Members: membersIDs,
				}, ws.ChatRoomAddedMeEventData{
					ChatRoom: room,
				})
				m.WsHub.Dispatch(ws.EventTypeChatRoomAddedMember, ws.Targets{
					Members: alreadyMemberIDs,
				}, ws.ChatRoomAddedMemberEventData{
					ChatRoomID:           room.ChatRoomID,
					Action:               action,
					ChatRoomActionID:     actAttr.ChatRoomActionID,
					ChatRoomActionTypeID: actAttr.ChatRoomActionTypeID,
					ActedAt:              actAttr.ActedAt,
				})
			}
		}(room, memberIDs, alreadyMemberIDs, action, actAttr)
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
	var inMembers bool
	for _, v := range memberIDs {
		if v == ownerID {
			inMembers = true
			break
		}
	}
	if inMembers {
		return 0, errhandle.NewCommonError(response.CannotDeleteSelfFromOrganization, nil)
	}
	owner, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(OrganizationBelongingTargetOwner)
		}
		return 0, fmt.Errorf("failed to find member: %w", err)
	}
	org, err := m.DB.FindOrganizationWithDetailWithSd(ctx, sd, organizationID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(OrganizationBelongingTargetOrganization)
		}
		return 0, fmt.Errorf("failed to find organization: %w", err)
	}
	if org.Grade.Valid {
		return 0, errhandle.NewCommonError(response.AttemptOperateGradeOrganization, nil)
	}
	if org.Group.Valid {
		return 0, errhandle.NewCommonError(response.AttemptOperateGroupOrganization, nil)
	}
	if org.IsPersonal {
		return 0, errhandle.NewCommonError(response.AttemptOperatePersonalOrganization, nil)
	}
	if org.IsWhole {
		return 0, errhandle.NewCommonError(response.AttemptOperateWholeOrganization, nil)
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
	belongingMembers, err := m.DB.GetMembersOnOrganizationWithSd(
		ctx,
		sd,
		organizationID,
		parameter.WhereMemberOnOrganizationParam{},
		parameter.MemberOnOrganizationOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get members on organization: %w", err)
	}
	removeMembers := make(map[uuid.UUID]struct{}, len(mm.Data))
	for _, v := range memberIDs {
		removeMembers[v] = struct{}{}
	}
	leftMemberIDs := make([]uuid.UUID, 0, len(belongingMembers.Data))
	for _, v := range belongingMembers.Data {
		if _, ok := removeMembers[v.Member.MemberID]; ok {
			continue
		}
		leftMemberIDs = append(leftMemberIDs, v.Member.MemberID)
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
		var action entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers
		var actAttr entity.ChatRoomAction
		if _, action, actAttr, err = removeMembersFromChatRoom(
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
		defer func(
			room entity.ChatRoom, membersIDs, leftMemberIDs []uuid.UUID,
			action entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers,
			actAttr entity.ChatRoomAction,
		) {
			if err == nil {
				m.WsHub.Dispatch(ws.EventTypeChatRoomRemovedMe, ws.Targets{
					Members: membersIDs,
				}, ws.ChatRoomRemovedMeEventData{
					ChatRoomID:           room.ChatRoomID,
					Action:               action,
					ChatRoomActionID:     actAttr.ChatRoomActionID,
					ChatRoomActionTypeID: actAttr.ChatRoomActionTypeID,
					ActedAt:              actAttr.ActedAt,
				})
				m.WsHub.Dispatch(ws.EventTypeChatRoomRemovedMember, ws.Targets{
					Members: leftMemberIDs,
				}, ws.ChatRoomRemovedMemberEventData{
					ChatRoomID:           room.ChatRoomID,
					Action:               action,
					ChatRoomActionID:     actAttr.ChatRoomActionID,
					ChatRoomActionTypeID: actAttr.ChatRoomActionTypeID,
					ActedAt:              actAttr.ActedAt,
				})
			}
		}(room, memberIDs, leftMemberIDs, action, actAttr)
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
			return 0, errhandle.NewModelNotFoundError(OrganizationBelongingTargetOwner)
		}
		return 0, fmt.Errorf("failed to find member: %w", err)
	}
	org, err := m.DB.FindOrganizationWithDetailWithSd(ctx, sd, organizationID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(OrganizationBelongingTargetOrganization)
		}
		return 0, fmt.Errorf("failed to find organization: %w", err)
	}
	if org.Grade.Valid {
		return 0, errhandle.NewCommonError(response.AttemptOperateGradeOrganization, nil)
	}
	if org.Group.Valid {
		return 0, errhandle.NewCommonError(response.AttemptOperateGroupOrganization, nil)
	}
	if org.IsPersonal {
		return 0, errhandle.NewCommonError(response.AttemptOperatePersonalOrganization, nil)
	}
	if org.IsWhole {
		return 0, errhandle.NewCommonError(response.AttemptOperateWholeOrganization, nil)
	}
	belongingMembers, err := m.DB.GetMembersOnOrganizationWithSd(
		ctx,
		sd,
		organizationID,
		parameter.WhereMemberOnOrganizationParam{},
		parameter.MemberOnOrganizationOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get members on organization: %w", err)
	}
	leftMemberIDs := make([]uuid.UUID, 0, len(belongingMembers.Data))
	for _, v := range belongingMembers.Data {
		if v.Member.MemberID == memberID {
			continue
		}
		leftMemberIDs = append(leftMemberIDs, v.Member.MemberID)
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
		var action entity.ChatRoomWithdrawActionWithMember
		var actAttr entity.ChatRoomAction
		if _, action, actAttr, err = withdrawMemberFromChatRoom(
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
		defer func(
			room entity.ChatRoom, leftMemberIDs []uuid.UUID, memberID uuid.UUID,
			action entity.ChatRoomWithdrawActionWithMember,
			actAttr entity.ChatRoomAction,
		) {
			if err == nil {
				m.WsHub.Dispatch(ws.EventTypeChatRoomWithdrawnMe, ws.Targets{
					Members: []uuid.UUID{memberID},
				}, ws.ChatRoomWithdrawnMeEventData{
					ChatRoomID:           room.ChatRoomID,
					Action:               action,
					ChatRoomActionID:     actAttr.ChatRoomActionID,
					ChatRoomActionTypeID: actAttr.ChatRoomActionTypeID,
					ActedAt:              actAttr.ActedAt,
				})
				m.WsHub.Dispatch(ws.EventTypeChatRoomWithdrawnMember, ws.Targets{
					Members: leftMemberIDs,
				}, ws.ChatRoomWithdrawnMemberEventData{
					ChatRoomID:           room.ChatRoomID,
					Action:               action,
					ChatRoomActionID:     actAttr.ChatRoomActionID,
					ChatRoomActionTypeID: actAttr.ChatRoomActionTypeID,
					ActedAt:              actAttr.ActedAt,
				})
			}
		}(room, leftMemberIDs, memberID, action, actAttr)
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
