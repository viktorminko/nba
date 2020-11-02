package simulation

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/stretchr/testify/assert"
	"github.com/viktorminko/nba/pkg/simulation/event"
	"github.com/viktorminko/nba/pkg/simulation/game"
	"sync"
	"testing"
	"time"
)

func TestStart(t *testing.T) {

	gameDuration := 6 * time.Second
	eventDuration := 2 * time.Second

	var wg sync.WaitGroup
	eventsCh, errCh := Start(
		context.Background(),
		&wg,
		[]*game.Game{
			//valid game
			{
				ID: "game1",
				Home: &game.Team{
					ID: "Chicago",
				},

				Guest: &game.Team{
					ID: "LA	",
				},
			},
			//invalid game
			nil,
			//valid game
			{
				ID: "game2",
				Home: &game.Team{
					ID: "Bulls",
				},

				Guest: &game.Team{
					ID: "Miami	",
				},
			},
			//invalid game empty ID
			{},
			//invalid game, empty home team
			{
				ID: "game3",
			},
			//invalid game, empty guest team
			{
				ID: "game4",
				Home: &game.Team{
					ID: "NY",
				},
			},
			//invalid game
			nil,
		},
		gameDuration,
		eventDuration,
	)

	var errCount int
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range errCh {
			errCount++
		}
	}()

	var (
		gameStatuses []event.GameStatus
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range eventsCh {
			var goal event.Goal
			err := gob.NewDecoder(bytes.NewBuffer(msg)).Decode(&goal)
			if err == nil {
				continue
			}

			var statusChange event.GameStatus
			err = gob.NewDecoder(bytes.NewBuffer(msg)).Decode(&statusChange)
			if err == nil {
				gameStatuses = append(gameStatuses, statusChange)
				continue
			}

			t.Error("unknown message type")
		}
	}()

	wg.Wait()

	//invalid games
	assert.Equal(t, 5, errCount)

	//goals might be empty, but statuses might not
	//each valid game have start and finish events
	assert.Equal(t, 4, len(gameStatuses))
}
