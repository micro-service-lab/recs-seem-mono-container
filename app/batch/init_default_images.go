package batch

import (
	"context"
	"fmt"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
)

// InitDefaultImages is a batch to initialize default images.
type InitDefaultImages struct {
	Manager *service.ManagerInterface
}

// Run initializes images.
func (c *InitDefaultImages) Run(ctx context.Context) error {
	entries, err := service.DefaultImages.ReadDir("static/images")
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}
	var sp []parameter.CreateImageSpecifyFilenameServiceParam
	for i, e := range entries {
		if e.IsDir() {
			continue
		}
		f, err := service.DefaultImages.Open("static/images/" + e.Name())
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()
		sp = append(sp, parameter.CreateImageSpecifyFilenameServiceParam{
			Origin:   f,
			Filename: service.DefaultImageKeys[i],
			Alias:    e.Name(),
		})
	}
	_, err = (*c.Manager).CreateImagesSpecifyFilename(ctx, entity.UUID{}, sp)
	if err != nil {
		return fmt.Errorf("failed to create images: %w", err)
	}
	return nil
}

// RunDiff run only if there is a difference.
func (c *InitDefaultImages) RunDiff(_ context.Context, _, _ bool) error {
	fmt.Println("Init Default Images does not support diff.")
	return nil
}
