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

func TestPgAdapter_Permission(t *testing.T) {
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
				fp, err := factory.Generator.NewPermissions(1)
				assert.NoError(t, err)
				fp, err = fp.WithPermissionCategory(
					createPermissionCategories(ctx, t, sd, adapter, 3),
				)
				assert.NoError(t, err)
				p := fp.ForCreateParam()[0]
				e, err := adapter.CreatePermissionWithSd(ctx, sd, p)
				assert.NoError(t, err)
				assert.NotEmpty(t, e.PermissionID)
				assert.Equal(t, p.Key, e.Key)
				assert.Equal(t, p.Name, e.Name)
				e, err = adapter.FindPermissionByIDWithSd(ctx, sd, e.PermissionID)
				assert.NoError(t, err)
				assert.NotEmpty(t, e.PermissionID)
				assert.Equal(t, p.Key, e.Key)
				assert.Equal(t, p.Name, e.Name)
				e, err = adapter.FindPermissionByKeyWithSd(ctx, sd, p.Key)
				assert.NoError(t, err)
				assert.NotEmpty(t, e.PermissionID)
				assert.Equal(t, p.Key, e.Key)
				assert.Equal(t, p.Name, e.Name)
				count, err := adapter.CountPermissionsWithSd(ctx, sd, parameter.WherePermissionParam{})
				assert.NoError(t, err)
				assert.Equal(t, int64(1), count)
				where1 := parameter.WherePermissionParam{
					WhereLikeName: true,
					SearchName:    fp[0].Name,
				}
				where2 := parameter.WherePermissionParam{
					WhereLikeName: true,
					SearchName:    "name",
				}
				count, err = adapter.CountPermissionsWithSd(ctx, sd, where1)
				assert.NoError(t, err)
				assert.Equal(t, fp.CountContainsName(fp[0].Name), count)
				count, err = adapter.CountPermissionsWithSd(ctx, sd, where2)
				assert.NoError(t, err)
				assert.Equal(t, fp.CountContainsName("name"), count)
				el, err := adapter.GetPermissionsWithSd(
					ctx,
					sd,
					where1,
					parameter.PermissionOrderMethodDefault,
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
				fp, err := factory.Generator.NewPermissions(3)
				assert.NoError(t, err)
				fp, err = fp.WithPermissionCategory(
					createPermissionCategories(ctx, t, sd, adapter, 5),
				)
				assert.NoError(t, err)
				ps := fp.ForCreateParam()
				count, err := adapter.CreatePermissionsWithSd(ctx, sd, ps)
				assert.NoError(t, err)
				assert.Equal(t, int64(3), count)
				count, err = adapter.CountPermissionsWithSd(ctx, sd, parameter.WherePermissionParam{})
				assert.NoError(t, err)
				assert.Equal(t, int64(3), count)
				where := parameter.WherePermissionParam{
					WhereLikeName: true,
					SearchName:    fp[0].Name,
				}
				where2 := parameter.WherePermissionParam{}
				el, err := adapter.GetPermissionsWithSd(
					ctx,
					sd,
					where,
					parameter.PermissionOrderMethodDefault,
					validNp,
					invalidCp,
					validWc,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, int(fp.CountContainsName(fp[0].Name)))
				assert.Equal(t, el.WithCount.Count, fp.CountContainsName(fp[0].Name))
				el, err = adapter.GetPermissionsWithSd(
					ctx,
					sd,
					where2,
					parameter.PermissionOrderMethodDefault,
					validNp,
					invalidCp,
					validWc,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, 3)
				assert.Equal(t, el.WithCount.Count, int64(3))
				on, err := adapter.GetPermissionsWithSd(
					ctx,
					sd,
					where2,
					parameter.PermissionOrderMethodName,
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
				el, err = adapter.GetPluralPermissionsWithSd(
					ctx,
					sd,
					[]uuid.UUID{el.Data[0].PermissionID, el.Data[1].PermissionID},
					validNp,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, 2)
				// delete
				err = adapter.DeletePermissionWithSd(ctx, sd, el.Data[0].PermissionID)
				assert.NoError(t, err)
				count, err = adapter.CountPermissionsWithSd(ctx, sd, parameter.WherePermissionParam{})
				assert.NoError(t, err)
				assert.Equal(t, int64(2), count)
				el, err = adapter.GetPermissionsWithSd(
					ctx,
					sd,
					where2,
					parameter.PermissionOrderMethodDefault,
					validNp,
					invalidCp,
					validWc,
				)
				assert.NoError(t, err)
				assert.Len(t, el.Data, 2)
				assert.Equal(t, el.WithCount.Count, int64(2))
				// update
				p := parameter.UpdatePermissionParams{
					Name:                 "name4",
					Key:                  "key4",
					Description:          "description4",
					PermissionCategoryID: el.Data[0].PermissionCategoryID,
				}
				e, err := adapter.UpdatePermissionWithSd(ctx, sd, el.Data[0].PermissionID, p)
				assert.NoError(t, err)
				assert.Equal(t, p.Name, e.Name)
				assert.Equal(t, p.Key, e.Key)
				assert.Equal(t, p.Description, e.Description)
				assert.Equal(t, p.PermissionCategoryID, e.PermissionCategoryID)
				e, err = adapter.FindPermissionByIDWithSd(ctx, sd, el.Data[0].PermissionID)
				assert.NoError(t, err)
				assert.Equal(t, p.Name, e.Name)
				// update by key
				p2 := parameter.UpdatePermissionByKeyParams{
					Name:                 "name5",
					Description:          "description5",
					PermissionCategoryID: el.Data[1].PermissionCategoryID,
				}
				e, err = adapter.UpdatePermissionByKeyWithSd(ctx, sd, el.Data[1].Key, p2)
				assert.NoError(t, err)
				assert.Equal(t, p2.Name, e.Name)
				assert.Equal(t, p2.Description, e.Description)
				assert.Equal(t, p2.PermissionCategoryID, e.PermissionCategoryID)
				e, err = adapter.FindPermissionByKeyWithSd(ctx, sd, el.Data[1].Key)
				assert.NoError(t, err)
				assert.Equal(t, p2.Name, e.Name)
				assert.Equal(t, p2.Description, e.Description)
				assert.Equal(t, p2.PermissionCategoryID, e.PermissionCategoryID)
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
