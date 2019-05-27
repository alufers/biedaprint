package core

import (
	"context"

	"github.com/pkg/errors"
)

func (r *subscriptionResolver) TrackedValueUpdated(ctx context.Context, name string) (<-chan interface{}, error) {
	tv, ok := r.App.TrackedValuesService.TrackedValues[name]
	tv.ValueMutex.RLock()
	defer tv.ValueMutex.RUnlock()
	if !ok {
		return nil, errors.New("tracked value with this name not found")
	}
	outChan := make(chan interface{})
	sub, can := tv.ValueUpdatedBroadcaster.Subscribe()
	go func() {
	Loop:
		for {
			select {
			case data := <-sub:
				outChan <- data
			case <-ctx.Done():
				can()
				break Loop
			}
		}
	}()

	return outChan, nil
}

func (r *queryResolver) TrackedValues(ctx context.Context) (resp []*TrackedValue, err error) {
	resp = []*TrackedValue{}
	for _, tv := range r.App.TrackedValuesService.TrackedValues {
		tv.ValueMutex.RLock()
		val := *tv.TrackedValue
		resp = append(resp, &val)
		tv.ValueMutex.RUnlock()
	}
	return resp, nil
}

func (r *queryResolver) TrackedValue(ctx context.Context, name string) (*TrackedValue, error) {
	tv, ok := r.App.TrackedValuesService.TrackedValues[name]
	tv.ValueMutex.RLock()
	defer tv.ValueMutex.RUnlock()
	if !ok {
		return nil, errors.New("tracked value with this name not found")
	}
	val := *tv.TrackedValue
	return &val, nil
}
