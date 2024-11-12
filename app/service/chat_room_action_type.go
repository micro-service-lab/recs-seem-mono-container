package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// ChatRoomActionTypeKey チャットルームアクションタイプキー。
type ChatRoomActionTypeKey string

const (
	// ChatRoomActionTypeKeyCreate チャットルームアクションタイプキー: 作成。
	ChatRoomActionTypeKeyCreate ChatRoomActionTypeKey = "create"
	// ChatRoomActionTypeKeyUpdateName チャットルームアクションタイプキー: 名前更新。
	ChatRoomActionTypeKeyUpdateName ChatRoomActionTypeKey = "update_name"
	// ChatRoomActionTypeKeyAddMember チャットルームアクションタイプキー: メンバー追加。
	ChatRoomActionTypeKeyAddMember ChatRoomActionTypeKey = "add_member"
	// ChatRoomActionTypeKeyRemoveMember チャットルームアクションタイプキー: メンバー削除。
	ChatRoomActionTypeKeyRemoveMember ChatRoomActionTypeKey = "remove_member"
	// ChatRoomActionTypeKeyWithdraw チャットルームアクションタイプキー: 退室。
	ChatRoomActionTypeKeyWithdraw ChatRoomActionTypeKey = "withdraw"
	// ChatRoomActionTypeKeyMessage チャットルームアクションタイプキー: メッセージ。
	ChatRoomActionTypeKeyMessage ChatRoomActionTypeKey = "message"
	// ChatRoomActionTypeKeyDeleteMessage チャットルームアクションタイプキー: メッセージ削除。
	ChatRoomActionTypeKeyDeleteMessage ChatRoomActionTypeKey = "delete_message"
)

// ChatRoomActionType チャットルームアクションタイプ。
type ChatRoomActionType struct {
	Key  string
	Name string
}

// ChatRoomActionTypes チャットルームアクションタイプ一覧。
var ChatRoomActionTypes = []ChatRoomActionType{
	{
		Key:  string(ChatRoomActionTypeKeyCreate),
		Name: "作成",
	},
	{
		Key:  string(ChatRoomActionTypeKeyUpdateName),
		Name: "名前更新",
	},
	{
		Key:  string(ChatRoomActionTypeKeyAddMember),
		Name: "メンバー追加",
	},
	{
		Key:  string(ChatRoomActionTypeKeyRemoveMember),
		Name: "メンバー削除",
	},
	{
		Key:  string(ChatRoomActionTypeKeyWithdraw),
		Name: "退室",
	},
	{
		Key:  string(ChatRoomActionTypeKeyMessage),
		Name: "メッセージ",
	},
	{
		Key:  string(ChatRoomActionTypeKeyDeleteMessage),
		Name: "メッセージ削除",
	},
}

// ManageChatRoomActionType チャットルームアクションタイプ管理サービス。
type ManageChatRoomActionType struct {
	DB store.Store
}

// CreateChatRoomActionType チャットルームアクションタイプを作成する。
func (m *ManageChatRoomActionType) CreateChatRoomActionType(
	ctx context.Context,
	name, key string,
) (entity.ChatRoomActionType, error) {
	p := parameter.CreateChatRoomActionTypeParam{
		Name: name,
		Key:  key,
	}
	e, err := m.DB.CreateChatRoomActionType(ctx, p)
	if err != nil {
		return entity.ChatRoomActionType{}, fmt.Errorf("failed to create chat room action type: %w", err)
	}
	return e, nil
}

// CreateChatRoomActionTypes チャットルームアクションタイプを複数作成する。
func (m *ManageChatRoomActionType) CreateChatRoomActionTypes(
	ctx context.Context, ps []parameter.CreateChatRoomActionTypeParam,
) (int64, error) {
	es, err := m.DB.CreateChatRoomActionTypes(ctx, ps)
	if err != nil {
		return 0, fmt.Errorf("failed to create chat room action types: %w", err)
	}
	return es, nil
}

// UpdateChatRoomActionType チャットルームアクションタイプを更新する。
func (m *ManageChatRoomActionType) UpdateChatRoomActionType(
	ctx context.Context, id uuid.UUID, name, key string,
) (entity.ChatRoomActionType, error) {
	p := parameter.UpdateChatRoomActionTypeParams{
		Name: name,
		Key:  key,
	}
	e, err := m.DB.UpdateChatRoomActionType(ctx, id, p)
	if err != nil {
		return entity.ChatRoomActionType{}, fmt.Errorf("failed to update chat room action type: %w", err)
	}
	return e, nil
}

// DeleteChatRoomActionType チャットルームアクションタイプを削除する。
func (m *ManageChatRoomActionType) DeleteChatRoomActionType(ctx context.Context, id uuid.UUID) (int64, error) {
	c, err := m.DB.DeleteChatRoomActionType(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room action type: %w", err)
	}
	return c, nil
}

// PluralDeleteChatRoomActionTypes チャットルームアクションタイプを複数削除する。
func (m *ManageChatRoomActionType) PluralDeleteChatRoomActionTypes(
	ctx context.Context, ids []uuid.UUID,
) (int64, error) {
	c, err := m.DB.PluralDeleteChatRoomActionTypes(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete chat room action types: %w", err)
	}
	return c, nil
}

// FindChatRoomActionTypeByID チャットルームアクションタイプをIDで取得する。
func (m *ManageChatRoomActionType) FindChatRoomActionTypeByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.ChatRoomActionType, error) {
	e, err := m.DB.FindChatRoomActionTypeByID(ctx, id)
	if err != nil {
		return entity.ChatRoomActionType{}, fmt.Errorf("failed to find chat room action type by id: %w", err)
	}
	return e, nil
}

// FindChatRoomActionTypeByKey チャットルームアクションタイプをキーで取得する。
func (m *ManageChatRoomActionType) FindChatRoomActionTypeByKey(
	ctx context.Context, key string,
) (entity.ChatRoomActionType, error) {
	e, err := m.DB.FindChatRoomActionTypeByKey(ctx, key)
	if err != nil {
		return entity.ChatRoomActionType{}, fmt.Errorf("failed to find chat room action type by key: %w", err)
	}
	return e, nil
}

// GetChatRoomActionTypes チャットルームアクションタイプを取得する。
func (m *ManageChatRoomActionType) GetChatRoomActionTypes(
	ctx context.Context,
	whereSearchName string,
	order parameter.ChatRoomActionTypeOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.ChatRoomActionType], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereChatRoomActionTypeParam{
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
	r, err := m.DB.GetChatRoomActionTypes(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.ChatRoomActionType]{}, fmt.Errorf("failed to get chat room action types: %w", err)
	}
	return r, nil
}

// GetChatRoomActionTypesCount チャットルームアクションタイプの数を取得する。
func (m *ManageChatRoomActionType) GetChatRoomActionTypesCount(
	ctx context.Context,
	whereSearchName string,
) (int64, error) {
	p := parameter.WhereChatRoomActionTypeParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	c, err := m.DB.CountChatRoomActionTypes(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get chat room action types count: %w", err)
	}
	return c, nil
}
