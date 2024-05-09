package batch

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
)

// InitEventTypes is a batch to initialize event types.
type InitEventTypes struct {
	Manager *service.ManagerInterface
}

// Run initializes event types.
func (c *InitEventTypes) Run(ctx context.Context) error {
	var as []parameter.CreateEventTypeParam
	for _, a := range service.EventTypes {
		as = append(as, parameter.CreateEventTypeParam{
			Name:  a.Name,
			Key:   a.Key,
			Color: a.Color,
		})
	}
	_, err := (*c.Manager).CreateEventTypes(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create event types: %w", err)
	}
	return nil
}

// RunDiff run only if there is a difference.
func (c *InitEventTypes) RunDiff(ctx context.Context, notDel, deepEqual bool) error {
	exists, err := (*c.Manager).GetEventTypes(
		ctx,
		"",
		parameter.EventTypeOrderMethodDefault,
		parameter.NonePagination,
		parameter.Limit(0),
		parameter.Cursor(""),
		parameter.Offset(0),
		parameter.WithCount(false),
	)
	if err != nil {
		return fmt.Errorf("failed to get event types: %w", err)
	}
	existData := make(map[uuid.UUID]service.EventType, len(exists.Data))
	existIDs := make([]uuid.UUID, len(exists.Data))
	existKey := make([]string, len(exists.Data))
	for i, a := range exists.Data {
		existData[a.EventTypeID] = service.EventType{
			Name:  a.Name,
			Key:   a.Key,
			Color: a.Color,
		}
		existIDs[i] = a.EventTypeID
		existKey[i] = a.Key
	}
	var as []parameter.CreateEventTypeParam
	for _, a := range service.EventTypes {
		matchIndex := contains(existKey, a.Key)
		if matchIndex == -1 {
			as = append(as, parameter.CreateEventTypeParam{
				Name:  a.Name,
				Key:   a.Key,
				Color: a.Color,
			})
		} else {
			var uid uuid.UUID
			existIDs, uid = removeUUID(existIDs, matchIndex)
			existKey, _ = removeString(existKey, matchIndex)
			if deepEqual {
				de := isDeepEqual(a, existData[uid])
				if !de {
					_, err = (*c.Manager).UpdateEventType(
						ctx,
						uid,
						a.Name,
						a.Key,
						a.Color,
					)
					if err != nil {
						return fmt.Errorf("failed to update event type: %w", err)
					}
				}
			}
			delete(existData, uid)
		}
	}
	_, err = (*c.Manager).CreateEventTypes(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create event types: %w", err)
	}
	if len(existIDs) > 0 && !notDel {
		err = (*c.Manager).PluralDeleteEventTypes(ctx, existIDs)
		if err != nil {
			return fmt.Errorf("failed to delete event types: %w", err)
		}
	}
	return nil
}
