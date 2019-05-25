package core

import (
	"sync"
	"time"
)

type TrackedValueInternal struct {
	*TrackedValue
	ValueMutex              *sync.RWMutex
	ValueUpdatedBroadcaster *EventBroadcaster
}

func NewTrackedValueInternal(inner *TrackedValue) *TrackedValueInternal {
	return &TrackedValueInternal{
		TrackedValue:            inner,
		ValueMutex:              &sync.RWMutex{},
		ValueUpdatedBroadcaster: NewEventBroadcaster(),
	}
}

func (tv *TrackedValueInternal) UpdateValue(val interface{}) {

	tv.ValueMutex.Lock()
	defer tv.ValueMutex.Unlock()
	if tv.MaxHistoryLength != 0 {
		if len(tv.History) >= tv.MaxHistoryLength {
			tv.History = append(tv.History[1:], val)
		} else {
			tv.History = append(tv.History, val)
		}
	}
	tv.Value = val
	now := time.Now()
	tv.LastUpdate = &now
	if tv.LastSent == nil || time.Now().Sub(*tv.LastSent) > time.Millisecond*time.Duration(tv.MinUpdateInterval) {

		tv.LastSent = &now
		tv.ValueUpdatedBroadcaster.Broadcast(val)
	}

}
