package service

import (
	"context"
	"github.com/pkg/errors"
	"github.com/viktorminko/nba/pkg/simulation/game"
	"github.com/viktorminko/nba/pkg/simulation/publisher"
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

//start goroutine to handle errors from channel
//goroutine finishes when channel is closed
func startErrorHandler(errCh <-chan error) {
	go func() {
		for err := range errCh {
			log.Println(err)
		}
	}()
}

func Start(ctx context.Context, initData io.Reader, eventsTopic transport.Transporter, gameDuration, eventDuration time.Duration) error {
	log.Println("Starting service")

	//read teams data from reader and prepare random games
	games, err := initSimulation(ctx, initData)
	if err != nil {
		return errors.Wrap(err, "init simulation")
	}

	var wg sync.WaitGroup

	//start simulation
	eventsCh, simulationErrch := simulation.Start(ctx, &wg, games, gameDuration, eventDuration)
	startErrorHandler(simulationErrch)

	//subscribe on events from simulation and start queue publisher
	errCh, err := publisher.Start(ctx, &wg, eventsCh, eventsTopic)
	if err != nil {
		return errors.Wrap(err, "start goals publisher")
	}
	startErrorHandler(errCh)

	wg.Wait()

	return nil
}
