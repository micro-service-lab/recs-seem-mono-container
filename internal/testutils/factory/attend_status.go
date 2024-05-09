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

type attendStatus struct {
	AttendStatusID uuid.UUID `faker:"uuid,unique"`
	Key            string    `faker:"word,unique"`
	Name           string    `faker:"word"`
}

// AttendStatuses is a slice of attendStatus.
type AttendStatuses []attendStatus

// NewAttendStatuses creates a new AttendStatuses factory.
func (f *Factory) NewAttendStatuses(num int) (AttendStatuses, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	d := make([]attendStatus, num)
	for i := 0; i < num; i++ { // Generate 5 structs having a unique word
		err := faker.FakeData(&d[i])
		if err != nil {
			return nil, fmt.Errorf("failed to generate fake data: %w", err)
		}
	}
	faker.ResetUnique() // Forget all generated unique values. Allows to start generating another unrelated dataset.
	return d, nil
}

// Copy returns a copy of AttendStatuses.
func (d AttendStatuses) Copy() AttendStatuses {
	return append(AttendStatuses{}, d...)
}

// LimitAndOffset returns a slice of AttendStatuses with the given limit and offset.
func (d AttendStatuses) LimitAndOffset(limit, offset int) AttendStatuses {
	if len(d) < offset {
		return AttendStatuses{}
	}
	if len(d) < offset+limit {
		return d[offset:]
	}
	return d[offset : offset+limit]
}

// ForCreateParam converts AttendStatuses to []parameter.CreateAttendStatusParam.
func (d AttendStatuses) ForCreateParam() []parameter.CreateAttendStatusParam {
	params := make([]parameter.CreateAttendStatusParam, len(d))
	for i, v := range d {
		params[i] = parameter.CreateAttendStatusParam{
			Key:  v.Key,
			Name: v.Name,
		}
	}
	return params
}

// ForEntity converts AttendStatuses to []entity.AttendStatus.
func (d AttendStatuses) ForEntity() []entity.AttendStatus {
	entities := make([]entity.AttendStatus, len(d))
	for i, v := range d {
		entities[i] = entity.AttendStatus{
			AttendStatusID: v.AttendStatusID,
			Key:            v.Key,
			Name:           v.Name,
		}
	}
	return entities
}

// CountContainsName returns the number of attendStatuses that contain the given name.
func (d AttendStatuses) CountContainsName(name string) int64 {
	count := int64(0)
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			count++
		}
	}
	return count
}

// FilterByName filters AttendStatuses by name.
func (d AttendStatuses) FilterByName(name string) AttendStatuses {
	var res AttendStatuses
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			res = append(res, v)
		}
	}
	return res
}

// OrderByNames sorts AttendStatuses by name.
func (d AttendStatuses) OrderByNames() AttendStatuses {
	var res AttendStatuses
	res = append(res, d...)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res
}

// OrderByReverseNames sorts AttendStatuses by name in reverse order.
func (d AttendStatuses) OrderByReverseNames() AttendStatuses {
	var res AttendStatuses
	res = append(res, d...)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})
	return res
}
