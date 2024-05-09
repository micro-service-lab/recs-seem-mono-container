package batch

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
)

// InitAttendanceTypes is a batch to initialize attendance types.
type InitAttendanceTypes struct {
	Manager *service.ManagerInterface
}

// Run initializes attendance types.
func (c *InitAttendanceTypes) Run(ctx context.Context) error {
	var as []parameter.CreateAttendanceTypeParam
	for _, a := range service.AttendanceTypes {
		as = append(as, parameter.CreateAttendanceTypeParam{
			Name:  a.Name,
			Key:   a.Key,
			Color: a.Color,
		})
	}
	_, err := (*c.Manager).CreateAttendanceTypes(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create attendance types: %w", err)
	}
	return nil
}

// RunDiff run only if there is a difference.
func (c *InitAttendanceTypes) RunDiff(ctx context.Context, notDel, deepEqual bool) error {
	exists, err := (*c.Manager).GetAttendanceTypes(
		ctx,
		"",
		parameter.AttendanceTypeOrderMethodDefault,
		parameter.NonePagination,
		parameter.Limit(0),
		parameter.Cursor(""),
		parameter.Offset(0),
		parameter.WithCount(false),
	)
	if err != nil {
		return fmt.Errorf("failed to get attendance types: %w", err)
	}
	existData := make(map[uuid.UUID]service.AttendanceType, len(exists.Data))
	existIDs := make([]uuid.UUID, len(exists.Data))
	existKey := make([]string, len(exists.Data))
	for i, a := range exists.Data {
		existData[a.AttendanceTypeID] = service.AttendanceType{
			Name:  a.Name,
			Key:   a.Key,
			Color: a.Color,
		}
		existIDs[i] = a.AttendanceTypeID
		existKey[i] = a.Key
	}
	var as []parameter.CreateAttendanceTypeParam
	for _, a := range service.AttendanceTypes {
		matchIndex := contains(existKey, a.Key)
		if matchIndex == -1 {
			as = append(as, parameter.CreateAttendanceTypeParam{
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
					_, err = (*c.Manager).UpdateAttendanceType(
						ctx,
						uid,
						a.Name,
						a.Key,
						a.Color,
					)
					if err != nil {
						return fmt.Errorf("failed to update attendance type: %w", err)
					}
				}
			}
			delete(existData, uid)
		}
	}
	_, err = (*c.Manager).CreateAttendanceTypes(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create attendance types: %w", err)
	}
	if len(existIDs) > 0 && !notDel {
		err = (*c.Manager).PluralDeleteAttendanceTypes(ctx, existIDs)
		if err != nil {
			return fmt.Errorf("failed to delete attendance types: %w", err)
		}
	}
	return nil
}
