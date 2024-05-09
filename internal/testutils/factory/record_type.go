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

type recordType struct {
	RecordTypeID uuid.UUID `faker:"uuid,unique"`
	Key          string    `faker:"word,unique"`
	Name         string    `faker:"word"`
}

// RecordType is a slice of recordType.
type RecordType []recordType

// NewRecordTypes creates a new RecordType factory.
func (f *Factory) NewRecordTypes(num int) (RecordType, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	d := make([]recordType, num)
	for i := 0; i < num; i++ { // Generate 5 structs having a unique word
		err := faker.FakeData(&d[i])
		if err != nil {
			return nil, fmt.Errorf("failed to generate fake data: %w", err)
		}
	}
	faker.ResetUnique() // Forget all generated unique values. Allows to start generating another unrelated dataset.
	return d, nil
}

// Copy returns a copy of RecordType.
func (d RecordType) Copy() RecordType {
	return append(RecordType{}, d...)
}

// LimitAndOffset returns a slice of RecordType with the given limit and offset.
func (d RecordType) LimitAndOffset(limit, offset int) RecordType {
	if len(d) < offset {
		return RecordType{}
	}
	if len(d) < offset+limit {
		return d[offset:]
	}
	return d[offset : offset+limit]
}

// ForCreateParam converts RecordType to []parameter.CreateRecordTypeParam.
func (d RecordType) ForCreateParam() []parameter.CreateRecordTypeParam {
	params := make([]parameter.CreateRecordTypeParam, len(d))
	for i, v := range d {
		params[i] = parameter.CreateRecordTypeParam{
			Key:  v.Key,
			Name: v.Name,
		}
	}
	return params
}

// ForEntity converts RecordType to []entity.RecordType.
func (d RecordType) ForEntity() []entity.RecordType {
	entities := make([]entity.RecordType, len(d))
	for i, v := range d {
		entities[i] = entity.RecordType{
			RecordTypeID: v.RecordTypeID,
			Key:          v.Key,
			Name:         v.Name,
		}
	}
	return entities
}

// CountContainsName returns the number of recordTypes that contain the given name.
func (d RecordType) CountContainsName(name string) int64 {
	count := int64(0)
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			count++
		}
	}
	return count
}

// FilterByName filters RecordType by name.
func (d RecordType) FilterByName(name string) RecordType {
	var res RecordType
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			res = append(res, v)
		}
	}
	return res
}

// OrderByNames sorts RecordType by name.
func (d RecordType) OrderByNames() RecordType {
	var res RecordType
	res = append(res, d...)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res
}

// OrderByReverseNames sorts RecordType by name in reverse order.
func (d RecordType) OrderByReverseNames() RecordType {
	var res RecordType
	res = append(res, d...)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})
	return res
}
