package core

import (
	"context"

	"github.com/pkg/errors"
)

func (r *mutationResolver) ConnectToSerial(ctx context.Context, void *bool) (*bool, error) {
	return nil, r.App.PrinterManager.ConnectToSerial()
}

func (r *mutationResolver) DisconnectFromSerial(ctx context.Context, void *bool) (*bool, error) {
	return nil, r.App.PrinterManager.DisconnectFromSerial()
}

func (r *subscriptionResolver) SerialConsoleData(ctx context.Context) (<-chan string, error) {
	outChan := make(chan string)
	sub, can := r.App.PrinterManager.consoleBroadcaster.Subscribe()
	go func() {
	Loop:
		for {
			select {
			case data := <-sub:
				if data != nil {
					outChan <- data.(string)
				}
			case <-ctx.Done():
				can()
				break Loop
			}
		}
	}()

	return outChan, nil
}

func (r *queryResolver) RecentCommands(ctx context.Context) ([]string, error) {
	return r.App.RecentCommandsManager.GetRecentCommands(), nil
}

func (r *mutationResolver) SendGcode(ctx context.Context, cmd string) (*bool, error) {
	r.App.PrinterManager.consoleWriteSem <- cmd + "\r\n"
	return nil, nil
}

func (r *mutationResolver) SendConsoleCommand(ctx context.Context, cmd string) (*bool, error) {
	r.App.PrinterManager.consoleWriteSem <- cmd + "\r\n"
	return nil, r.App.RecentCommandsManager.AddRecentCommand(cmd)
}

func (r *queryResolver) ScrollbackBuffer(ctx context.Context) (string, error) {
	r.App.PrinterManager.scrollbackBufferMutex.Lock()
	defer r.App.PrinterManager.scrollbackBufferMutex.Unlock()
	if r.App.PrinterManager.scrollbackBuffer == nil {
		return "", errors.New("No scrollback buffer")
	}
	return r.App.PrinterManager.scrollbackBuffer.String(), nil
}