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

type mimeType struct {
	MimeTypeID uuid.UUID `faker:"uuid,unique"`
	Key        string    `faker:"word,unique"`
	Name       string    `faker:"word"`
	Kind       string    `faker:"word"`
}

// MimeType is a slice of mimeType.
type MimeType []mimeType

// NewMimeTypes creates a new MimeType factory.
func (f *Factory) NewMimeTypes(num int) (MimeType, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	d := make([]mimeType, num)
	for i := 0; i < num; i++ { // Generate 5 structs having a unique word
		err := faker.FakeData(&d[i])
		if err != nil {
			return nil, fmt.Errorf("failed to generate fake data: %w", err)
		}
	}
	faker.ResetUnique() // Forget all generated unique values. Allows to start generating another unrelated dataset.
	return d, nil
}

// Copy returns a copy of MimeType.
func (d MimeType) Copy() MimeType {
	return append(MimeType{}, d...)
}

// LimitAndOffset returns a slice of MimeType with the given limit and offset.
func (d MimeType) LimitAndOffset(limit, offset int) MimeType {
	if len(d) < offset {
		return MimeType{}
	}
	if len(d) < offset+limit {
		return d[offset:]
	}
	return d[offset : offset+limit]
}

// ForCreateParam converts MimeType to []parameter.CreateMimeTypeParam.
func (d MimeType) ForCreateParam() []parameter.CreateMimeTypeParam {
	params := make([]parameter.CreateMimeTypeParam, len(d))
	for i, v := range d {
		params[i] = parameter.CreateMimeTypeParam{
			Key:  v.Key,
			Name: v.Name,
			Kind: v.Kind,
		}
	}
	return params
}

// ForEntity converts MimeType to []entity.MimeType.
func (d MimeType) ForEntity() []entity.MimeType {
	entities := make([]entity.MimeType, len(d))
	for i, v := range d {
		entities[i] = entity.MimeType{
			MimeTypeID: v.MimeTypeID,
			Key:        v.Key,
			Name:       v.Name,
			Kind:       v.Kind,
		}
	}
	return entities
}

// CountContainsName returns the number of mimeTypes that contain the given name.
func (d MimeType) CountContainsName(name string) int64 {
	count := int64(0)
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			count++
		}
	}
	return count
}

// FilterByName filters MimeType by name.
func (d MimeType) FilterByName(name string) MimeType {
	var res MimeType
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			res = append(res, v)
		}
	}
	return res
}

// OrderByNames sorts MimeType by name.
func (d MimeType) OrderByNames() MimeType {
	var res MimeType
	res = append(res, d...)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res
}

// OrderByReverseNames sorts MimeType by name in reverse order.
func (d MimeType) OrderByReverseNames() MimeType {
	var res MimeType
	res = append(res, d...)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})
	return res
}
