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

func TestPgAdapter_AttendanceType(t *testing.T) {
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
				fp, err := factory.Generator.NewAttendanceTypes(1)
				assert.NoError(t, err)
				p := fp.ForCreateParam()[0]
				e, err := adapter.CreateAttendanceTypeWithSd(ctx, sd, p)
				assert.NoError(t, err)
				assert.NotEmpty(t, e.AttendanceTypeID)
				assert.Equal(t, p.Key, e.Key)
				assert.Equal(t, p.Name, e.Name)
				e, err = adapter.FindAttendanceTypeByIDWithSd(ctx, sd, e.AttendanceTypeID)
				assert.NoError(t, err)
				assert.NotEmpty(t, e.AttendanceTypeID)
				assert.Equal(t, p.Key, e.Key)
				assert.Equal(t, p.Name, e.Name)
				e, err = adapter.FindAttendanceTypeByKeyWithSd(ctx, sd, p.Key)
				assert.NoError(t, err)
				assert.NotEmpty(t, e.AttendanceTypeID)
				assert.Equal(t, p.Key, e.Key)
				assert.Equal(t, p.Name, e.Name)
				count, err := adapter.CountAttendanceTypesWithSd(ctx, sd, parameter.WhereAttendanceTypeParam{})
				assert.NoError(t, err)
				assert.Equal(t, int64(1), count)
				where1 := parameter.WhereAttendanceTypeParam{
					WhereLikeName: true,
					SearchName:    fp[0].Name,
				}
				where2 := parameter.WhereAttendanceTypeParam{
					WhereLikeName: true,
					SearchName:    "name",
				}
				count, err = adapter.CountAttendanceTypesWithSd(ctx, sd, where1)
				assert.NoError(t, err)
				assert.Equal(t, fp.CountContainsName(fp[0].Name), count)
				count, err = adapter.CountAttendanceTypesWithSd(ctx, sd, where2)
				assert.NoError(t, err)
				assert.Equal(t, fp.CountContainsName("name"), count)
				el, err := adapter.GetAttendanceTypesWithSd(
					ctx,
					sd,
					where1,
					parameter.AttendanceTypeOrderMethodDefault,
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
				fp, err := factory.Generator.NewAttendanceTypes(3)
				assert.NoError(t, err)
				ps := fp.ForCreateParam()
				count, err := adapter.CreateAttendanceTypesWithSd(ctx, sd, ps)
				assert.NoError(t, err)
				assert.Equal(t, int64(3), count)
				count, err = adapter.CountAttendanceTypesWithSd(ctx, sd, parameter.WhereAttendanceTypeParam{})
				assert.NoError(t, err)
				assert.Equal(t, int64(3), count)
				where := parameter.WhereAttendanceTypeParam{
					WhereLikeName: true,
					SearchName:    fp[0].Name,
				}
				where2 := parameter.WhereAttendanceTypeParam{}
				el, err := adapter.GetAttendanceTypesWithSd(
					ctx,
					sd,
					where,
					parameter.AttendanceTypeOrderMethodDefault,
					validNp,
					invalidCp,
					validWc,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, int(fp.CountContainsName(fp[0].Name)))
				assert.Equal(t, el.WithCount.Count, fp.CountContainsName(fp[0].Name))
				el, err = adapter.GetAttendanceTypesWithSd(
					ctx,
					sd,
					where2,
					parameter.AttendanceTypeOrderMethodDefault,
					validNp,
					invalidCp,
					validWc,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, 3)
				assert.Equal(t, el.WithCount.Count, int64(3))
				on, err := adapter.GetAttendanceTypesWithSd(
					ctx,
					sd,
					where2,
					parameter.AttendanceTypeOrderMethodName,
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
				el, err = adapter.GetPluralAttendanceTypesWithSd(
					ctx,
					sd,
					[]uuid.UUID{el.Data[0].AttendanceTypeID, el.Data[1].AttendanceTypeID},
					parameter.AttendanceTypeOrderMethodDefault,
					validNp,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, 2)
				// delete
				_, err = adapter.DeleteAttendanceTypeWithSd(ctx, sd, el.Data[0].AttendanceTypeID)
				assert.NoError(t, err)
				count, err = adapter.CountAttendanceTypesWithSd(ctx, sd, parameter.WhereAttendanceTypeParam{})
				assert.NoError(t, err)
				assert.Equal(t, int64(2), count)
				el, err = adapter.GetAttendanceTypesWithSd(
					ctx,
					sd,
					where2,
					parameter.AttendanceTypeOrderMethodDefault,
					validNp,
					invalidCp,
					validWc,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, 2)
				assert.Equal(t, el.WithCount.Count, int64(2))
				// update
				p := parameter.UpdateAttendanceTypeParams{
					Name:  "name4",
					Key:   "key4",
					Color: "#000000",
				}
				e, err := adapter.UpdateAttendanceTypeWithSd(ctx, sd, el.Data[0].AttendanceTypeID, p)
				assert.NoError(t, err)
				assert.Equal(t, p.Name, e.Name)
				assert.Equal(t, p.Key, e.Key)
				assert.Equal(t, p.Color, e.Color)
				e, err = adapter.FindAttendanceTypeByIDWithSd(ctx, sd, el.Data[0].AttendanceTypeID)
				assert.NoError(t, err)
				assert.Equal(t, p.Name, e.Name)
				assert.Equal(t, p.Key, e.Key)
				assert.Equal(t, p.Color, e.Color)
				// update by key
				p2 := parameter.UpdateAttendanceTypeByKeyParams{
					Name:  "name5",
					Color: "#000001",
				}
				e, err = adapter.UpdateAttendanceTypeByKeyWithSd(ctx, sd, el.Data[1].Key, p2)
				assert.NoError(t, err)
				assert.Equal(t, p2.Name, e.Name)
				assert.Equal(t, p2.Color, e.Color)
				e, err = adapter.FindAttendanceTypeByKeyWithSd(ctx, sd, el.Data[1].Key)
				assert.NoError(t, err)
				assert.Equal(t, p2.Name, e.Name)
				assert.Equal(t, p2.Color, e.Color)
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
