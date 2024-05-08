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

type attendanceType struct {
	AttendanceTypeID uuid.UUID `faker:"uuid,unique"`
	Key              string    `faker:"word,unique"`
	Name             string    `faker:"word"`
	Color            string    `faker:"word"`
}

// AttendanceType is a slice of attendanceType.
type AttendanceType []attendanceType

// NewAttendanceTypes creates a new AttendanceType factory.
func NewAttendanceTypes(num int) (AttendanceType, error) {
	d := make([]attendanceType, num)
	for i := 0; i < num; i++ { // Generate 5 structs having a unique word
		err := faker.FakeData(&d[i])
		if err != nil {
			return nil, fmt.Errorf("failed to generate fake data: %w", err)
		}
	}
	faker.ResetUnique() // Forget all generated unique values. Allows to start generating another unrelated dataset.
	return d, nil
}

// Copy returns a copy of AttendanceType.
func (d AttendanceType) Copy() AttendanceType {
	return append(AttendanceType{}, d...)
}

// LimitAndOffset returns a slice of AttendanceType with the given limit and offset.
func (d AttendanceType) LimitAndOffset(limit, offset int) AttendanceType {
	if len(d) < offset {
		return AttendanceType{}
	}
	if len(d) < offset+limit {
		return d[offset:]
	}
	return d[offset : offset+limit]
}

// ForCreateParam converts AttendanceType to []parameter.CreateAttendanceTypeParam.
func (d AttendanceType) ForCreateParam() []parameter.CreateAttendanceTypeParam {
	params := make([]parameter.CreateAttendanceTypeParam, len(d))
	for i, v := range d {
		params[i] = parameter.CreateAttendanceTypeParam{
			Key:   v.Key,
			Name:  v.Name,
			Color: v.Color,
		}
	}
	return params
}

// ForEntity converts AttendanceType to []entity.AttendanceType.
func (d AttendanceType) ForEntity() []entity.AttendanceType {
	entities := make([]entity.AttendanceType, len(d))
	for i, v := range d {
		entities[i] = entity.AttendanceType{
			AttendanceTypeID: v.AttendanceTypeID,
			Key:              v.Key,
			Name:             v.Name,
			Color:            v.Color,
		}
	}
	return entities
}

// CountContainsName returns the number of attendanceTypes that contain the given name.
func (d AttendanceType) CountContainsName(name string) int64 {
	count := int64(0)
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			count++
		}
	}
	return count
}

// FilterByName filters AttendanceType by name.
func (d AttendanceType) FilterByName(name string) AttendanceType {
	var res AttendanceType
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			res = append(res, v)
		}
	}
	return res
}

// OrderByNames sorts AttendanceType by name.
func (d AttendanceType) OrderByNames() AttendanceType {
	var res AttendanceType
	res = append(res, d...)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res
}

// OrderByReverseNames sorts AttendanceType by name in reverse order.
func (d AttendanceType) OrderByReverseNames() AttendanceType {
	var res AttendanceType
	res = append(res, d...)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})
	return res
}
