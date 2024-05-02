package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

func TestManageAttendStatus_CreateAttendStatus(t *testing.T) {
	storeMock := &store.StoreMock{
		CreateAttendStatusFunc: func(
			_ context.Context, p parameter.CreateAttendStatusParam,
		) (entity.AttendStatus, error) {
			return entity.AttendStatus{
				AttendStatusID: uuid.New(),
				Name:           p.Name,
				Key:            p.Key,
			}, nil
		},
	}
	s := service.ManageAttendStatus{
		DB: storeMock,
	}
	ctx := context.Background()
	name := "name"
	key := "key"
	attendStatus, err := s.CreateAttendStatus(ctx, name, key)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, attendStatus.AttendStatusID)
	assert.Equal(t, name, attendStatus.Name)
	assert.Equal(t, key, attendStatus.Key)
}

func TestManageAttendStatus_CreateAttendStatuses(t *testing.T) {
	storeMock := &store.StoreMock{
		CreateAttendStatusesFunc: func(
			_ context.Context, ps []parameter.CreateAttendStatusParam,
		) (int64, error) {
			return int64(len(ps)), nil
		},
	}
	s := service.ManageAttendStatus{
		DB: storeMock,
	}
	ctx := context.Background()
	ps := []parameter.CreateAttendStatusParam{
		{Name: "name1", Key: "key1"},
		{Name: "name2", Key: "key2"},
	}
	n, err := s.CreateAttendStatuses(ctx, ps)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(ps)), n)
}

func TestManageAttendStatus_UpdateAttendStatus(t *testing.T) {
	storeMock := &store.StoreMock{
		UpdateAttendStatusFunc: func(
			_ context.Context, id uuid.UUID, p parameter.UpdateAttendStatusParams,
		) (entity.AttendStatus, error) {
			return entity.AttendStatus{
				AttendStatusID: id,
				Name:           p.Name,
				Key:            p.Key,
			}, nil
		},
	}
	s := service.ManageAttendStatus{
		DB: storeMock,
	}
	ctx := context.Background()
	id := uuid.New()
	name := "update name"
	key := "update key"
	attendStatus, err := s.UpdateAttendStatus(ctx, id, name, key)
	assert.NoError(t, err)
	assert.Equal(t, id, attendStatus.AttendStatusID)
	assert.Equal(t, name, attendStatus.Name)
	assert.Equal(t, key, attendStatus.Key)
}

func TestManageAttendStatus_DeleteAttendStatus(t *testing.T) {
	storeMock := &store.StoreMock{
		DeleteAttendStatusFunc: func(_ context.Context, _ uuid.UUID) error {
			return nil
		},
	}
	s := service.ManageAttendStatus{
		DB: storeMock,
	}
	ctx := context.Background()
	id := uuid.New()
	err := s.DeleteAttendStatus(ctx, id)
	assert.NoError(t, err)
}

func TestManageAttendStatus_FindAttendStatusByID(t *testing.T) {
	storeMock := &store.StoreMock{
		FindAttendStatusByIDFunc: func(_ context.Context, id uuid.UUID) (entity.AttendStatus, error) {
			return entity.AttendStatus{
				AttendStatusID: id,
				Name:           "name",
				Key:            "key",
			}, nil
		},
	}
	s := service.ManageAttendStatus{
		DB: storeMock,
	}
	ctx := context.Background()
	id := uuid.New()
	attendStatus, err := s.FindAttendStatusByID(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, id, attendStatus.AttendStatusID)
	assert.Equal(t, "name", attendStatus.Name)
	assert.Equal(t, "key", attendStatus.Key)
}

func TestManageAttendStatus_GetAttendStatuses(t *testing.T) {
	data := []entity.AttendStatus{
		{
			AttendStatusID: uuid.New(),
			Name:           "name1",
			Key:            "key1",
		},
		{
			AttendStatusID: uuid.New(),
			Name:           "name2",
			Key:            "key2",
		},
	}
	storeMock := &store.StoreMock{
		GetAttendStatusesFunc: func(
			_ context.Context,
			_ parameter.WhereAttendStatusParam,
			_ parameter.AttendStatusOrderMethod,
			_ store.NumberedPaginationParam,
			_ store.CursorPaginationParam,
			_ store.WithCountParam,
		) (store.ListResult[entity.AttendStatus], error) {
			return store.ListResult[entity.AttendStatus]{
				Data: data,
				WithCount: store.WithCountAttribute{
					Count: int64(len(data)),
				},
			}, nil
		},
	}
	s := service.ManageAttendStatus{
		DB: storeMock,
	}
	ctx := context.Background()
	attendStatuses, err := s.GetAttendStatuses(
		ctx, "", parameter.AttendStatusOrderMethodDefault, parameter.NonePagination, 0, "", 0, true)
	assert.NoError(t, err)
	assert.Len(t, attendStatuses.Data, 2)
	assert.NotEqual(t, uuid.Nil, attendStatuses.Data[0].AttendStatusID)
	assert.Equal(t, data[0].Name, attendStatuses.Data[0].Name)
	assert.Equal(t, data[0].Key, attendStatuses.Data[0].Key)
	assert.NotEqual(t, uuid.Nil, attendStatuses.Data[1].AttendStatusID)
	assert.Equal(t, data[1].Name, attendStatuses.Data[1].Name)
	assert.Equal(t, data[1].Key, attendStatuses.Data[1].Key)
	assert.Equal(t, int64(2), attendStatuses.WithCount.Count)
}

func TestManageAttendStatus_GetAttendStatusesCount(t *testing.T) {
	s := service.ManageAttendStatus{
		DB: &store.StoreMock{
			CountAttendStatusesFunc: func(_ context.Context, _ parameter.WhereAttendStatusParam) (int64, error) {
				return 2, nil
			},
		},
	}
	ctx := context.Background()
	count, err := s.GetAttendStatusesCount(ctx, "")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}
