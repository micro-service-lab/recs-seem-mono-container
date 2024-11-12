package batch

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/app/storage"
)

// InitGrades is a batch to initialize grades.
type InitGrades struct {
	Manager *service.ManagerInterface
	Storage storage.Storage
}

// Run initializes grades.
func (c *InitGrades) Run(ctx context.Context) error {
	var as []parameter.CreateGradeServiceParam
	shuffleImages := service.DefaultImageKeys
	service.ArrayShuffle(shuffleImages)
	for i, a := range service.Grades {
		url, err := c.Storage.GetURLFromKey(ctx, shuffleImages[i%len(shuffleImages)])
		if err != nil {
			return fmt.Errorf("failed to get url from key: %w", err)
		}
		attachableItem, err := (*c.Manager).FindAttachableItemByURL(ctx, url)
		if err != nil {
			return fmt.Errorf("failed to find attachable item by url: %w", err)
		}
		var imageID entity.UUID
		if attachableItem.Image.Valid {
			imageID = entity.UUID{
				Bytes: attachableItem.Image.Entity.ImageID,
				Valid: true,
			}
		}
		as = append(as, parameter.CreateGradeServiceParam{
			Name:         a.Name,
			Key:          a.Key,
			Description:  entity.String{String: a.Description, Valid: true},
			Color:        entity.String{String: a.Color, Valid: true},
			CoverImageID: imageID,
		})
	}
	_, err := (*c.Manager).CreateGrades(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create grades: %w", err)
	}
	return nil
}

// RunDiff run only if there is a difference.
func (c *InitGrades) RunDiff(ctx context.Context, notDel, deepEqual bool) error {
	shuffleImages := service.DefaultImageKeys
	service.ArrayShuffle(shuffleImages)
	exists, err := (*c.Manager).GetGradesWithOrganization(
		ctx,
		parameter.GradeOrderMethodDefault,
		parameter.NonePagination,
		parameter.Limit(0),
		parameter.Cursor(""),
		parameter.Offset(0),
		parameter.WithCount(false),
	)
	if err != nil {
		return fmt.Errorf("failed to get grades: %w", err)
	}
	existData := make(map[uuid.UUID]service.Grade, len(exists.Data))
	existIDs := make([]uuid.UUID, len(exists.Data))
	existKey := make([]string, len(exists.Data))
	for i, a := range exists.Data {
		existData[a.GradeID] = service.Grade{
			Name:        a.Organization.Name,
			Key:         a.Key,
			Description: a.Organization.Description.String,
			Color:       a.Organization.Color.String,
		}
		existIDs[i] = a.GradeID
		existKey[i] = a.Key
	}
	var as []parameter.CreateGradeServiceParam
	for i, a := range service.Grades {
		url, err := c.Storage.GetURLFromKey(ctx, shuffleImages[i%len(shuffleImages)])
		if err != nil {
			return fmt.Errorf("failed to get url from key: %w", err)
		}
		attachableItem, err := (*c.Manager).FindAttachableItemByURL(ctx, url)
		if err != nil {
			return fmt.Errorf("failed to find attachable item by url: %w", err)
		}
		var imageID entity.UUID
		if attachableItem.Image.Valid {
			imageID = entity.UUID{
				Bytes: attachableItem.Image.Entity.ImageID,
				Valid: true,
			}
		}
		matchIndex := contains(existKey, a.Key)
		if matchIndex == -1 {
			as = append(as, parameter.CreateGradeServiceParam{
				Name:         a.Name,
				Key:          a.Key,
				Description:  entity.String{String: a.Name, Valid: true},
				Color:        entity.String{String: a.Color, Valid: true},
				CoverImageID: imageID,
			})
		} else {
			var uid uuid.UUID
			existIDs, uid = removeUUID(existIDs, matchIndex)
			existKey, _ = removeString(existKey, matchIndex)
			if deepEqual {
				de := isDeepEqual(a, existData[uid])
				if !de {
					_, err = (*c.Manager).UpdateGrade(
						ctx,
						uid,
						a.Name,
						entity.String{String: a.Description, Valid: true},
						entity.String{String: a.Color, Valid: true},
						imageID,
					)
					if err != nil {
						return fmt.Errorf("failed to update grade: %w", err)
					}
				}
			}
			delete(existData, uid)
		}
	}
	_, err = (*c.Manager).CreateGrades(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create grades: %w", err)
	}
	if len(existIDs) > 0 && !notDel {
		_, err = (*c.Manager).PluralDeleteGrades(ctx, existIDs)
		if err != nil {
			return fmt.Errorf("failed to delete grades: %w", err)
		}
	}
	return nil
}
