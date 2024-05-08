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

type eventType struct {
	EventTypeID uuid.UUID `faker:"uuid,unique"`
	Key         string    `faker:"word,unique"`
	Name        string    `faker:"word"`
	Color       string    `faker:"word"`
}

// EventType is a slice of eventType.
type EventType []eventType

// NewEventTypes creates a new EventType factory.
func NewEventTypes(num int) (EventType, error) {
	d := make([]eventType, num)
	for i := 0; i < num; i++ { // Generate 5 structs having a unique word
		err := faker.FakeData(&d[i])
		if err != nil {
			return nil, fmt.Errorf("failed to generate fake data: %w", err)
		}
	}
	faker.ResetUnique() // Forget all generated unique values. Allows to start generating another unrelated dataset.
	return d, nil
}

// Copy returns a copy of EventType.
func (d EventType) Copy() EventType {
	return append(EventType{}, d...)
}

// LimitAndOffset returns a slice of EventType with the given limit and offset.
func (d EventType) LimitAndOffset(limit, offset int) EventType {
	if len(d) < offset {
		return EventType{}
	}
	if len(d) < offset+limit {
		return d[offset:]
	}
	return d[offset : offset+limit]
}

// ForCreateParam converts EventType to []parameter.CreateEventTypeParam.
func (d EventType) ForCreateParam() []parameter.CreateEventTypeParam {
	params := make([]parameter.CreateEventTypeParam, len(d))
	for i, v := range d {
		params[i] = parameter.CreateEventTypeParam{
			Key:   v.Key,
			Name:  v.Name,
			Color: v.Color,
		}
	}
	return params
}

// ForEntity converts EventType to []entity.EventType.
func (d EventType) ForEntity() []entity.EventType {
	entities := make([]entity.EventType, len(d))
	for i, v := range d {
		entities[i] = entity.EventType{
			EventTypeID: v.EventTypeID,
			Key:         v.Key,
			Name:        v.Name,
			Color:       v.Color,
		}
	}
	return entities
}

// CountContainsName returns the number of eventTypes that contain the given name.
func (d EventType) CountContainsName(name string) int64 {
	count := int64(0)
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			count++
		}
	}
	return count
}

// FilterByName filters EventType by name.
func (d EventType) FilterByName(name string) EventType {
	var res EventType
	for _, v := range d {
		if strings.Contains(v.Name, name) {
			res = append(res, v)
		}
	}
	return res
}

// OrderByNames sorts EventType by name.
func (d EventType) OrderByNames() EventType {
	var res EventType
	res = append(res, d...)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res
}

// OrderByReverseNames sorts EventType by name in reverse order.
func (d EventType) OrderByReverseNames() EventType {
	var res EventType
	res = append(res, d...)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})
	return res
}
