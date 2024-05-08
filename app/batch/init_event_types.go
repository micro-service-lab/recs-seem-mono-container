package batch

import (
	"context"
	"fmt"

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
