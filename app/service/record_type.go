package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// RecordTypeKey 議事録タイプキー。
type RecordTypeKey string

const (
	// RecordTypeKeyMeeting ミーティング議事録。
	RecordTypeKeyMeeting RecordTypeKey = "meeting"
)

// RecordType 議事録タイプ。
type RecordType struct {
	Key   string
	Name  string
	Color string
}

// RecordTypes 議事録タイプ一覧。
var RecordTypes = []RecordType{
	{Key: string(RecordTypeKeyMeeting), Name: "ミーティング", Color: "#FFB866"},
}

// ManageRecordType 議事録タイプ管理サービス。
type ManageRecordType struct {
	DB store.Store
}

// CreateRecordType 議事録タイプを作成する。
func (m *ManageRecordType) CreateRecordType(
	ctx context.Context,
	name, key string,
) (entity.RecordType, error) {
	p := parameter.CreateRecordTypeParam{
		Name: name,
		Key:  key,
	}
	e, err := m.DB.CreateRecordType(ctx, p)
	if err != nil {
		return entity.RecordType{}, fmt.Errorf("failed to create record type: %w", err)
	}
	return e, nil
}

// CreateRecordTypes 議事録タイプを複数作成する。
func (m *ManageRecordType) CreateRecordTypes(
	ctx context.Context, ps []parameter.CreateRecordTypeParam,
) (int64, error) {
	es, err := m.DB.CreateRecordTypes(ctx, ps)
	if err != nil {
		return 0, fmt.Errorf("failed to create record types: %w", err)
	}
	return es, nil
}

// UpdateRecordType 議事録タイプを更新する。
func (m *ManageRecordType) UpdateRecordType(
	ctx context.Context, id uuid.UUID, name, key string,
) (entity.RecordType, error) {
	p := parameter.UpdateRecordTypeParams{
		Name: name,
		Key:  key,
	}
	e, err := m.DB.UpdateRecordType(ctx, id, p)
	if err != nil {
		return entity.RecordType{}, fmt.Errorf("failed to update record type: %w", err)
	}
	return e, nil
}

// DeleteRecordType 議事録タイプを削除する。
func (m *ManageRecordType) DeleteRecordType(ctx context.Context, id uuid.UUID) (int64, error) {
	c, err := m.DB.DeleteRecordType(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete record type: %w", err)
	}
	return c, nil
}

// PluralDeleteRecordTypes 議事録タイプを複数削除する。
func (m *ManageRecordType) PluralDeleteRecordTypes(ctx context.Context, ids []uuid.UUID) (int64, error) {
	c, err := m.DB.PluralDeleteRecordTypes(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete record types: %w", err)
	}
	return c, nil
}

// FindRecordTypeByID 議事録タイプをIDで取得する。
func (m *ManageRecordType) FindRecordTypeByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.RecordType, error) {
	e, err := m.DB.FindRecordTypeByID(ctx, id)
	if err != nil {
		return entity.RecordType{}, fmt.Errorf("failed to find record type by id: %w", err)
	}
	return e, nil
}

// FindRecordTypeByKey 議事録タイプをキーで取得する。
func (m *ManageRecordType) FindRecordTypeByKey(ctx context.Context, key string) (entity.RecordType, error) {
	e, err := m.DB.FindRecordTypeByKey(ctx, key)
	if err != nil {
		return entity.RecordType{}, fmt.Errorf("failed to find record type by key: %w", err)
	}
	return e, nil
}

// GetRecordTypes 議事録タイプを取得する。
func (m *ManageRecordType) GetRecordTypes(
	ctx context.Context,
	whereSearchName string,
	order parameter.RecordTypeOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.RecordType], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereRecordTypeParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset)},
			Limit:  entity.Int{Int64: int64(limit)},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit)},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetRecordTypes(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.RecordType]{}, fmt.Errorf("failed to get record types: %w", err)
	}
	return r, nil
}

// GetRecordTypesCount 議事録タイプの数を取得する。
func (m *ManageRecordType) GetRecordTypesCount(
	ctx context.Context,
	whereSearchName string,
) (int64, error) {
	p := parameter.WhereRecordTypeParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	c, err := m.DB.CountRecordTypes(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get record types count: %w", err)
	}
	return c, nil
}
