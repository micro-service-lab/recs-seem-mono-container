package batch

import (
	"context"
	"errors"
	"fmt"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/app/storage"
)

// InitWholeOrganization is a batch to initialize grades.
type InitWholeOrganization struct {
	Manager *service.ManagerInterface
	Storage storage.Storage
}

// Run initializes grades.
func (c *InitWholeOrganization) Run(ctx context.Context) error {
	shuffleImages := service.DefaultImageKeys
	service.ArrayShuffle(shuffleImages)
	url, err := c.Storage.GetURLFromKey(ctx, shuffleImages[0])
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
	_, err = (*c.Manager).CreateWholeOrganization(
		ctx,
		service.WholeOrganization.Name,
		entity.String{String: service.WholeOrganization.Description, Valid: true},
		entity.String{String: service.WholeOrganization.Color, Valid: true},
		imageID,
	)
	if err != nil {
		return fmt.Errorf("failed to create grades: %w", err)
	}
	return nil
}

// RunDiff run only if there is a difference.
func (c *InitWholeOrganization) RunDiff(ctx context.Context, _, deepEqual bool) error {
	exist, err := (*c.Manager).FindWholeOrganization(ctx)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if !errors.As(err, &e) {
			return fmt.Errorf("failed to find grades: %w", err)
		}
		return c.Run(ctx)
	}
	shuffleImages := service.DefaultImageKeys
	service.ArrayShuffle(shuffleImages)
	existData := service.Organization{
		Name:        exist.Name,
		Description: exist.Description.String,
		Color:       exist.Color.String,
	}
	url, err := c.Storage.GetURLFromKey(ctx, shuffleImages[0])
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
	if deepEqual {
		de := isDeepEqual(service.WholeOrganization, existData)
		if !de {
			_, err = (*c.Manager).UpdateWholeOrganization(
				ctx,
				service.WholeOrganization.Name,
				entity.String{String: service.WholeOrganization.Description, Valid: true},
				entity.String{String: service.WholeOrganization.Color, Valid: true},
				imageID,
			)
			if err != nil {
				return fmt.Errorf("failed to update grade: %w", err)
			}
		}
	}
	return nil
}
