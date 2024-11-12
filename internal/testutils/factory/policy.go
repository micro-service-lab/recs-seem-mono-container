package factory

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

type policy struct {
	PolicyID       uuid.UUID      `faker:"uuid,unique"`
	Key            string         `faker:"word,unique"`
	Name           string         `faker:"word"`
	Description    string         `faker:"sentence"`
	PolicyCategory policyCategory `faker:"-"`
}

// Policy is a slice of policy.
type Policy []policy

// NewPolicies creates a new Policy factory.
func (f *Factory) NewPolicies(num int) (Policy, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	d := make([]policy, num)
	for i := 0; i < num; i++ { // Generate 5 structs having a unique word
		err := faker.FakeData(&d[i])
		if err != nil {
			return Policy{}, fmt.Errorf("failed to generate fake data: %w", err)
		}
	}
	faker.ResetUnique() // Forget all generated unique values. Allows to start generating another unrelated Dataset.
	return d, nil
}

// WithPolicyCategory adds a policyCategory to Policy.
func (d Policy) WithPolicyCategory(categories []policyCategory) (Policy, error) {
	res := d.Copy()
	for i := range res {
		randI, err := rand.Int(rand.Reader, big.NewInt(int64(len(categories))))
		if err != nil {
			return Policy{}, fmt.Errorf("failed to generate fake data: %w", err)
		}
		res[i].PolicyCategory = categories[randI.Int64()]
	}
	return res, nil
}

// Copy returns a copy of Policy.
func (d Policy) Copy() Policy {
	return append(Policy{}, d...)
}

// LimitAndOffset returns a slice of Policy with the given limit and offset.
func (d Policy) LimitAndOffset(limit, offset int) Policy {
	if len(d) < offset {
		return d
	}
	if len(d) < offset+limit {
		return d[offset:]
	}
	return d[offset : offset+limit]
}

// ForCreateParam converts Policy to []parameter.CreatePolicyParam.
func (d Policy) ForCreateParam() []parameter.CreatePolicyParam {
	params := make([]parameter.CreatePolicyParam, len(d))
	for i, v := range d {
		params[i] = parameter.CreatePolicyParam{
			Key:              v.Key,
			Name:             v.Name,
			Description:      v.Description,
			PolicyCategoryID: v.PolicyCategory.PolicyCategoryID,
		}
	}
	return params
}

// ForEntity converts Policy to []entity.Policy.
func (d Policy) ForEntity() []entity.Policy {
	entities := make([]entity.Policy, len(d))
	for i, v := range d {
		entities[i] = entity.Policy{
			PolicyID:         v.PolicyID,
			Key:              v.Key,
			Name:             v.Name,
			Description:      v.Description,
			PolicyCategoryID: v.PolicyCategory.PolicyCategoryID,
		}
	}
	return entities
}

// ForEntityWithPolicyCategory converts Policy to []entity.PolicyWithCategory.
func (d Policy) ForEntityWithPolicyCategory() []entity.PolicyWithCategory {
	entities := make([]entity.PolicyWithCategory, len(d))
	for i, v := range d {
		entities[i] = entity.PolicyWithCategory{
			Policy: entity.Policy{
				PolicyID:         v.PolicyID,
				Key:              v.Key,
				Name:             v.Name,
				Description:      v.Description,
				PolicyCategoryID: v.PolicyCategory.PolicyCategoryID,
			},
			PolicyCategory: entity.PolicyCategory{
				PolicyCategoryID: v.PolicyCategory.PolicyCategoryID,
				Name:             v.PolicyCategory.Name,
				Description:      v.PolicyCategory.Description,
			},
		}
	}
	return entities
}

// CountContainsName returns the number of policyCategories that contain the given name.
func (d Policy) CountContainsName(name string) int64 {
	count := int64(0)
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			count++
		}
	}
	return count
}

// FilterByName filters Policy by name.
func (d Policy) FilterByName(name string) Policy {
	var res Policy
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			res = append(res, v)
		}
	}
	return res
}

// FilterByPolicyCategories filters Policy by policyCategories.
func (d Policy) FilterByPolicyCategories(categoryIDs []uuid.UUID) Policy {
	var res Policy
	for _, v := range d {
		for _, id := range categoryIDs {
			if v.PolicyCategory.PolicyCategoryID == id {
				res = append(res, v)
			}
		}
	}
	return res
}

// FilterByIDs filters Policy by policyIDs.
func (d Policy) FilterByIDs(policyIDs []uuid.UUID) Policy {
	var res Policy
	for _, v := range d {
		for _, id := range policyIDs {
			if v.PolicyID == id {
				res = append(res, v)
			}
		}
	}
	return res
}

// Delete deletes the policy with the given policyID.
func (d Policy) Delete(policyID uuid.UUID) Policy {
	var res Policy
	for _, v := range d {
		if v.PolicyID != policyID {
			res = append(res, v)
		}
	}
	return res
}

// Update updates the policy with the given policyID.
func (d Policy) Update(policyID uuid.UUID, update parameter.UpdatePolicyParams) Policy {
	var res Policy
	for _, v := range d {
		if v.PolicyID == policyID {
			v.Name = update.Name
			v.Key = update.Key
			v.Description = update.Description
			res = append(res, v)
		} else {
			res = append(res, v)
		}
	}
	return res
}

// UpdateByKey updates the policy with the given key.
func (d Policy) UpdateByKey(key string, update parameter.UpdatePolicyByKeyParams) Policy {
	var res Policy
	for _, v := range d {
		if v.Key == key {
			v.Name = update.Name
			v.Description = update.Description
			res = append(res, v)
		} else {
			res = append(res, v)
		}
	}
	return res
}

// Len returns the length of Policy.
func (d Policy) Len() int64 {
	return int64(len(d))
}

// OrderByNames sorts Policy by name.
func (d Policy) OrderByNames() Policy {
	res := d.Copy()
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res
}

// OrderByReverseNames sorts Policy by name in reverse order.
func (d Policy) OrderByReverseNames() Policy {
	res := d.Copy()
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})
	return res
}
