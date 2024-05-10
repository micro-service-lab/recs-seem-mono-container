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

type permission struct {
	PermissionID       uuid.UUID          `faker:"uuid,unique"`
	Key                string             `faker:"word,unique"`
	Name               string             `faker:"word"`
	Description        string             `faker:"sentence"`
	PermissionCategory permissionCategory `faker:"-"`
}

// Permission is a struct for Permission factory.
type Permission struct {
	f    *Factory
	data []permission
}

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
	faker.ResetUnique() // Forget all generated unique values. Allows to start generating another unrelated dataset.
	return Permission{data: d, f: f}, nil
}

// WithPermissionCategory adds a permissionCategory to Permission.
func (d Permission) WithPermissionCategory(categoryNum int) (Permission, error) {
	cd, err := d.f.NewPermissionCategories(categoryNum)
	if err != nil {
		return Permission{}, fmt.Errorf("failed to generate fake data: %w", err)
	}
	res := d.Copy()
	for i := range res.data {
		randI, err := faker.RandomInt(0, categoryNum-1, 1)
		if err != nil {
			return Permission{}, fmt.Errorf("failed to generate fake data: %w", err)
		}
		res.data[i].PermissionCategory = cd[randI[0]]
	}
	return res, nil
}

// Copy returns a copy of Permission.
func (d Permission) Copy() Permission {
	var res Permission
	res.data = append(res.data, d.data...)
	return res
}

// LimitAndOffset returns a slice of Permission with the given limit and offset.
func (d Permission) LimitAndOffset(limit, offset int) Permission {
	if len(d.data) < offset {
		return d
	}
	if len(d.data) < offset+limit {
		return Permission{data: d.data[offset:], f: d.f}
	}
	return Permission{data: d.data[offset : offset+limit], f: d.f}
}

// ForCreateParam converts Permission to []parameter.CreatePermissionParam.
func (d Permission) ForCreateParam() []parameter.CreatePermissionParam {
	params := make([]parameter.CreatePermissionParam, len(d.data))
	for i, v := range d.data {
		params[i] = parameter.CreatePermissionParam{
			Key:         v.Key,
			Name:        v.Name,
			Description: v.Description,
		}
	}
	return params
}

// ForEntity converts Permission to []entity.Permission.
func (d Permission) ForEntity() []entity.Permission {
	entities := make([]entity.Permission, len(d.data))
	for i, v := range d.data {
		entities[i] = entity.Permission{
			PermissionID: v.PermissionID,
			Key:          v.Key,
			Name:         v.Name,
			Description:  v.Description,
		}
	}
	return entities
}

// CountContainsName returns the number of permissionCategories that contain the given name.
func (d Permission) CountContainsName(name string) int64 {
	count := int64(0)
	for _, v := range d.data {
		if strings.Contains(v.Name, name) {
			count++
		}
	}
	return count
}

// FilterByName filters Permission by name.
func (d Permission) FilterByName(name string) Permission {
	var res Permission
	res.f = d.f
	for _, v := range d.data {
		if strings.Contains(v.Name, name) {
			res.data = append(res.data, v)
		}
	}
	return res
}

// OrderByNames sorts Permission by name.
func (d Permission) OrderByNames() Permission {
	var res Permission
	res.f = d.f
	res.data = append(res.data, d.data...)
	sort.Slice(res.data, func(i, j int) bool {
		return res.data[i].Name < res.data[j].Name
	})
	return res
}

// OrderByReverseNames sorts Permission by name in reverse order.
func (d Permission) OrderByReverseNames() Permission {
	var res Permission
	res.data = append(res.data, d.data...)
	sort.Slice(res.data, func(i, j int) bool {
		return res.data[i].Name > res.data[j].Name
	})
	return res
}
