package batch

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
)

// InitChatRoomActionTypes is a batch to initialize chat room action types.
type InitChatRoomActionTypes struct {
	Manager *service.ManagerInterface
}

// Run initializes chat room action types.
func (c *InitChatRoomActionTypes) Run(ctx context.Context) error {
	var as []parameter.CreateChatRoomActionTypeParam
	for _, a := range service.ChatRoomActionTypes {
		as = append(as, parameter.CreateChatRoomActionTypeParam{
			Name: a.Name,
			Key:  a.Key,
		})
	}
	_, err := (*c.Manager).CreateChatRoomActionTypes(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create chat room action types: %w", err)
	}
	return nil
}

// RunDiff run only if there is a difference.
func (c *InitChatRoomActionTypes) RunDiff(ctx context.Context, notDel, deepEqual bool) error {
	exists, err := (*c.Manager).GetChatRoomActionTypes(
		ctx,
		"",
		parameter.ChatRoomActionTypeOrderMethodDefault,
		parameter.NonePagination,
		parameter.Limit(0),
		parameter.Cursor(""),
		parameter.Offset(0),
		parameter.WithCount(false),
	)
	if err != nil {
		return fmt.Errorf("failed to get chat room action types: %w", err)
	}
	existData := make(map[uuid.UUID]service.ChatRoomActionType, len(exists.Data))
	existIDs := make([]uuid.UUID, len(exists.Data))
	existKey := make([]string, len(exists.Data))
	for i, a := range exists.Data {
		existData[a.ChatRoomActionTypeID] = service.ChatRoomActionType{
			Name: a.Name,
			Key:  a.Key,
		}
		existIDs[i] = a.ChatRoomActionTypeID
		existKey[i] = a.Key
	}
	var as []parameter.CreateChatRoomActionTypeParam
	for _, a := range service.ChatRoomActionTypes {
		matchIndex := contains(existKey, a.Key)
		if matchIndex == -1 {
			as = append(as, parameter.CreateChatRoomActionTypeParam{
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
					_, err = (*c.Manager).UpdateChatRoomActionType(
						ctx,
						uid,
						a.Name,
						a.Key,
					)
					if err != nil {
						return fmt.Errorf("failed to update chat room action type: %w", err)
					}
				}
			}
			delete(existData, uid)
		}
	}
	_, err = (*c.Manager).CreateChatRoomActionTypes(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create chat room action types: %w", err)
	}
	if len(existIDs) > 0 && !notDel {
		_, err = (*c.Manager).PluralDeleteChatRoomActionTypes(ctx, existIDs)
		if err != nil {
			return fmt.Errorf("failed to delete chat room action types: %w", err)
		}
	}
	return nil
}
