package profiler

import (
	"context"
	"fmt"

	"cloud.google.com/go/profiler"
	"github.com/hsmtkk/bookish-pancake/utilgcp"
)

func Start(ctx context.Context) error {
	projectID, err := utilgcp.ProjectID(ctx)
	if err != nil {
		return err
	}
	if err := profiler.Start(profiler.Config{ProjectID: projectID}); err != nil {
		return fmt.Errorf("profiler.Start failed; %w", err)
	}
	return nil
}
