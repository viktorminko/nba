package service

import (
	"context"
	"github.com/pkg/errors"
	"github.com/viktorminko/nba/pkg/simulation/game"
	"github.com/viktorminko/nba/pkg/simulation/pubsub"
	"github.com/viktorminko/nba/pkg/simulation/simulation"
	"github.com/viktorminko/nba/pkg/simulation/transport"
	"io"
	"log"
	"sync"
	"time"
)

func initSimulation(ctx context.Context, r io.Reader) ([]*game.Game, error) {
	games, err := simulation.Init(r)
	if err != nil {
		return nil, errors.Wrap(err, "init games")
	}

	return games, nil
}

func startErrorHandler(errCh <-chan error) {
	go func() {
		for {
			select {
			case err, ok := <-errCh:
				if !ok {
					return
				}
				log.Println(err)
			}
		}
	}()
}

func Start(ctx context.Context, initData io.Reader, eventsTopic transport.Transporter, gameDuration, eventDuration time.Duration) error {
	log.Println("Starting service")

	//randomize team pairs
	games, err := initSimulation(ctx, initData)
	if err != nil {
		return errors.Wrap(err, "init simulation")
	}

	var wg sync.WaitGroup

	eventsCh, simulationErrch := simulation.Start(ctx, &wg, games, gameDuration, eventDuration)
	startErrorHandler(simulationErrch)

	errCh, err := pubsub.StartPub(ctx, &wg, eventsCh, eventsTopic)
	if err != nil {
		return errors.Wrap(err, "start goals publisher")
	}
	startErrorHandler(errCh)

	wg.Wait()

	return nil
}
