package batch

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
)

// InitPermissionCategories is a batch to initialize permission categories.
type InitPermissionCategories struct {
	Manager *service.ManagerInterface
}

// Run initializes permission categories.
func (c *InitPermissionCategories) Run(ctx context.Context) error {
	var as []parameter.CreatePermissionCategoryParam
	for _, a := range service.PermissionCategories {
		as = append(as, parameter.CreatePermissionCategoryParam{
			Name:        a.Name,
			Key:         a.Key,
			Description: a.Description,
		})
	}
	_, err := (*c.Manager).CreatePermissionCategories(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create permission categories: %w", err)
	}
	return nil
}

// RunDiff run only if there is a difference.
func (c *InitPermissionCategories) RunDiff(ctx context.Context, notDel, deepEqual bool) error {
	exists, err := (*c.Manager).GetPermissionCategories(
		ctx,
		"",
		parameter.PermissionCategoryOrderMethodDefault,
		parameter.NonePagination,
		parameter.Limit(0),
		parameter.Cursor(""),
		parameter.Offset(0),
		parameter.WithCount(false),
	)
	if err != nil {
		return fmt.Errorf("failed to get permission categories: %w", err)
	}
	existData := make(map[uuid.UUID]service.PermissionCategory, len(exists.Data))
	existIDs := make([]uuid.UUID, len(exists.Data))
	existKey := make([]string, len(exists.Data))
	for i, a := range exists.Data {
		existData[a.PermissionCategoryID] = service.PermissionCategory{
			Name:        a.Name,
			Key:         a.Key,
			Description: a.Description,
		}
		existIDs[i] = a.PermissionCategoryID
		existKey[i] = a.Key
	}
	var as []parameter.CreatePermissionCategoryParam
	for _, a := range service.PermissionCategories {
		matchIndex := contains(existKey, a.Key)
		if matchIndex == -1 {
			as = append(as, parameter.CreatePermissionCategoryParam{
				Name:        a.Name,
				Key:         a.Key,
				Description: a.Description,
			})
		} else {
			var uid uuid.UUID
			existIDs, uid = removeUUID(existIDs, matchIndex)
			existKey, _ = removeString(existKey, matchIndex)
			if deepEqual {
				de := isDeepEqual(a, existData[uid])
				if !de {
					_, err = (*c.Manager).UpdatePermissionCategory(
						ctx,
						uid,
						a.Name,
						a.Key,
						a.Description,
					)
					if err != nil {
						return fmt.Errorf("failed to update permission category: %w", err)
					}
				}
			}
			delete(existData, uid)
		}
	}
	_, err = (*c.Manager).CreatePermissionCategories(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create permission categories: %w", err)
	}
	if len(existIDs) > 0 && !notDel {
		err = (*c.Manager).PluralDeletePermissionCategories(ctx, existIDs)
		if err != nil {
			return fmt.Errorf("failed to delete permission categories: %w", err)
		}
	}
	return nil
}
