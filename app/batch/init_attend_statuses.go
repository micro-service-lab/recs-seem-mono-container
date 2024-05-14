package batch

import (
	"context"
	"fmt"

	"github.com/google/uuid"

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

// RunDiff run only if there is a difference.
func (c *InitAttendStatuses) RunDiff(ctx context.Context, notDel, deepEqual bool) error {
	exists, err := (*c.Manager).GetAttendStatuses(
		ctx,
		"",
		parameter.AttendStatusOrderMethodDefault,
		parameter.NonePagination,
		parameter.Limit(0),
		parameter.Cursor(""),
		parameter.Offset(0),
		parameter.WithCount(false),
	)
	if err != nil {
		return fmt.Errorf("failed to get attend statuses: %w", err)
	}
	existData := make(map[uuid.UUID]service.AttendStatus, len(exists.Data))
	existIDs := make([]uuid.UUID, len(exists.Data))
	existKey := make([]string, len(exists.Data))
	for i, a := range exists.Data {
		existData[a.AttendStatusID] = service.AttendStatus{
			Name: a.Name,
			Key:  a.Key,
		}
		existIDs[i] = a.AttendStatusID
		existKey[i] = a.Key
	}
	var as []parameter.CreateAttendStatusParam
	for _, a := range service.AttendStatuses {
		matchIndex := contains(existKey, a.Key)
		if matchIndex == -1 {
			as = append(as, parameter.CreateAttendStatusParam{
				Name: a.Name,
				Key:  a.Key,
			})
		} else {
			var uid uuid.UUID
			existIDs, uid = removeUUID(existIDs, matchIndex)
			existKey, _ = removeString(existKey, matchIndex)
			if deepEqual {
				de := isDeepEqual(a, existData[uid])
				if !de {
					_, err = (*c.Manager).UpdateAttendStatus(
						ctx,
						uid,
						a.Name,
						a.Key,
					)
					if err != nil {
						return fmt.Errorf("failed to update attend status: %w", err)
					}
				}
			}
			delete(existData, uid)
		}
	}
	_, err = (*c.Manager).CreateAttendStatuses(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create attend statuses: %w", err)
	}
	if len(existIDs) > 0 && !notDel {
		_, err = (*c.Manager).PluralDeleteAttendStatuses(ctx, existIDs)
		if err != nil {
			return fmt.Errorf("failed to delete attend statuses: %w", err)
		}
	}
	return nil
}
