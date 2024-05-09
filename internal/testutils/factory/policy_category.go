package factory

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

type policyCategory struct {
	PolicyCategoryID uuid.UUID `faker:"uuid,unique"`
	Key              string    `faker:"word,unique"`
	Name             string    `faker:"word"`
	Description      string    `faker:"sentence"`
}

// PolicyCategory is a slice of policyCategory.
type PolicyCategory []policyCategory

// NewPolicyCategories creates a new PolicyCategory factory.
func (f *Factory) NewPolicyCategories(num int) (PolicyCategory, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	d := make([]policyCategory, num)
	for i := 0; i < num; i++ { // Generate 5 structs having a unique word
		err := faker.FakeData(&d[i])
		if err != nil {
			return nil, fmt.Errorf("failed to generate fake data: %w", err)
		}
	}
	faker.ResetUnique() // Forget all generated unique values. Allows to start generating another unrelated dataset.
	return d, nil
}

// Copy returns a copy of PolicyCategory.
func (d PolicyCategory) Copy() PolicyCategory {
	return append(PolicyCategory{}, d...)
}

// LimitAndOffset returns a slice of PolicyCategory with the given limit and offset.
func (d PolicyCategory) LimitAndOffset(limit, offset int) PolicyCategory {
	if len(d) < offset {
		return PolicyCategory{}
	}
	if len(d) < offset+limit {
		return d[offset:]
	}
	return d[offset : offset+limit]
}

// ForCreateParam converts PolicyCategory to []parameter.CreatePolicyCategoryParam.
func (d PolicyCategory) ForCreateParam() []parameter.CreatePolicyCategoryParam {
	params := make([]parameter.CreatePolicyCategoryParam, len(d))
	for i, v := range d {
		params[i] = parameter.CreatePolicyCategoryParam{
			Key:         v.Key,
			Name:        v.Name,
			Description: v.Description,
		}
	}
	return params
}

// ForEntity converts PolicyCategory to []entity.PolicyCategory.
func (d PolicyCategory) ForEntity() []entity.PolicyCategory {
	entities := make([]entity.PolicyCategory, len(d))
	for i, v := range d {
		entities[i] = entity.PolicyCategory{
			PolicyCategoryID: v.PolicyCategoryID,
			Key:              v.Key,
			Name:             v.Name,
			Description:      v.Description,
		}
	}
	return entities
}

// CountContainsName returns the number of policyCategories that contain the given name.
func (d PolicyCategory) CountContainsName(name string) int64 {
	count := int64(0)
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			count++
		}
	}
	return count
}

// FilterByName filters PolicyCategory by name.
func (d PolicyCategory) FilterByName(name string) PolicyCategory {
	var res PolicyCategory
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			res = append(res, v)
		}
	}
	return res
}

// OrderByNames sorts PolicyCategory by name.
func (d PolicyCategory) OrderByNames() PolicyCategory {
	var res PolicyCategory
	res = append(res, d...)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res
}

// OrderByReverseNames sorts PolicyCategory by name in reverse order.
func (d PolicyCategory) OrderByReverseNames() PolicyCategory {
	var res PolicyCategory
	res = append(res, d...)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})
	return res
}
