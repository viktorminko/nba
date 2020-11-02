package publisher

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/viktorminko/nba/pkg/simulation/transport"
	"sync"
)

//Start starts a publisher to publish events coming from simulation to the queue
//represented by Transporter interface. It works as a proxy between simulation and queue.
//Message is read form dataCh and send to the queue.
//It returns a channel where it sends errors occurred while processing messages.
func Start(ctx context.Context, wg *sync.WaitGroup, dataCh <-chan []byte, transport transport.Transporter) (<-chan error, error) {
	errCh := make(chan error)

	if transport == nil {
		return nil, errors.New("transport is nil")
	}

	wg.Add(1)
	go func() {
		defer func() {
			close(errCh)
			wg.Done()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case b, ok := <-dataCh:
				if !ok {
					return
				}

				if err := transport.Transport(bytes.NewReader(b)); err != nil {
					errCh <- errors.Wrap(err, "transport message")
				}
			}
		}

	}()

	return errCh, nil
}
