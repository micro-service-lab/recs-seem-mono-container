package batch

import (
	"context"
	"fmt"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
)

// InitAttendStatuses is a batch to initialize attend statuses.
type InitAttendStatuses struct {
	Manager *service.ManagerInterface
}

// Run initializes attend statuses.
func (c *InitAttendStatuses) Run(ctx context.Context) error {
	var as []parameter.CreateAttendStatusParam
	for _, a := range service.AttendStatuses {
		as = append(as, parameter.CreateAttendStatusParam{
			Name: a.Name,
			Key:  a.Key,
		})
	}
	_, err := (*c.Manager).CreateAttendStatuses(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create attend statuses: %w", err)
	}
	return nil
}
