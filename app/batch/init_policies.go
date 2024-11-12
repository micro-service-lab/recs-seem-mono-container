package batch

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
)

// InitPolicies is a batch to initialize policy.
type InitPolicies struct {
	Manager *service.ManagerInterface
}

// Run initializes policy.
func (c *InitPolicies) Run(ctx context.Context) error {
	searchPC, _, err := searchPolicyCategoryID(ctx, c.Manager)
	if err != nil {
		return fmt.Errorf("failed to search policy category: %w", err)
	}
	var as []parameter.CreatePolicyParam
	for _, a := range service.Policies {
		cID := searchPC(a.PolicyCategoryKey)
		if cID == uuid.Nil {
			return fmt.Errorf("failed to find policy category: %s", a.PolicyCategoryKey)
		}
		as = append(as, parameter.CreatePolicyParam{
			Name:             a.Name,
			Key:              a.Key,
			Description:      a.Description,
			PolicyCategoryID: cID,
		})
	}
	_, err = (*c.Manager).CreatePolicies(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create policy: %w", err)
	}
	return nil
}

// RunDiff run only if there is a difference.
func (c *InitPolicies) RunDiff(ctx context.Context, notDel, deepEqual bool) error {
	searchPcID, searchPcKey, err := searchPolicyCategoryID(ctx, c.Manager)
	if err != nil {
		return fmt.Errorf("failed to search policy category: %w", err)
	}
	exists, err := (*c.Manager).GetPolicies(
		ctx,
		"",
		[]uuid.UUID{},
		parameter.PolicyOrderMethodDefault,
		parameter.NonePagination,
		parameter.Limit(0),
		parameter.Cursor(""),
		parameter.Offset(0),
		parameter.WithCount(false),
	)
	if err != nil {
		return fmt.Errorf("failed to get policy: %w", err)
	}
	existData := make(map[uuid.UUID]service.Policy, len(exists.Data))
	existIDs := make([]uuid.UUID, len(exists.Data))
	existKey := make([]string, len(exists.Data))
	for i, a := range exists.Data {
		key := searchPcKey(a.PolicyCategoryID)
		if key == "" {
			return fmt.Errorf("failed to find policy category: %s", a.PolicyCategoryID)
		}
		existData[a.PolicyID] = service.Policy{
			Name:              a.Name,
			Key:               a.Key,
			Description:       a.Description,
			PolicyCategoryKey: key,
		}
		existIDs[i] = a.PolicyID
		existKey[i] = a.Key
	}
	var as []parameter.CreatePolicyParam
	for _, a := range service.Policies {
		matchIndex := contains(existKey, a.Key)
		pcID := searchPcID(a.PolicyCategoryKey)
		if pcID == uuid.Nil {
			return fmt.Errorf("failed to find policy category: %s", a.PolicyCategoryKey)
		}
		if matchIndex == -1 {
			as = append(as, parameter.CreatePolicyParam{
				Name:             a.Name,
				Key:              a.Key,
				Description:      a.Description,
				PolicyCategoryID: pcID,
			})
		} else {
			var uid uuid.UUID
			existIDs, uid = removeUUID(existIDs, matchIndex)
			existKey, _ = removeString(existKey, matchIndex)
			if deepEqual {
				de := isDeepEqual(a, existData[uid])
				if !de {
					_, err = (*c.Manager).UpdatePolicy(
						ctx,
						uid,
						a.Name,
						a.Key,
						a.Description,
						pcID,
					)
					if err != nil {
						return fmt.Errorf("failed to update policy category: %w", err)
					}
				}
			}
			delete(existData, uid)
		}
	}
	_, err = (*c.Manager).CreatePolicies(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create policy: %w", err)
	}
	if len(existIDs) > 0 && !notDel {
		_, err = (*c.Manager).PluralDeletePolicies(ctx, existIDs)
		if err != nil {
			return fmt.Errorf("failed to delete policy: %w", err)
		}
	}
	return nil
}

func searchPolicyCategoryID(
	ctx context.Context,
	c *service.ManagerInterface,
) (func(key service.PolicyCategoryKey) uuid.UUID,
	func(id uuid.UUID) service.PolicyCategoryKey,
	error,
) {
	pc, err := (*c).GetPolicyCategories(
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
		return nil, nil, fmt.Errorf("failed to get policy category: %w", err)
	}
	searchPcID := func(key service.PolicyCategoryKey) uuid.UUID {
		for _, a := range pc.Data {
			if a.Key == string(key) {
				return a.PolicyCategoryID
			}
		}
		return uuid.Nil
	}
	searchPCKey := func(id uuid.UUID) service.PolicyCategoryKey {
		for _, a := range pc.Data {
			if a.PolicyCategoryID == id {
				return service.PolicyCategoryKey(a.Key)
			}
		}
		return ""
	}
	return searchPcID, searchPCKey, nil
}
