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

func TestPgAdapter_RecordType(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
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
				fp, err := factory.Generator.NewRecordTypes(1)
				assert.NoError(t, err)
				p := fp.ForCreateParam()[0]
				e, err := adapter.CreateRecordTypeWithSd(ctx, sd, p)
				assert.NoError(t, err)
				assert.NotEmpty(t, e.RecordTypeID)
				assert.Equal(t, p.Key, e.Key)
				assert.Equal(t, p.Name, e.Name)
				e, err = adapter.FindRecordTypeByIDWithSd(ctx, sd, e.RecordTypeID)
				assert.NoError(t, err)
				assert.NotEmpty(t, e.RecordTypeID)
				assert.Equal(t, p.Key, e.Key)
				assert.Equal(t, p.Name, e.Name)
				e, err = adapter.FindRecordTypeByKeyWithSd(ctx, sd, p.Key)
				assert.NoError(t, err)
				assert.NotEmpty(t, e.RecordTypeID)
				assert.Equal(t, p.Key, e.Key)
				assert.Equal(t, p.Name, e.Name)
				count, err := adapter.CountRecordTypesWithSd(ctx, sd, parameter.WhereRecordTypeParam{})
				assert.NoError(t, err)
				assert.Equal(t, int64(1), count)
				where1 := parameter.WhereRecordTypeParam{
					WhereLikeName: true,
					SearchName:    fp[0].Name,
				}
				where2 := parameter.WhereRecordTypeParam{
					WhereLikeName: true,
					SearchName:    "name",
				}
				count, err = adapter.CountRecordTypesWithSd(ctx, sd, where1)
				assert.NoError(t, err)
				assert.Equal(t, fp.CountContainsName(fp[0].Name), count)
				count, err = adapter.CountRecordTypesWithSd(ctx, sd, where2)
				assert.NoError(t, err)
				assert.Equal(t, fp.CountContainsName("name"), count)
				el, err := adapter.GetRecordTypesWithSd(
					ctx,
					sd,
					where1,
					parameter.RecordTypeOrderMethodDefault,
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
				fp, err := factory.Generator.NewRecordTypes(3)
				assert.NoError(t, err)
				ps := fp.ForCreateParam()
				count, err := adapter.CreateRecordTypesWithSd(ctx, sd, ps)
				assert.NoError(t, err)
				assert.Equal(t, int64(3), count)
				count, err = adapter.CountRecordTypesWithSd(ctx, sd, parameter.WhereRecordTypeParam{})
				assert.NoError(t, err)
				assert.Equal(t, int64(3), count)
				where := parameter.WhereRecordTypeParam{
					WhereLikeName: true,
					SearchName:    fp[0].Name,
				}
				where2 := parameter.WhereRecordTypeParam{}
				el, err := adapter.GetRecordTypesWithSd(
					ctx,
					sd,
					where,
					parameter.RecordTypeOrderMethodDefault,
					validNp,
					invalidCp,
					validWc,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, int(fp.CountContainsName(fp[0].Name)))
				assert.Equal(t, el.WithCount.Count, fp.CountContainsName(fp[0].Name))
				el, err = adapter.GetRecordTypesWithSd(
					ctx,
					sd,
					where2,
					parameter.RecordTypeOrderMethodDefault,
					validNp,
					invalidCp,
					validWc,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, 3)
				assert.Equal(t, el.WithCount.Count, int64(3))
				on, err := adapter.GetRecordTypesWithSd(
					ctx,
					sd,
					where2,
					parameter.RecordTypeOrderMethodName,
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
				el, err = adapter.GetPluralRecordTypesWithSd(
					ctx,
					sd,
					[]uuid.UUID{el.Data[0].RecordTypeID, el.Data[1].RecordTypeID},
					validNp,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, 2)
				// delete
				_, err = adapter.DeleteRecordTypeWithSd(ctx, sd, el.Data[0].RecordTypeID)
				assert.NoError(t, err)
				count, err = adapter.CountRecordTypesWithSd(ctx, sd, parameter.WhereRecordTypeParam{})
				assert.NoError(t, err)
				assert.Equal(t, int64(2), count)
				el, err = adapter.GetRecordTypesWithSd(
					ctx,
					sd,
					where2,
					parameter.RecordTypeOrderMethodDefault,
					validNp,
					invalidCp,
					validWc,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, 2)
				assert.Equal(t, el.WithCount.Count, int64(2))
				// update
				p := parameter.UpdateRecordTypeParams{
					Name: "name4",
					Key:  "key4",
				}
				e, err := adapter.UpdateRecordTypeWithSd(ctx, sd, el.Data[0].RecordTypeID, p)
				assert.NoError(t, err)
				assert.Equal(t, p.Name, e.Name)
				assert.Equal(t, p.Key, e.Key)
				e, err = adapter.FindRecordTypeByIDWithSd(ctx, sd, el.Data[0].RecordTypeID)
				assert.NoError(t, err)
				assert.Equal(t, p.Name, e.Name)
				assert.Equal(t, p.Key, e.Key)
				// update by key
				p2 := parameter.UpdateRecordTypeByKeyParams{
					Name: "name5",
				}
				e, err = adapter.UpdateRecordTypeByKeyWithSd(ctx, sd, el.Data[1].Key, p2)
				assert.NoError(t, err)
				assert.Equal(t, p2.Name, e.Name)
				e, err = adapter.FindRecordTypeByKeyWithSd(ctx, sd, el.Data[1].Key)
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
