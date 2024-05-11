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

type permission struct {
	PermissionID       uuid.UUID          `faker:"uuid,unique"`
	Key                string             `faker:"word,unique"`
	Name               string             `faker:"word"`
	Description        string             `faker:"sentence"`
	PermissionCategory permissionCategory `faker:"-"`
}

// Permission is a slice of permission.
type Permission []permission

// NewPermissions creates a new Permission factory.
func (f *Factory) NewPermissions(num int) (Permission, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	d := make([]permission, num)
	for i := 0; i < num; i++ { // Generate 5 structs having a unique word
		err := faker.FakeData(&d[i])
		if err != nil {
			return Permission{}, fmt.Errorf("failed to generate fake data: %w", err)
		}
	}
	faker.ResetUnique() // Forget all generated unique values. Allows to start generating another unrelated Dataset.
	return d, nil
}

// WithPermissionCategory adds a permissionCategory to Permission.
func (d Permission) WithPermissionCategory(categories []permissionCategory) (Permission, error) {
	res := d.Copy()
	for i := range res {
		randI, err := rand.Int(rand.Reader, big.NewInt(int64(len(categories))))
		if err != nil {
			return Permission{}, fmt.Errorf("failed to generate fake data: %w", err)
		}
		res[i].PermissionCategory = categories[randI.Int64()]
	}
	return res, nil
}

// Copy returns a copy of Permission.
func (d Permission) Copy() Permission {
	return append(Permission{}, d...)
}

// LimitAndOffset returns a slice of Permission with the given limit and offset.
func (d Permission) LimitAndOffset(limit, offset int) Permission {
	if len(d) < offset {
		return d
	}
	if len(d) < offset+limit {
		return d[offset:]
	}
	return d[offset : offset+limit]
}

// ForCreateParam converts Permission to []parameter.CreatePermissionParam.
func (d Permission) ForCreateParam() []parameter.CreatePermissionParam {
	params := make([]parameter.CreatePermissionParam, len(d))
	for i, v := range d {
		params[i] = parameter.CreatePermissionParam{
			Key:                  v.Key,
			Name:                 v.Name,
			Description:          v.Description,
			PermissionCategoryID: v.PermissionCategory.PermissionCategoryID,
		}
	}
	return params
}

// ForEntity converts Permission to []entity.Permission.
func (d Permission) ForEntity() []entity.Permission {
	entities := make([]entity.Permission, len(d))
	for i, v := range d {
		entities[i] = entity.Permission{
			PermissionID:         v.PermissionID,
			Key:                  v.Key,
			Name:                 v.Name,
			Description:          v.Description,
			PermissionCategoryID: v.PermissionCategory.PermissionCategoryID,
		}
	}
	return entities
}

// CountContainsName returns the number of permissionCategories that contain the given name.
func (d Permission) CountContainsName(name string) int64 {
	count := int64(0)
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			count++
		}
	}
	return count
}

// FilterByName filters Permission by name.
func (d Permission) FilterByName(name string) Permission {
	var res Permission
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			res = append(res, v)
		}
	}
	return res
}

// FilterByPermissionCategories filters Permission by permissionCategories.
func (d Permission) FilterByPermissionCategories(categoryIDs []uuid.UUID) Permission {
	var res Permission
	for _, v := range d {
		for _, id := range categoryIDs {
			if v.PermissionCategory.PermissionCategoryID == id {
				res = append(res, v)
			}
		}
	}
	return res
}

// OrderByNames sorts Permission by name.
func (d Permission) OrderByNames() Permission {
	res := d.Copy()
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res
}

// OrderByReverseNames sorts Permission by name in reverse order.
func (d Permission) OrderByReverseNames() Permission {
	res := d.Copy()
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})
	return res
}
