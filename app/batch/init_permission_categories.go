package batch

import (
	"context"
	"fmt"

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
