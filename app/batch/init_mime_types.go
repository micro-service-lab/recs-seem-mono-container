package batch

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
)

// InitMimeTypes is a batch to initialize mime types.
type InitMimeTypes struct {
	Manager *service.ManagerInterface
}

// Run initializes mime types.
func (c *InitMimeTypes) Run(ctx context.Context) error {
	var as []parameter.CreateMimeTypeParam
	for _, a := range service.MimeTypes {
		as = append(as, parameter.CreateMimeTypeParam{
			Name: a.Name,
			Key:  a.Key,
			Kind: a.Kind,
		})
	}
	_, err := (*c.Manager).CreateMimeTypes(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create mime types: %w", err)
	}
	return nil
}

// RunDiff run only if there is a difference.
func (c *InitMimeTypes) RunDiff(ctx context.Context, notDel, deepEqual bool) error {
	exists, err := (*c.Manager).GetMimeTypes(
		ctx,
		"",
		parameter.MimeTypeOrderMethodDefault,
		parameter.NonePagination,
		parameter.Limit(0),
		parameter.Cursor(""),
		parameter.Offset(0),
		parameter.WithCount(false),
	)
	if err != nil {
		return fmt.Errorf("failed to get mime types: %w", err)
	}
	existData := make(map[uuid.UUID]service.MimeType, len(exists.Data))
	existIDs := make([]uuid.UUID, len(exists.Data))
	existKey := make([]string, len(exists.Data))
	for i, a := range exists.Data {
		existData[a.MimeTypeID] = service.MimeType{
			Name: a.Name,
			Key:  a.Key,
			Kind: a.Kind,
		}
		existIDs[i] = a.MimeTypeID
		existKey[i] = a.Key
	}
	var as []parameter.CreateMimeTypeParam
	for _, a := range service.MimeTypes {
		matchIndex := contains(existKey, a.Key)
		if matchIndex == -1 {
			as = append(as, parameter.CreateMimeTypeParam{
				Name: a.Name,
				Key:  a.Key,
				Kind: a.Kind,
			})
		} else {
			var uid uuid.UUID
			existIDs, uid = removeUUID(existIDs, matchIndex)
			existKey, _ = removeString(existKey, matchIndex)
			if deepEqual {
				de := isDeepEqual(a, existData[uid])
				if !de {
					_, err = (*c.Manager).UpdateMimeType(
						ctx,
						uid,
						a.Name,
						a.Key,
						a.Kind,
					)
					if err != nil {
						return fmt.Errorf("failed to update mime type: %w", err)
					}
				}
			}
			delete(existData, uid)
		}
	}
	_, err = (*c.Manager).CreateMimeTypes(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create mime types: %w", err)
	}
	if len(existIDs) > 0 && !notDel {
		err = (*c.Manager).PluralDeleteMimeTypes(ctx, existIDs)
		if err != nil {
			return fmt.Errorf("failed to delete mime types: %w", err)
		}
	}
	return nil
}
