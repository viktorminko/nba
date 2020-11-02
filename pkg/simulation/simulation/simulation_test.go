package simulation

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
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
	eventsCh, _ := Start(
		context.Background(),
		&wg,
		[]*game.Game{
			{
				Home: &game.Team{
					ID: "Chicago",
				},

				Guest: &game.Team{
					ID: "LA	",
				},
			},
		},
		gameDuration,
		eventDuration,
	)

	var (
		goals        []event.Goal
		gameStatuses []event.GameStatus
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case msg, ok := <-eventsCh:
				if !ok {
					return
				}

				var goal event.Goal
				err := gob.NewDecoder(bytes.NewBuffer(msg)).Decode(&goal)
				if err == nil {
					goals = append(goals, goal)
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
		}
	}()

	wg.Wait()
	assert.NotEmpty(t, goals)
	assert.NotEmpty(t, gameStatuses)

	fmt.Println(goals)
	fmt.Println(gameStatuses)
}
