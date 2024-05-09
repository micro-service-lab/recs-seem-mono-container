package pgadapter_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/testutils/factory"
)

func TestPgAdapter_MimeType(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	adapter := NewDummyPgAdapter(t)

	t.Cleanup(func() {
		err := adapter.Cleanup(ctx)
		require.NoError(t, err)
	})

	series := []SeriesTest{
		{
			Name: "create one",
			Test: func(t *testing.T) {
				sd, err := adapter.Begin(ctx)
				assert.NoError(t, err)
				t.Cleanup(func() {
					err := adapter.Rollback(ctx, sd)
					require.NoError(t, err)
				})
				fp, err := factory.Generator.NewMimeTypes(1)
				assert.NoError(t, err)
				p := fp.ForCreateParam()[0]
				e, err := adapter.CreateMimeTypeWithSd(ctx, sd, p)
				assert.NoError(t, err)
				assert.NotEmpty(t, e.MimeTypeID)
				assert.Equal(t, p.Key, e.Key)
				assert.Equal(t, p.Name, e.Name)
				e, err = adapter.FindMimeTypeByIDWithSd(ctx, sd, e.MimeTypeID)
				assert.NoError(t, err)
				assert.NotEmpty(t, e.MimeTypeID)
				assert.Equal(t, p.Key, e.Key)
				assert.Equal(t, p.Name, e.Name)
				e, err = adapter.FindMimeTypeByKeyWithSd(ctx, sd, p.Key)
				assert.NoError(t, err)
				assert.NotEmpty(t, e.MimeTypeID)
				assert.Equal(t, p.Key, e.Key)
				assert.Equal(t, p.Name, e.Name)
				count, err := adapter.CountMimeTypesWithSd(ctx, sd, parameter.WhereMimeTypeParam{})
				assert.NoError(t, err)
				assert.Equal(t, int64(1), count)
				where1 := parameter.WhereMimeTypeParam{
					WhereLikeName: true,
					SearchName:    fp[0].Name,
				}
				where2 := parameter.WhereMimeTypeParam{
					WhereLikeName: true,
					SearchName:    "name",
				}
				count, err = adapter.CountMimeTypesWithSd(ctx, sd, where1)
				assert.NoError(t, err)
				assert.Equal(t, fp.CountContainsName(fp[0].Name), count)
				count, err = adapter.CountMimeTypesWithSd(ctx, sd, where2)
				assert.NoError(t, err)
				assert.Equal(t, fp.CountContainsName("name"), count)
				el, err := adapter.GetMimeTypesWithSd(
					ctx,
					sd,
					where1,
					parameter.MimeTypeOrderMethodDefault,
					validNp,
					invalidCp,
					validWc,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, 1)
				assert.Equal(t, p.Key, el.Data[0].Key)
				assert.Equal(t, el.WithCount.Count, int64(1))
			},
		},
		{
			Name: "create plural",
			Test: func(t *testing.T) {
				sd, err := adapter.Begin(ctx)
				assert.NoError(t, err)
				t.Cleanup(func() {
					err := adapter.Rollback(ctx, sd)
					require.NoError(t, err)
				})
				fp, err := factory.Generator.NewMimeTypes(3)
				assert.NoError(t, err)
				ps := fp.ForCreateParam()
				count, err := adapter.CreateMimeTypesWithSd(ctx, sd, ps)
				assert.NoError(t, err)
				assert.Equal(t, int64(3), count)
				count, err = adapter.CountMimeTypesWithSd(ctx, sd, parameter.WhereMimeTypeParam{})
				assert.NoError(t, err)
				assert.Equal(t, int64(3), count)
				where := parameter.WhereMimeTypeParam{
					WhereLikeName: true,
					SearchName:    fp[0].Name,
				}
				where2 := parameter.WhereMimeTypeParam{}
				el, err := adapter.GetMimeTypesWithSd(
					ctx,
					sd,
					where,
					parameter.MimeTypeOrderMethodDefault,
					validNp,
					invalidCp,
					validWc,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, int(fp.CountContainsName(fp[0].Name)))
				assert.Equal(t, el.WithCount.Count, fp.CountContainsName(fp[0].Name))
				el, err = adapter.GetMimeTypesWithSd(
					ctx,
					sd,
					where2,
					parameter.MimeTypeOrderMethodDefault,
					validNp,
					invalidCp,
					validWc,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, 3)
				assert.Equal(t, el.WithCount.Count, int64(3))
				on, err := adapter.GetMimeTypesWithSd(
					ctx,
					sd,
					where2,
					parameter.MimeTypeOrderMethodName,
					validNp,
					invalidCp,
					validWc,
				)
				assert.NoError(t, err)
				assert.Len(t, on.Data, 3)
				assert.Equal(t, on.WithCount.Count, int64(3))
				assert.Equal(t, on.Data[0].Name, fp.OrderByNames()[0].Name)
				assert.Equal(t, on.Data[1].Name, fp.OrderByNames()[1].Name)
				assert.Equal(t, on.Data[2].Name, fp.OrderByNames()[2].Name)
				// getPlural
				el, err = adapter.GetPluralMimeTypesWithSd(
					ctx,
					sd,
					[]uuid.UUID{el.Data[0].MimeTypeID, el.Data[1].MimeTypeID},
					validNp,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, 2)
				// delete
				err = adapter.DeleteMimeTypeWithSd(ctx, sd, el.Data[0].MimeTypeID)
				assert.NoError(t, err)
				count, err = adapter.CountMimeTypesWithSd(ctx, sd, parameter.WhereMimeTypeParam{})
				assert.NoError(t, err)
				assert.Equal(t, int64(2), count)
				el, err = adapter.GetMimeTypesWithSd(
					ctx,
					sd,
					where2,
					parameter.MimeTypeOrderMethodDefault,
					validNp,
					invalidCp,
					validWc,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, 2)
				assert.Equal(t, el.WithCount.Count, int64(2))
				// update
				p := parameter.UpdateMimeTypeParams{
					Name: "name4",
					Key:  "key4",
				}
				e, err := adapter.UpdateMimeTypeWithSd(ctx, sd, el.Data[0].MimeTypeID, p)
				assert.NoError(t, err)
				assert.Equal(t, p.Name, e.Name)
				assert.Equal(t, p.Key, e.Key)
				e, err = adapter.FindMimeTypeByIDWithSd(ctx, sd, el.Data[0].MimeTypeID)
				assert.NoError(t, err)
				assert.Equal(t, p.Name, e.Name)
				// update by key
				p2 := parameter.UpdateMimeTypeByKeyParams{
					Name: "name5",
				}
				e, err = adapter.UpdateMimeTypeByKeyWithSd(ctx, sd, el.Data[1].Key, p2)
				assert.NoError(t, err)
				assert.Equal(t, p2.Name, e.Name)
				e, err = adapter.FindMimeTypeByKeyWithSd(ctx, sd, el.Data[1].Key)
				assert.NoError(t, err)
				assert.Equal(t, p2.Name, e.Name)
			},
		},
	}

	for _, s := range series {
		ss := s
		t.Run(s.Name, func(t *testing.T) {
			t.Parallel()
			ss.Test(t)
		})
	}
}
