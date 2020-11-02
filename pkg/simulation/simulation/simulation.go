package simulation

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/pkg/errors"
	"github.com/viktorminko/nba/pkg/simulation/event"
	"github.com/viktorminko/nba/pkg/simulation/game"
	"log"
	"math/rand"
	"sync"
	"time"
)

func serialize(v interface{}) ([]byte, error) {
	var bts bytes.Buffer
	enc := gob.NewEncoder(&bts)
	if err := enc.Encode(v); err != nil {
		return nil, errors.Wrap(err, "encode data")
	}

	return bts.Bytes(), nil
}

func validateGame(g *game.Game) error {
	if g == nil {
		return errors.New("game is nil")
	}

	if len(g.ID) == 0 {
		return errors.New("game id is empty")
	}

	if g.Guest == nil {
		return errors.New("game guest team is nil")
	}

	if g.Home == nil {
		return errors.New("game home team is nil")
	}

	//@TODO add more validation here
	return nil
}

func startGame(
	ctx context.Context,
	wg *sync.WaitGroup,
	curGame *game.Game,
	gameDuration, eventDuration time.Duration,
	errCh chan<- error,
	eventsCh chan<- []byte,
	finishCh chan<- struct{}) {

	//indicate that game was finished
	defer func() {
		finishCh <- struct{}{}
		wg.Done()
	}()

	//if game is invalid, return immediately
	if err := validateGame(curGame); err != nil {
		errCh <- errors.Wrap(err, "validate game")
		return
	}

	log.Println("game started")
	defer log.Println("game finished")

	defer func() {
		b, err := serialize(event.GameStatus{
			Game:   curGame,
			Status: event.GameStatusFinished,
		})
		if err != nil {
			errCh <- errors.Wrap(err, "serialize game status event")
		}

		eventsCh <- b
	}()

	bts, err := serialize(event.GameStatus{
		Game:   curGame,
		Status: event.GameStatusStarted,
	})
	if err != nil {
		errCh <- errors.Wrap(err, "serialize game status event")
	}
	eventsCh <- bts

	eventTicker := time.NewTicker(eventDuration)
	defer eventTicker.Stop()

	timeAfter := time.After(gameDuration)
	rand.Seed(time.Now().UnixNano())
	for {
		select {
		case <-ctx.Done():
			return
		case <-timeAfter:
			//finish game
			return
		case <-eventTicker.C:
			//either one of teams scored, or neither scored
			//nolint
			r := rand.Intn(3)

			if r == 0 {
				break
			}

			curGoal := event.Goal{
				GameID:       curGame.ID,
				TeamScoredID: curGame.Home.ID,
				HomeTeamID:   curGame.Home.ID,
				GuestTeamID:  curGame.Guest.ID,
			}

			if r == 2 {
				curGoal.TeamScoredID = curGame.Guest.ID
			}

			//nolint gosec
			curGoal.Value = rand.Intn(2) + 2

			b, err := serialize(curGoal)
			if err != nil {
				errCh <- errors.Wrap(err, "serialize goal event")
			}

			eventsCh <- b
		}
	}
}

//Start starts main simulation loop for provided games, duration of the game and duration of events being sent.
//It returns channel to read events from and channel to read errors happened in simulation.
func Start(ctx context.Context, wg *sync.WaitGroup, games []*game.Game, gameDuration, eventDuration time.Duration) (
	<-chan []byte,
	<-chan error) {
	log.Println("simulation started")

	eventsCh := make(chan []byte)
	errCh := make(chan error)
	finishCh := make(chan struct{})

	for i := range games {
		select {
		case <-ctx.Done():
			return nil, nil
		default:

		}

		wg.Add(1)
		go startGame(ctx, wg, games[i], gameDuration, eventDuration, errCh, eventsCh, finishCh)
	}

	//we need to wait until all games are finished in a non blocking manner
	//when games are finished we need to close all corresponded channels
	go func() {
		defer func() {
			close(eventsCh)
			close(errCh)
			close(finishCh)
		}()

		for i := 0; i < len(games); i++ {
			select {
			case <-ctx.Done():
				return
			case <-finishCh:

			}
		}
		log.Println("simulation finished")
	}()

	return eventsCh, errCh
}
