package factory

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-faker/faker/v4"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

type attendStatus struct {
	Key  string `faker:"word,unique"`
	Name string `faker:"word"`
}

// AttendStatuses is a slice of attendStatus.
type AttendStatuses []attendStatus

// NewAttendStatuses creates a new AttendStatuses factory.
func NewAttendStatuses(num int) (AttendStatuses, error) {
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

// OrderByNames sorts AttendStatuses by name.
func (d AttendStatuses) OrderByNames() AttendStatuses {
	sort.Slice(d, func(i, j int) bool {
		return d[i].Name < d[j].Name
	})
	return d
}
