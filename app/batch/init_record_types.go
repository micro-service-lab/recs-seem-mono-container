package batch

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
)

// InitRecordTypes is a batch to initialize record types.
type InitRecordTypes struct {
	Manager *service.ManagerInterface
}

// Run initializes record types.
func (c *InitRecordTypes) Run(ctx context.Context) error {
	var as []parameter.CreateRecordTypeParam
	for _, a := range service.RecordTypes {
		as = append(as, parameter.CreateRecordTypeParam{
			Name: a.Name,
			Key:  a.Key,
		})
	}
	_, err := (*c.Manager).CreateRecordTypes(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create record types: %w", err)
	}
	return nil
}

// RunDiff run only if there is a difference.
func (c *InitRecordTypes) RunDiff(ctx context.Context, notDel, deepEqual bool) error {
	exists, err := (*c.Manager).GetRecordTypes(
		ctx,
		"",
		parameter.RecordTypeOrderMethodDefault,
		parameter.NonePagination,
		parameter.Limit(0),
		parameter.Cursor(""),
		parameter.Offset(0),
		parameter.WithCount(false),
	)
	if err != nil {
		return fmt.Errorf("failed to get record types: %w", err)
	}
	existData := make(map[uuid.UUID]service.RecordType, len(exists.Data))
	existIDs := make([]uuid.UUID, len(exists.Data))
	existKey := make([]string, len(exists.Data))
	for i, a := range exists.Data {
		existData[a.RecordTypeID] = service.RecordType{
			Name: a.Name,
			Key:  a.Key,
		}
		existIDs[i] = a.RecordTypeID
		existKey[i] = a.Key
	}
	var as []parameter.CreateRecordTypeParam
	for _, a := range service.RecordTypes {
		matchIndex := contains(existKey, a.Key)
		if matchIndex == -1 {
			as = append(as, parameter.CreateRecordTypeParam{
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
					_, err = (*c.Manager).UpdateRecordType(
						ctx,
						uid,
						a.Name,
						a.Key,
					)
					if err != nil {
						return fmt.Errorf("failed to update record type: %w", err)
					}
				}
			}
			delete(existData, uid)
		}
	}
	_, err = (*c.Manager).CreateRecordTypes(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create record types: %w", err)
	}
	if len(existIDs) > 0 && !notDel {
		_, err = (*c.Manager).PluralDeleteRecordTypes(ctx, existIDs)
		if err != nil {
			return fmt.Errorf("failed to delete record types: %w", err)
		}
	}
	return nil
}
