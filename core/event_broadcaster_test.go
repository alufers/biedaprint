package core

import (
	"testing"
	"time"
)

func TestBroadcasts(t *testing.T) {
	eb := NewEventBroadcaster()
	lis := make(chan bool)
	go func() {
		ev, _ := eb.Subscribe()
		lis <- true

		select {
		case data := <-ev:
			if data != "dupa" {
				t.Fatalf("invalid data ddd")
			}
			lis <- true

		case <-time.After(time.Millisecond * 100):
			t.Fatalf("didn't transmit")
			lis <- true
		}
	}()
	<-lis
	eb.Broadcast("dupa")
	<-lis
}

func TestCancels(t *testing.T) {
	eb := NewEventBroadcaster()
	lis := make(chan bool)
	go func() {
		ev, can := eb.Subscribe()
		can()
		lis <- true

		select {
		case _, ok := <-ev:
			if ok {
				panic("bad")
			}
			lis <- true
		case <-time.After(time.Millisecond * 100):

			lis <- true
		}
	}()
	<-lis
	<-time.After(time.Millisecond * 100)
	eb.Broadcast("dupa")
	<-lis
}
