package batch

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
)

// InitPolicyCategories is a batch to initialize policy categories.
type InitPolicyCategories struct {
	Manager *service.ManagerInterface
}

// Run initializes policy categories.
func (c *InitPolicyCategories) Run(ctx context.Context) error {
	var as []parameter.CreatePolicyCategoryParam
	for _, a := range service.PolicyCategories {
		as = append(as, parameter.CreatePolicyCategoryParam{
			Name:        a.Name,
			Key:         a.Key,
			Description: a.Description,
		})
	}
	_, err := (*c.Manager).CreatePolicyCategories(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create policy categories: %w", err)
	}
	return nil
}

// RunDiff run only if there is a difference.
func (c *InitPolicyCategories) RunDiff(ctx context.Context, notDel, deepEqual bool) error {
	exists, err := (*c.Manager).GetPolicyCategories(
		ctx,
		"",
		parameter.PolicyCategoryOrderMethodDefault,
		parameter.NonePagination,
		parameter.Limit(0),
		parameter.Cursor(""),
		parameter.Offset(0),
		parameter.WithCount(false),
	)
	if err != nil {
		return fmt.Errorf("failed to get policy categories: %w", err)
	}
	existData := make(map[uuid.UUID]service.PolicyCategory, len(exists.Data))
	existIDs := make([]uuid.UUID, len(exists.Data))
	existKey := make([]string, len(exists.Data))
	for i, a := range exists.Data {
		existData[a.PolicyCategoryID] = service.PolicyCategory{
			Name:        a.Name,
			Key:         a.Key,
			Description: a.Description,
		}
		existIDs[i] = a.PolicyCategoryID
		existKey[i] = a.Key
	}
	var as []parameter.CreatePolicyCategoryParam
	for _, a := range service.PolicyCategories {
		matchIndex := contains(existKey, a.Key)
		if matchIndex == -1 {
			as = append(as, parameter.CreatePolicyCategoryParam{
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
					_, err = (*c.Manager).UpdatePolicyCategory(
						ctx,
						uid,
						a.Name,
						a.Key,
						a.Description,
					)
					if err != nil {
						return fmt.Errorf("failed to update policy category: %w", err)
					}
				}
			}
			delete(existData, uid)
		}
	}
	_, err = (*c.Manager).CreatePolicyCategories(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create policy categories: %w", err)
	}
	if len(existIDs) > 0 && !notDel {
		_, err = (*c.Manager).PluralDeletePolicyCategories(ctx, existIDs)
		if err != nil {
			return fmt.Errorf("failed to delete policy categories: %w", err)
		}
	}
	return nil
}
