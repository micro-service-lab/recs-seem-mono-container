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

type permissionCategory struct {
	PermissionCategoryID uuid.UUID `faker:"uuid,unique"`
	Key                  string    `faker:"word,unique"`
	Name                 string    `faker:"word"`
	Description          string    `faker:"sentence"`
}

// PermissionCategory is a slice of permissionCategory.
type PermissionCategory []permissionCategory

// NewPermissionCategories creates a new PermissionCategory factory.
func NewPermissionCategories(num int) (PermissionCategory, error) {
	d := make([]permissionCategory, num)
	for i := 0; i < num; i++ { // Generate 5 structs having a unique word
		err := faker.FakeData(&d[i])
		if err != nil {
			return nil, fmt.Errorf("failed to generate fake data: %w", err)
		}
	}
	faker.ResetUnique() // Forget all generated unique values. Allows to start generating another unrelated dataset.
	return d, nil
}

// Copy returns a copy of PermissionCategory.
func (d PermissionCategory) Copy() PermissionCategory {
	return append(PermissionCategory{}, d...)
}

// LimitAndOffset returns a slice of PermissionCategory with the given limit and offset.
func (d PermissionCategory) LimitAndOffset(limit, offset int) PermissionCategory {
	if len(d) < offset {
		return PermissionCategory{}
	}
	if len(d) < offset+limit {
		return d[offset:]
	}
	return d[offset : offset+limit]
}

// ForCreateParam converts PermissionCategory to []parameter.CreatePermissionCategoryParam.
func (d PermissionCategory) ForCreateParam() []parameter.CreatePermissionCategoryParam {
	params := make([]parameter.CreatePermissionCategoryParam, len(d))
	for i, v := range d {
		params[i] = parameter.CreatePermissionCategoryParam{
			Key:         v.Key,
			Name:        v.Name,
			Description: v.Description,
		}
	}
	return params
}

// ForEntity converts PermissionCategory to []entity.PermissionCategory.
func (d PermissionCategory) ForEntity() []entity.PermissionCategory {
	entities := make([]entity.PermissionCategory, len(d))
	for i, v := range d {
		entities[i] = entity.PermissionCategory{
			PermissionCategoryID: v.PermissionCategoryID,
			Key:                  v.Key,
			Name:                 v.Name,
			Description:          v.Description,
		}
	}
	return entities
}

// CountContainsName returns the number of permissionCategories that contain the given name.
func (d PermissionCategory) CountContainsName(name string) int64 {
	count := int64(0)
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			count++
		}
	}
	return count
}

// FilterByName filters PermissionCategory by name.
func (d PermissionCategory) FilterByName(name string) PermissionCategory {
	var res PermissionCategory
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			res = append(res, v)
		}
	}
	return res
}

// OrderByNames sorts PermissionCategory by name.
func (d PermissionCategory) OrderByNames() PermissionCategory {
	var res PermissionCategory
	res = append(res, d...)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res
}

// OrderByReverseNames sorts PermissionCategory by name in reverse order.
func (d PermissionCategory) OrderByReverseNames() PermissionCategory {
	var res PermissionCategory
	res = append(res, d...)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})
	return res
}
