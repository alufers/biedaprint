package core

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

func (r *mutationResolver) StartPrintJob(ctx context.Context, id int) (*bool, error) {
	meta, err := r.App.GcodeFileMetaRepositoryService.GetOneByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to load gcode meta: %w", err)
	}
	job := &PrintJobInternal{
		app: r.App,
		PrintJob: &PrintJob{
			GcodeMeta:   meta,
			StartedTime: time.Now(),
		},
	}

	select {
	case r.App.PrinterService.printJobSem <- job:
	default:
		return nil, errors.New("serial writer busy with antoher job")
	}
	return nil, nil
}

func (r *mutationResolver) AbortPrintJob(ctx context.Context, void *bool) (*bool, error) {
	select {
	case r.App.PrinterService.abortPrintSem <- true:
	default:
	}
	return nil, nil
}

func (r *queryResolver) CurrentPrintJob(ctx context.Context) (*PrintJob, error) {
	panic("not implemented")
}

func (r *subscriptionResolver) CurrentPrintJobUpdated(ctx context.Context) (<-chan *PrintJob, error) {
	panic("not implemented")
}
