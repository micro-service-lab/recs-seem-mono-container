package batch

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
)

// InitPermissions is a batch to initialize permission.
type InitPermissions struct {
	Manager *service.ManagerInterface
}

// Run initializes permission.
func (c *InitPermissions) Run(ctx context.Context) error {
	searchPC, _, err := searchPermissionCategoryID(ctx, c.Manager)
	if err != nil {
		return fmt.Errorf("failed to search permission category: %w", err)
	}
	var as []parameter.CreatePermissionParam
	for _, a := range service.Permissions {
		cID := searchPC(a.PermissionCategoryKey)
		if cID == uuid.Nil {
			return fmt.Errorf("failed to find permission category: %s", a.PermissionCategoryKey)
		}
		as = append(as, parameter.CreatePermissionParam{
			Name:                 a.Name,
			Key:                  a.Key,
			Description:          a.Description,
			PermissionCategoryID: cID,
		})
	}
	_, err = (*c.Manager).CreatePermissions(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}
	return nil
}

// RunDiff run only if there is a difference.
func (c *InitPermissions) RunDiff(ctx context.Context, notDel, deepEqual bool) error {
	searchPcID, searchPcKey, err := searchPermissionCategoryID(ctx, c.Manager)
	if err != nil {
		return fmt.Errorf("failed to search permission category: %w", err)
	}
	exists, err := (*c.Manager).GetPermissions(
		ctx,
		"",
		[]uuid.UUID{},
		parameter.PermissionOrderMethodDefault,
		parameter.NonePagination,
		parameter.Limit(0),
		parameter.Cursor(""),
		parameter.Offset(0),
		parameter.WithCount(false),
	)
	if err != nil {
		return fmt.Errorf("failed to get permission: %w", err)
	}
	existData := make(map[uuid.UUID]service.Permission, len(exists.Data))
	existIDs := make([]uuid.UUID, len(exists.Data))
	existKey := make([]string, len(exists.Data))
	for i, a := range exists.Data {
		key := searchPcKey(a.PermissionCategoryID)
		if key == "" {
			return fmt.Errorf("failed to find permission category: %s", a.PermissionCategoryID)
		}
		existData[a.PermissionID] = service.Permission{
			Name:                  a.Name,
			Key:                   a.Key,
			Description:           a.Description,
			PermissionCategoryKey: key,
		}
		existIDs[i] = a.PermissionID
		existKey[i] = a.Key
	}
	var as []parameter.CreatePermissionParam
	for _, a := range service.Permissions {
		matchIndex := contains(existKey, a.Key)
		pcID := searchPcID(a.PermissionCategoryKey)
		if pcID == uuid.Nil {
			return fmt.Errorf("failed to find permission category: %s", a.PermissionCategoryKey)
		}
		if matchIndex == -1 {
			as = append(as, parameter.CreatePermissionParam{
				Name:                 a.Name,
				Key:                  a.Key,
				Description:          a.Description,
				PermissionCategoryID: pcID,
			})
		} else {
			var uid uuid.UUID
			existIDs, uid = removeUUID(existIDs, matchIndex)
			existKey, _ = removeString(existKey, matchIndex)
			if deepEqual {
				de := isDeepEqual(a, existData[uid])
				if !de {
					_, err = (*c.Manager).UpdatePermission(
						ctx,
						uid,
						a.Name,
						a.Key,
						a.Description,
						pcID,
					)
					if err != nil {
						return fmt.Errorf("failed to update permission category: %w", err)
					}
				}
			}
			delete(existData, uid)
		}
	}
	_, err = (*c.Manager).CreatePermissions(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}
	if len(existIDs) > 0 && !notDel {
		err = (*c.Manager).PluralDeletePermissions(ctx, existIDs)
		if err != nil {
			return fmt.Errorf("failed to delete permission: %w", err)
		}
	}
	return nil
}

func searchPermissionCategoryID(
	ctx context.Context,
	c *service.ManagerInterface,
) (func(key service.PermissionCategoryKey) uuid.UUID,
	func(id uuid.UUID) service.PermissionCategoryKey,
	error,
) {
	pc, err := (*c).GetPermissionCategories(
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
		return nil, nil, fmt.Errorf("failed to get permission category: %w", err)
	}
	searchPcID := func(key service.PermissionCategoryKey) uuid.UUID {
		for _, a := range pc.Data {
			if a.Key == string(key) {
				return a.PermissionCategoryID
			}
		}
		return uuid.Nil
	}
	searchPCKey := func(id uuid.UUID) service.PermissionCategoryKey {
		for _, a := range pc.Data {
			if a.PermissionCategoryID == id {
				return service.PermissionCategoryKey(a.Key)
			}
		}
		return ""
	}
	return searchPcID, searchPCKey, nil
}
