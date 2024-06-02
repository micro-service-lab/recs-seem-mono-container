package pgadapter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/query"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

func convReadableMessageOnMember(e query.GetReadableMessagesOnMemberRow) entity.ReadableMessageOnMember {
	return entity.ReadableMessageOnMember{
		Message: entity.Message{
			MessageID:        e.MessageID,
			ChatRoomActionID: e.ChatRoomActionID,
			SenderID:         entity.UUID(e.SenderID),
			Body:             e.Body,
			PostedAt:         e.PostedAt,
			LastEditedAt:     e.LastEditedAt,
		},
		ReadAt: entity.Timestamptz{
			Time:             e.ReadAt.Time,
			InfinityModifier: entity.InfinityModifier(e.ReadAt.InfinityModifier),
			Valid:            e.ReadAt.Valid,
		},
	}
}

func convReadableMemberOnMessage(e query.GetReadableMembersOnMessageRow) entity.ReadableMemberOnMessage {
	return entity.ReadableMemberOnMessage{
		Member: entity.MemberCard{
			MemberID:  e.MemberID,
			Name:      e.Name,
			Email:     e.Email,
			FirstName: entity.String(e.FirstName),
			LastName:  entity.String(e.LastName),
			ProfileImage: entity.NullableEntity[entity.ImageWithAttachableItem]{
				Valid: e.ProfileImageID.Valid,
				Entity: entity.ImageWithAttachableItem{
					ImageID: e.ProfileImageID.Bytes,
					Height:  entity.Float(e.ProfileImageHeight),
					Width:   entity.Float(e.ProfileImageWidth),
					AttachableItem: entity.AttachableItem{
						AttachableItemID: e.ProfileImageAttachableItemID.Bytes,
						OwnerID:          entity.UUID(e.ProfileImageOwnerID),
						FromOuter:        e.ProfileImageFromOuter.Bool,
						URL:              e.ProfileImageUrl.String,
						Alias:            e.ProfileImageAlias.String,
						Size:             entity.Float(e.ProfileImageSize),
						MimeTypeID:       e.ProfileImageMimeTypeID.Bytes,
					},
				},
			},
			GradeID: e.GradeID,
			GroupID: e.GroupID,
		},
		ReadAt: entity.Timestamptz{
			Time:             e.ReadAt.Time,
			InfinityModifier: entity.InfinityModifier(e.ReadAt.InfinityModifier),
			Valid:            e.ReadAt.Valid,
		},
	}
}

func countReadableMembersOnMessage(
	ctx context.Context, qtx *query.Queries, messageID uuid.UUID, where parameter.WhereReadableMemberOnMessageParam,
) (int64, error) {
	p := query.CountReadableMembersOnMessageParams{
		MessageID:      messageID,
		WhereLikeName:  where.WhereLikeName,
		SearchName:     where.SearchName,
		WhereIsRead:    where.WhereIsRead,
		WhereIsNotRead: where.WhereIsNotRead,
	}
	c, err := qtx.CountReadableMembersOnMessage(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count readable members on message: %w", err)
	}
	return c, nil
}

// CountReadableMembersOnMessage メッセージ上のメンバー数を取得する。
func (a *PgAdapter) CountReadableMembersOnMessage(
	ctx context.Context, messageID uuid.UUID, where parameter.WhereReadableMemberOnMessageParam,
) (int64, error) {
	return countReadableMembersOnMessage(ctx, a.query, messageID, where)
}

// CountReadableMembersOnMessageWithSd SD付きでメッセージ上のメンバー数を取得する。
func (a *PgAdapter) CountReadableMembersOnMessageWithSd(
	ctx context.Context, sd store.Sd, messageID uuid.UUID, where parameter.WhereReadableMemberOnMessageParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countReadableMembersOnMessage(ctx, qtx, messageID, where)
}

func countReadableMessagesOnChatRoomAndMember(
	ctx context.Context, qtx *query.Queries, chatRoomID, memberID uuid.UUID,
	where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
) (int64, error) {
	p := query.CountReadableMessagesOnChatRoomAndMemberParams{
		ChatRoomID:     chatRoomID,
		MemberID:       memberID,
		WhereIsRead:    where.WhereIsRead,
		WhereIsNotRead: where.WhereIsNotRead,
	}
	c, err := qtx.CountReadableMessagesOnChatRoomAndMember(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count readable messages on readable message and member: %w", err)
	}
	return c, nil
}

// CountReadableMessagesOnChatRoomAndMember チャットルーム、メンバー上のメッセージ数を取得する。
func (a *PgAdapter) CountReadableMessagesOnChatRoomAndMember(
	ctx context.Context, chatRoomID, memberID uuid.UUID,
	where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
) (int64, error) {
	return countReadableMessagesOnChatRoomAndMember(ctx, a.query, chatRoomID, memberID, where)
}

// CountReadableMessagesOnChatRoomAndMemberWithSd SD付きでチャットルーム、メンバー上のメッセージ数を取得する。
func (a *PgAdapter) CountReadableMessagesOnChatRoomAndMemberWithSd(
	ctx context.Context, sd store.Sd, chatRoomID, memberID uuid.UUID,
	where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countReadableMessagesOnChatRoomAndMember(ctx, qtx, chatRoomID, memberID, where)
}

func countReadsOnMessages(
	ctx context.Context, qtx *query.Queries, messageIDs []uuid.UUID, where parameter.WhereReadsOnMessageParam,
) ([]entity.ReadReceiptGroupByMessage, error) {
	p := query.CountReadsOnMessagesParams{
		MessageIds:     messageIDs,
		WhereIsRead:    where.WhereIsRead,
		WhereIsNotRead: where.WhereIsNotRead,
	}
	rs, err := qtx.CountReadsOnMessages(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("failed to count reads on messages: %w", err)
	}
	res := make([]entity.ReadReceiptGroupByMessage, 0, len(rs))
	for _, r := range rs {
		res = append(res, entity.ReadReceiptGroupByMessage{
			MessageID: r.MessageID,
			Count:     r.Count,
		})
	}
	return res, nil
}

// CountReadsOnMessages メッセージ上の既読数を取得する。
func (a *PgAdapter) CountReadsOnMessages(
	ctx context.Context, messageIDs []uuid.UUID, where parameter.WhereReadsOnMessageParam,
) ([]entity.ReadReceiptGroupByMessage, error) {
	return countReadsOnMessages(ctx, a.query, messageIDs, where)
}

// CountReadsOnMessagesWithSd SD付きでメッセージ上の既読数を取得する。
func (a *PgAdapter) CountReadsOnMessagesWithSd(
	ctx context.Context, sd store.Sd, messageIDs []uuid.UUID, where parameter.WhereReadsOnMessageParam,
) ([]entity.ReadReceiptGroupByMessage, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return nil, store.ErrNotFoundDescriptor
	}
	return countReadsOnMessages(ctx, qtx, messageIDs, where)
}

func countReadableMessagesOnChatRooms(
	ctx context.Context, qtx *query.Queries, chatRoomIDs []uuid.UUID,
	where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
) ([]entity.ReadReceiptGroupByChatRoom, error) {
	p := query.CountReadableMessagesOnChatRoomsParams{
		ChatRoomIds:    chatRoomIDs,
		WhereIsRead:    where.WhereIsRead,
		WhereIsNotRead: where.WhereIsNotRead,
	}
	rs, err := qtx.CountReadableMessagesOnChatRooms(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("failed to count readable messages on chat rooms: %w", err)
	}
	res := make([]entity.ReadReceiptGroupByChatRoom, 0, len(rs))
	for _, r := range rs {
		res = append(res, entity.ReadReceiptGroupByChatRoom{
			ChatRoomID: r.ChatRoomID.Bytes,
			Count:      r.Count,
		})
	}
	return res, nil
}

// CountReadableMessagesOnChatRooms チャットルーム上のメッセージ数を取得する。
func (a *PgAdapter) CountReadableMessagesOnChatRooms(
	ctx context.Context, chatRoomIDs []uuid.UUID,
	where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
) ([]entity.ReadReceiptGroupByChatRoom, error) {
	return countReadableMessagesOnChatRooms(ctx, a.query, chatRoomIDs, where)
}

// CountReadableMessagesOnChatRoomsWithSd SD付きでチャットルーム上のメッセージ数を取得する。
func (a *PgAdapter) CountReadableMessagesOnChatRoomsWithSd(
	ctx context.Context, sd store.Sd, chatRoomIDs []uuid.UUID,
	where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
) ([]entity.ReadReceiptGroupByChatRoom, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return nil, store.ErrNotFoundDescriptor
	}
	return countReadableMessagesOnChatRooms(ctx, qtx, chatRoomIDs, where)
}

func countReadableMessagesOnChatRoomsAndMember(
	ctx context.Context, qtx *query.Queries, chatRoomIDs []uuid.UUID, memberID uuid.UUID,
	where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
) ([]entity.ReadReceiptGroupByChatRoom, error) {
	p := query.CountReadableMessagesOnChatRoomsAndMemberParams{
		ChatRoomIds:    chatRoomIDs,
		MemberID:       memberID,
		WhereIsRead:    where.WhereIsRead,
		WhereIsNotRead: where.WhereIsNotRead,
	}
	rs, err := qtx.CountReadableMessagesOnChatRoomsAndMember(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("failed to count readable message on readable message and member: %w", err)
	}
	res := make([]entity.ReadReceiptGroupByChatRoom, 0, len(rs))
	for _, r := range rs {
		res = append(res, entity.ReadReceiptGroupByChatRoom{
			ChatRoomID: r.ChatRoomID.Bytes,
			Count:      r.Count,
		})
	}
	return res, nil
}

// CountReadableMessagesOnChatRoomsAndMember 複数のチャットルーム、メンバー上のメッセージ数を取得する。
func (a *PgAdapter) CountReadableMessagesOnChatRoomsAndMember(
	ctx context.Context, chatRoomIDs []uuid.UUID, memberID uuid.UUID,
	where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
) ([]entity.ReadReceiptGroupByChatRoom, error) {
	return countReadableMessagesOnChatRoomsAndMember(ctx, a.query, chatRoomIDs, memberID, where)
}

// CountReadableMessagesOnChatRoomsAndMemberWithSd SD付きで複数のチャットルーム、メンバー上のメッセージ数を取得する。
func (a *PgAdapter) CountReadableMessagesOnChatRoomsAndMemberWithSd(
	ctx context.Context, sd store.Sd, chatRoomIDs []uuid.UUID, memberID uuid.UUID,
	where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
) ([]entity.ReadReceiptGroupByChatRoom, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return nil, store.ErrNotFoundDescriptor
	}
	return countReadableMessagesOnChatRoomsAndMember(ctx, qtx, chatRoomIDs, memberID, where)
}

func countReadableMessagesOnMember(
	ctx context.Context, qtx *query.Queries, memberID uuid.UUID,
	where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
) (int64, error) {
	p := query.CountReadableMessagesOnMemberParams{
		MemberID:       memberID,
		WhereIsRead:    where.WhereIsRead,
		WhereIsNotRead: where.WhereIsNotRead,
	}
	c, err := qtx.CountReadableMessagesOnMember(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count readable messages on member: %w", err)
	}
	return c, nil
}

// CountReadableMessagesOnMember メンバー上のメッセージ数を取得する。
func (a *PgAdapter) CountReadableMessagesOnMember(
	ctx context.Context, memberID uuid.UUID,
	where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
) (int64, error) {
	return countReadableMessagesOnMember(ctx, a.query, memberID, where)
}

// CountReadableMessagesOnMemberWithSd SD付きでメンバー上のメッセージ数を取得する。
func (a *PgAdapter) CountReadableMessagesOnMemberWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
	where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countReadableMessagesOnMember(ctx, qtx, memberID, where)
}

func createReadReceipt(
	ctx context.Context, qtx *query.Queries, param parameter.CreateReadReceiptParam,
) (entity.ReadReceipt, error) {
	p := query.CreateReadReceiptParams{
		MessageID: param.MessageID,
		MemberID:  param.MemberID,
		ReadAt: pgtype.Timestamptz{
			Time:             param.ReadAt.Time,
			InfinityModifier: pgtype.InfinityModifier(param.ReadAt.InfinityModifier),
			Valid:            param.ReadAt.Valid,
		},
	}
	r, err := qtx.CreateReadReceipt(ctx, p)
	if err != nil {
		return entity.ReadReceipt{}, fmt.Errorf("failed to create read receipt: %w", err)
	}
	return entity.ReadReceipt{
		MessageID: r.MessageID,
		MemberID:  r.MemberID,
		ReadAt: entity.Timestamptz{
			Time:             r.ReadAt.Time,
			InfinityModifier: entity.InfinityModifier(r.ReadAt.InfinityModifier),
			Valid:            r.ReadAt.Valid,
		},
	}, nil
}

// CreateReadReceipt 既読情報を作成する。
func (a *PgAdapter) CreateReadReceipt(
	ctx context.Context, param parameter.CreateReadReceiptParam,
) (entity.ReadReceipt, error) {
	return createReadReceipt(ctx, a.query, param)
}

// CreateReadReceiptWithSd SD付きで既読情報を作成する。
func (a *PgAdapter) CreateReadReceiptWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateReadReceiptParam,
) (entity.ReadReceipt, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ReadReceipt{}, store.ErrNotFoundDescriptor
	}
	return createReadReceipt(ctx, qtx, param)
}

func createReadReceipts(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateReadReceiptParam,
) (int64, error) {
	ps := make([]query.CreateReadReceiptsParams, 0, len(params))
	for _, p := range params {
		ps = append(ps, query.CreateReadReceiptsParams{
			MessageID: p.MessageID,
			MemberID:  p.MemberID,
			ReadAt: pgtype.Timestamptz{
				Time:             p.ReadAt.Time,
				InfinityModifier: pgtype.InfinityModifier(p.ReadAt.InfinityModifier),
				Valid:            p.ReadAt.Valid,
			},
		})
	}
	c, err := qtx.CreateReadReceipts(ctx, ps)
	if err != nil {
		return 0, fmt.Errorf("failed to create read receipts: %w", err)
	}
	return c, nil
}

// CreateReadReceipts 複数の既読情報を作成する。
func (a *PgAdapter) CreateReadReceipts(
	ctx context.Context, params []parameter.CreateReadReceiptParam,
) (int64, error) {
	return createReadReceipts(ctx, a.query, params)
}

// CreateReadReceiptsWithSd SD付きで複数の既読情報を作成する。
func (a *PgAdapter) CreateReadReceiptsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateReadReceiptParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createReadReceipts(ctx, qtx, params)
}

func readReceipt(
	ctx context.Context, qtx *query.Queries, param parameter.ReadReceiptParam,
) (entity.ReadReceipt, error) {
	p := query.ReadReceiptParams{
		MessageID: param.MessageID,
		MemberID:  param.MemberID,
		ReadAt: pgtype.Timestamptz{
			Time:             param.ReadAt.Time,
			InfinityModifier: pgtype.InfinityModifier(param.ReadAt.InfinityModifier),
			Valid:            param.ReadAt.Valid,
		},
	}
	r, err := qtx.ReadReceipt(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ReadReceipt{}, errhandle.NewModelNotFoundError("read receipt")
		}
		return entity.ReadReceipt{}, fmt.Errorf("failed to mark readable members on message: %w", err)
	}
	return entity.ReadReceipt{
		MessageID: r.MessageID,
		MemberID:  r.MemberID,
		ReadAt: entity.Timestamptz{
			Time:             r.ReadAt.Time,
			InfinityModifier: entity.InfinityModifier(r.ReadAt.InfinityModifier),
			Valid:            r.ReadAt.Valid,
		},
	}, nil
}

// ReadReceipt 既読にする。
func (a *PgAdapter) ReadReceipt(
	ctx context.Context, param parameter.ReadReceiptParam,
) (entity.ReadReceipt, error) {
	return readReceipt(ctx, a.query, param)
}

// ReadReceiptWithSd SD付きで既読にする。
func (a *PgAdapter) ReadReceiptWithSd(
	ctx context.Context, sd store.Sd, param parameter.ReadReceiptParam,
) (entity.ReadReceipt, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ReadReceipt{}, store.ErrNotFoundDescriptor
	}
	return readReceipt(ctx, qtx, param)
}

func readReceipts(
	ctx context.Context, qtx *query.Queries, param parameter.ReadReceiptsParam,
) (int64, error) {
	p := query.ReadReceiptsParams{
		MessageIds: param.MessageIDs,
		MemberID:   param.MemberID,
		ReadAt: pgtype.Timestamptz{
			Time:             param.ReadAt.Time,
			InfinityModifier: pgtype.InfinityModifier(param.ReadAt.InfinityModifier),
			Valid:            param.ReadAt.Valid,
		},
	}
	c, err := qtx.ReadReceipts(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to mark readable members on message: %w", err)
	}
	if c != int64(len(param.MessageIDs)) {
		return 0, errhandle.NewModelNotFoundError("read receipt")
	}
	return c, nil
}

// ReadReceipts 複数既読にする。
func (a *PgAdapter) ReadReceipts(
	ctx context.Context, param parameter.ReadReceiptsParam,
) (int64, error) {
	return readReceipts(ctx, a.query, param)
}

// ReadReceiptsWithSd SD付きで複数既読にする。
func (a *PgAdapter) ReadReceiptsWithSd(
	ctx context.Context, sd store.Sd, param parameter.ReadReceiptsParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return readReceipts(ctx, qtx, param)
}

func readReceiptsOnMember(
	ctx context.Context, qtx *query.Queries, memberID uuid.UUID, readAt time.Time,
) (int64, error) {
	p := query.ReadReceiptsOnMemberParams{
		MemberID: memberID,
		ReadAt:   pgtype.Timestamptz{Time: readAt, Valid: true},
	}
	c, err := qtx.ReadReceiptsOnMember(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to mark readable members on message: %w", err)
	}
	return c, nil
}

// ReadReceiptsOnMember メンバー上の既読をする。
func (a *PgAdapter) ReadReceiptsOnMember(
	ctx context.Context, memberID uuid.UUID, readAt time.Time,
) (int64, error) {
	return readReceiptsOnMember(ctx, a.query, memberID, readAt)
}

// ReadReceiptsOnMemberWithSd SD付きでメンバー上の既読をする。
func (a *PgAdapter) ReadReceiptsOnMemberWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID, readAt time.Time,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return readReceiptsOnMember(ctx, qtx, memberID, readAt)
}

func readReceiptsOnChatRoomAndMember(
	ctx context.Context,
	qtx *query.Queries,
	chatRoomID, memberID uuid.UUID,
	readAt time.Time,
) (int64, error) {
	p := query.ReadReceiptsOnChatRoomAndMemberParams{
		ChatRoomID: chatRoomID,
		MemberID:   memberID,
		ReadAt:     pgtype.Timestamptz{Time: readAt, Valid: true},
	}
	c, err := qtx.ReadReceiptsOnChatRoomAndMember(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to mark readable members on message: %w", err)
	}
	return c, nil
}

// ReadReceiptsOnChatRoomAndMember チャットルーム、メンバー上の既読をする。
func (a *PgAdapter) ReadReceiptsOnChatRoomAndMember(
	ctx context.Context,
	chatRoomID, memberID uuid.UUID,
	readAt time.Time,
) (int64, error) {
	return readReceiptsOnChatRoomAndMember(ctx, a.query, chatRoomID, memberID, readAt)
}

// ReadReceiptsOnChatRoomAndMemberWithSd SD付きでチャットルーム、メンバー上の既読をする。
func (a *PgAdapter) ReadReceiptsOnChatRoomAndMemberWithSd(
	ctx context.Context,
	sd store.Sd,
	chatRoomID, memberID uuid.UUID,
	readAt time.Time,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return readReceiptsOnChatRoomAndMember(ctx, qtx, chatRoomID, memberID, readAt)
}

func existsReadReceipt(
	ctx context.Context, qtx *query.Queries, memberID, messageID uuid.UUID, where parameter.WhereExistsReadReceiptParam,
) (bool, error) {
	p := query.ExistsReadReceiptParams{
		MemberID:       memberID,
		MessageID:      messageID,
		WhereIsRead:    where.WhereIsRead,
		WhereIsNotRead: where.WhereIsNotRead,
	}
	r, err := qtx.ExistsReadReceipt(ctx, p)
	if err != nil {
		return false, fmt.Errorf("failed to check read receipt: %w", err)
	}
	return r, nil
}

// ExistsReadReceipt 既読情報が存在するか確認する。
func (a *PgAdapter) ExistsReadReceipt(
	ctx context.Context, memberID, messageID uuid.UUID, where parameter.WhereExistsReadReceiptParam,
) (bool, error) {
	return existsReadReceipt(ctx, a.query, memberID, messageID, where)
}

// ExistsReadReceiptWithSd SD付きで既読情報が存在するか確認する。
func (a *PgAdapter) ExistsReadReceiptWithSd(
	ctx context.Context, sd store.Sd, memberID, messageID uuid.UUID, where parameter.WhereExistsReadReceiptParam,
) (bool, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return false, store.ErrNotFoundDescriptor
	}
	return existsReadReceipt(ctx, qtx, memberID, messageID, where)
}

func findReadReceipt(
	ctx context.Context, qtx *query.Queries, memberID, messageID uuid.UUID,
) (entity.ReadReceipt, error) {
	p := query.FindReadReceiptParams{
		MemberID:  memberID,
		MessageID: messageID,
	}
	r, err := qtx.FindReadReceipt(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ReadReceipt{}, errhandle.NewModelNotFoundError("read receipt")
		}
		return entity.ReadReceipt{}, fmt.Errorf("failed to find read receipt: %w", err)
	}
	return entity.ReadReceipt{
		MessageID: r.MessageID,
		MemberID:  r.MemberID,
		ReadAt: entity.Timestamptz{
			Time:             r.ReadAt.Time,
			InfinityModifier: entity.InfinityModifier(r.ReadAt.InfinityModifier),
			Valid:            r.ReadAt.Valid,
		},
	}, nil
}

// FindReadReceipt 既読情報を取得する。
func (a *PgAdapter) FindReadReceipt(
	ctx context.Context, memberID, messageID uuid.UUID,
) (entity.ReadReceipt, error) {
	return findReadReceipt(ctx, a.query, memberID, messageID)
}

// FindReadReceiptWithSd SD付きで既読情報を取得する。
func (a *PgAdapter) FindReadReceiptWithSd(
	ctx context.Context, sd store.Sd, memberID, messageID uuid.UUID,
) (entity.ReadReceipt, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ReadReceipt{}, store.ErrNotFoundDescriptor
	}
	return findReadReceipt(ctx, qtx, memberID, messageID)
}

func getReadableMessagesOnMember(
	ctx context.Context,
	qtx *query.Queries,
	memberID uuid.UUID,
	where parameter.WhereReadableMessageOnMemberParam,
	order parameter.ReadableMessageOnMemberOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ReadableMessageOnMember], error) {
	eConvFunc := func(e entity.ReadableMessageOnMemberForQuery) (entity.ReadableMessageOnMember, error) {
		return e.ReadableMessageOnMember, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountReadableMessagesOnMemberParams{
			MemberID:       memberID,
			WhereIsRead:    where.WhereIsRead,
			WhereIsNotRead: where.WhereIsNotRead,
		}
		r, err := qtx.CountReadableMessagesOnMember(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count readable messages on member: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.ReadableMessageOnMemberForQuery, error) {
		p := query.GetReadableMessagesOnMemberParams{
			MemberID:       memberID,
			WhereIsRead:    where.WhereIsRead,
			WhereIsNotRead: where.WhereIsNotRead,
			OrderMethod:    orderMethod,
		}
		r, err := qtx.GetReadableMessagesOnMember(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.ReadableMessageOnMemberForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get readable messages on member: %w", err)
		}
		fq := make([]entity.ReadableMessageOnMemberForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.ReadableMessageOnMemberForQuery{
				Pkey:                    entity.Int(e.TMessagesPkey),
				ReadableMessageOnMember: convReadableMessageOnMember(query.GetReadableMessagesOnMemberRow(e)),
			}
		}
		return fq, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.ReadableMessageOnMemberForQuery, error) {
		var readCursor time.Time
		var err error
		switch subCursor {
		case parameter.ReadableMessageOnMemberReadAtCursorKey:
			cv, ok := subCursorValue.(string)
			readCursor, err = time.Parse(time.RFC3339, cv)
			if !ok || err != nil {
				readCursor = time.Time{}
			}
		}
		p := query.GetReadableMessagesOnMemberUseKeysetPaginateParams{
			MemberID:        memberID,
			WhereIsRead:     where.WhereIsRead,
			WhereIsNotRead:  where.WhereIsNotRead,
			OrderMethod:     orderMethod,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			Limit:           limit,
			ReadAtCursor: pgtype.Timestamptz{
				Time:  readCursor,
				Valid: !readCursor.IsZero(),
			},
		}
		r, err := qtx.GetReadableMessagesOnMemberUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get readable messages on member: %w", err)
		}
		fq := make([]entity.ReadableMessageOnMemberForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.ReadableMessageOnMemberForQuery{
				Pkey:                    entity.Int(e.TMessagesPkey),
				ReadableMessageOnMember: convReadableMessageOnMember(query.GetReadableMessagesOnMemberRow(e)),
			}
		}
		return fq, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.ReadableMessageOnMemberForQuery, error) {
		p := query.GetReadableMessagesOnMemberUseNumberedPaginateParams{
			MemberID:       memberID,
			WhereIsRead:    where.WhereIsRead,
			WhereIsNotRead: where.WhereIsNotRead,
			OrderMethod:    orderMethod,
			Limit:          limit,
			Offset:         offset,
		}
		r, err := qtx.GetReadableMessagesOnMemberUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get readable messages on member: %w", err)
		}
		fq := make([]entity.ReadableMessageOnMemberForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.ReadableMessageOnMemberForQuery{
				Pkey:                    entity.Int(e.TMessagesPkey),
				ReadableMessageOnMember: convReadableMessageOnMember(query.GetReadableMessagesOnMemberRow(e)),
			}
		}
		return fq, nil
	}
	selector := func(subCursor string, e entity.ReadableMessageOnMemberForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.ReadableMessageOnMemberDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.ReadableMessageOnMemberReadAtCursorKey:
			return entity.Int(e.Pkey), e.ReadableMessageOnMember.ReadAt
		}
		return entity.Int(e.Pkey), nil
	}

	res, err := store.RunListQuery(
		ctx,
		order,
		np,
		cp,
		wc,
		eConvFunc,
		runCFunc,
		runQFunc,
		runQCPFunc,
		runQNPFunc,
		selector,
	)
	if err != nil {
		return store.ListResult[entity.ReadableMessageOnMember]{},
			fmt.Errorf("failed to get readable messages on member: %w", err)
	}
	return res, nil
}

// GetReadableMessagesOnMember メンバー上のメッセージを取得する。
func (a *PgAdapter) GetReadableMessagesOnMember(
	ctx context.Context, memberID uuid.UUID, where parameter.WhereReadableMessageOnMemberParam,
	order parameter.ReadableMessageOnMemberOrderMethod,
	np store.NumberedPaginationParam, cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.ReadableMessageOnMember], error) {
	return getReadableMessagesOnMember(ctx, a.query, memberID, where, order, np, cp, wc)
}

// GetReadableMessagesOnMemberWithSd SD付きでメンバー上のメッセージを取得する。
func (a *PgAdapter) GetReadableMessagesOnMemberWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
	where parameter.WhereReadableMessageOnMemberParam, order parameter.ReadableMessageOnMemberOrderMethod,
	np store.NumberedPaginationParam, cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.ReadableMessageOnMember], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ReadableMessageOnMember]{}, store.ErrNotFoundDescriptor
	}
	return getReadableMessagesOnMember(ctx, qtx, memberID, where, order, np, cp, wc)
}

func getPluralReadableMessagesOnMember(
	ctx context.Context, qtx *query.Queries, memberIDs []uuid.UUID,
	orderMethod parameter.ReadableMessageOnMemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ReadableMessageOnMember], error) {
	var e []query.GetPluralReadableMessagesOnMemberRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralReadableMessagesOnMember(ctx, query.GetPluralReadableMessagesOnMemberParams{
			MemberIds:   memberIDs,
			OrderMethod: orderMethod.GetStringValue(),
		})
	} else {
		var qe []query.GetPluralReadableMessagesOnMemberUseNumberedPaginateRow
		qe, err = qtx.GetPluralReadableMessagesOnMemberUseNumberedPaginate(
			ctx, query.GetPluralReadableMessagesOnMemberUseNumberedPaginateParams{
				MemberIds:   memberIDs,
				Limit:       int32(np.Limit.Int64),
				Offset:      int32(np.Offset.Int64),
				OrderMethod: orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralReadableMessagesOnMemberRow, len(qe))
		for i, v := range qe {
			e[i] = query.GetPluralReadableMessagesOnMemberRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ReadableMessageOnMember]{},
			fmt.Errorf("failed to get readable messages on member: %w", err)
	}
	entities := make([]entity.ReadableMessageOnMember, len(e))
	for i, v := range e {
		entities[i] = convReadableMessageOnMember(query.GetReadableMessagesOnMemberRow(v))
	}
	return store.ListResult[entity.ReadableMessageOnMember]{Data: entities}, nil
}

// GetPluralReadableMessagesOnMember メンバー上の複数のメッセージを取得する。
func (a *PgAdapter) GetPluralReadableMessagesOnMember(
	ctx context.Context, memberIDs []uuid.UUID,
	np store.NumberedPaginationParam, order parameter.ReadableMessageOnMemberOrderMethod,
) (store.ListResult[entity.ReadableMessageOnMember], error) {
	return getPluralReadableMessagesOnMember(ctx, a.query, memberIDs, order, np)
}

// GetPluralReadableMessagesOnMemberWithSd SD付きでメンバー上の複数のメッセージを取得する。
func (a *PgAdapter) GetPluralReadableMessagesOnMemberWithSd(
	ctx context.Context, sd store.Sd, memberIDs []uuid.UUID,
	np store.NumberedPaginationParam, order parameter.ReadableMessageOnMemberOrderMethod,
) (store.ListResult[entity.ReadableMessageOnMember], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ReadableMessageOnMember]{}, store.ErrNotFoundDescriptor
	}
	return getPluralReadableMessagesOnMember(ctx, qtx, memberIDs, order, np)
}

func getReadableMembersOnMessage(
	ctx context.Context,
	qtx *query.Queries,
	messageID uuid.UUID,
	where parameter.WhereReadableMemberOnMessageParam,
	order parameter.ReadableMemberOnMessageOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ReadableMemberOnMessage], error) {
	eConvFunc := func(e entity.ReadableMemberOnMessageForQuery) (entity.ReadableMemberOnMessage, error) {
		return e.ReadableMemberOnMessage, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountReadableMembersOnMessageParams{
			MessageID:      messageID,
			WhereLikeName:  where.WhereLikeName,
			SearchName:     where.SearchName,
			WhereIsRead:    where.WhereIsRead,
			WhereIsNotRead: where.WhereIsNotRead,
		}
		r, err := qtx.CountReadableMembersOnMessage(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count readable members on message: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.ReadableMemberOnMessageForQuery, error) {
		p := query.GetReadableMembersOnMessageParams{
			MessageID:      messageID,
			WhereLikeName:  where.WhereLikeName,
			SearchName:     where.SearchName,
			WhereIsRead:    where.WhereIsRead,
			WhereIsNotRead: where.WhereIsNotRead,
			OrderMethod:    orderMethod,
		}
		r, err := qtx.GetReadableMembersOnMessage(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.ReadableMemberOnMessageForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get readable members on message: %w", err)
		}
		fq := make([]entity.ReadableMemberOnMessageForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.ReadableMemberOnMessageForQuery{
				Pkey:                    entity.Int(e.MMembersPkey),
				ReadableMemberOnMessage: convReadableMemberOnMessage(query.GetReadableMembersOnMessageRow(e)),
			}
		}
		return fq, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.ReadableMemberOnMessageForQuery, error) {
		var nameCursor string
		var readCursor time.Time
		var err error
		switch subCursor {
		case parameter.ReadableMemberOnMessageNameCursorKey:
			cv, ok := subCursorValue.(string)
			nameCursor = cv
			if !ok {
				nameCursor = ""
			}
		case parameter.ReadableMemberOnMessageReadAtCursorKey:
			cv, ok := subCursorValue.(string)
			readCursor, err = time.Parse(time.RFC3339, cv)
			if !ok || err != nil {
				readCursor = time.Time{}
			}
		}
		p := query.GetReadableMembersOnMessageUseKeysetPaginateParams{
			MessageID:       messageID,
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			WhereIsRead:     where.WhereIsRead,
			WhereIsNotRead:  where.WhereIsNotRead,
			OrderMethod:     orderMethod,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			Limit:           limit,
			NameCursor:      nameCursor,
			ReadAtCursor: pgtype.Timestamptz{
				Time:  readCursor,
				Valid: !readCursor.IsZero(),
			},
		}
		r, err := qtx.GetReadableMembersOnMessageUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get readable members on message: %w", err)
		}
		fq := make([]entity.ReadableMemberOnMessageForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.ReadableMemberOnMessageForQuery{
				Pkey:                    entity.Int(e.MMembersPkey),
				ReadableMemberOnMessage: convReadableMemberOnMessage(query.GetReadableMembersOnMessageRow(e)),
			}
		}
		return fq, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.ReadableMemberOnMessageForQuery, error) {
		p := query.GetReadableMembersOnMessageUseNumberedPaginateParams{
			MessageID:      messageID,
			WhereLikeName:  where.WhereLikeName,
			SearchName:     where.SearchName,
			WhereIsRead:    where.WhereIsRead,
			WhereIsNotRead: where.WhereIsNotRead,
			OrderMethod:    orderMethod,
			Limit:          limit,
			Offset:         offset,
		}
		r, err := qtx.GetReadableMembersOnMessageUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get readable members on message: %w", err)
		}
		fq := make([]entity.ReadableMemberOnMessageForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.ReadableMemberOnMessageForQuery{
				Pkey:                    entity.Int(e.MMembersPkey),
				ReadableMemberOnMessage: convReadableMemberOnMessage(query.GetReadableMembersOnMessageRow(e)),
			}
		}
		return fq, nil
	}
	selector := func(subCursor string, e entity.ReadableMemberOnMessageForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.ReadableMemberOnMessageDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.ReadableMemberOnMessageNameCursorKey:
			return entity.Int(e.Pkey), e.ReadableMemberOnMessage.Member.Name
		case parameter.ReadableMemberOnMessageReadAtCursorKey:
			return entity.Int(e.Pkey), e.ReadableMemberOnMessage.ReadAt
		}
		return entity.Int(e.Pkey), nil
	}

	res, err := store.RunListQuery(
		ctx,
		order,
		np,
		cp,
		wc,
		eConvFunc,
		runCFunc,
		runQFunc,
		runQCPFunc,
		runQNPFunc,
		selector,
	)
	if err != nil {
		return store.ListResult[entity.ReadableMemberOnMessage]{},
			fmt.Errorf("failed to get readable members on message: %w", err)
	}
	return res, nil
}

// GetReadableMembersOnMessage メッセージ上のメンバーを取得する。
func (a *PgAdapter) GetReadableMembersOnMessage(
	ctx context.Context, memberID uuid.UUID, where parameter.WhereReadableMemberOnMessageParam,
	order parameter.ReadableMemberOnMessageOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.ReadableMemberOnMessage], error) {
	return getReadableMembersOnMessage(ctx, a.query, memberID, where, order, np, cp, wc)
}

// GetReadableMembersOnMessageWithSd SD付きでメッセージ上のメンバーを取得する。
func (a *PgAdapter) GetReadableMembersOnMessageWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
	where parameter.WhereReadableMemberOnMessageParam, order parameter.ReadableMemberOnMessageOrderMethod,
	np store.NumberedPaginationParam, cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.ReadableMemberOnMessage], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ReadableMemberOnMessage]{}, store.ErrNotFoundDescriptor
	}
	return getReadableMembersOnMessage(ctx, qtx, memberID, where, order, np, cp, wc)
}

func getPluralReadableMembersOnMessage(
	ctx context.Context, qtx *query.Queries, messageIDs []uuid.UUID,
	orderMethod parameter.ReadableMemberOnMessageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ReadableMemberOnMessage], error) {
	var e []query.GetPluralReadableMembersOnMessageRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralReadableMembersOnMessage(ctx, query.GetPluralReadableMembersOnMessageParams{
			MessageIds:  messageIDs,
			OrderMethod: orderMethod.GetStringValue(),
		})
	} else {
		var qe []query.GetPluralReadableMembersOnMessageUseNumberedPaginateRow
		qe, err = qtx.GetPluralReadableMembersOnMessageUseNumberedPaginate(
			ctx, query.GetPluralReadableMembersOnMessageUseNumberedPaginateParams{
				MessageIds:  messageIDs,
				Limit:       int32(np.Limit.Int64),
				Offset:      int32(np.Offset.Int64),
				OrderMethod: orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralReadableMembersOnMessageRow, len(qe))
		for i, v := range qe {
			e[i] = query.GetPluralReadableMembersOnMessageRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ReadableMemberOnMessage]{},
			fmt.Errorf("failed to get readable members on message: %w", err)
	}
	entities := make([]entity.ReadableMemberOnMessage, len(e))
	for i, v := range e {
		entities[i] = convReadableMemberOnMessage(query.GetReadableMembersOnMessageRow(v))
	}
	return store.ListResult[entity.ReadableMemberOnMessage]{Data: entities}, nil
}

// GetPluralReadableMembersOnMessage メッセージ上の複数のメンバーを取得する。
func (a *PgAdapter) GetPluralReadableMembersOnMessage(
	ctx context.Context, messageIDs []uuid.UUID, np store.NumberedPaginationParam,
	order parameter.ReadableMemberOnMessageOrderMethod,
) (store.ListResult[entity.ReadableMemberOnMessage], error) {
	return getPluralReadableMembersOnMessage(ctx, a.query, messageIDs, order, np)
}

// GetPluralReadableMembersOnMessageWithSd SD付きでメッセージ上の複数のメンバーを取得する。
func (a *PgAdapter) GetPluralReadableMembersOnMessageWithSd(
	ctx context.Context, sd store.Sd, messageIDs []uuid.UUID,
	np store.NumberedPaginationParam, order parameter.ReadableMemberOnMessageOrderMethod,
) (store.ListResult[entity.ReadableMemberOnMessage], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ReadableMemberOnMessage]{}, store.ErrNotFoundDescriptor
	}
	return getPluralReadableMembersOnMessage(ctx, qtx, messageIDs, order, np)
}
