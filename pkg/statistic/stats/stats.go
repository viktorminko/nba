package stats

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/pkg/errors"
	"github.com/viktorminko/nba/pkg/simulation/event"
	"github.com/viktorminko/nba/pkg/simulation/game"
	"log"
	"sync"
	"time"
)

//TeamScore represents score of the team
type TeamScore struct {
	Team  *game.Team
	Score int
}

//Game represents statistics for one game
type Game struct {
	//Game ID
	ID string
	//Home team and its score
	Home *TeamScore
	//Guest team and its score
	Guest *TeamScore
	//Game status e.g. Finished
	Status event.Status
	//When game started
	TimeStarted time.Time
	//When game finished
	TimeFinished time.Time
	//When last event in the game occurred, contains amount of time
	//spend before game start and event occurred
	LastEventSinceStart time.Duration
}

//Statistics
type Stats struct {
	sync.Mutex
	//Map of the games in current statistic
	Games map[string]*Game
	//Total scores of home and guest teams
	TotalHome, TotalGuest int
}

//New returns new Statistic
func New() *Stats {
	return &Stats{
		Games: make(map[string]*Game),
	}
}

//Update statistics on Goal event
//safe for concurrent calls
func (s *Stats) handleGoalEvent(e *event.Goal) {
	s.Lock()
	defer func() {
		s.Unlock()
		log.Println("stats updated with goal")
	}()
	gameID := e.GameID
	if _, ok := s.Games[gameID]; !ok {
		s.Games[gameID] = &Game{
			ID: gameID,
			Home: &TeamScore{
				Team: &game.Team{
					ID: e.HomeTeamID,
				},
			},
			Guest: &TeamScore{
				Team: &game.Team{
					ID: e.GuestTeamID,
				},
			},
		}
	}

	cur := s.Games[gameID]
	if cur.Home.Team.ID == e.TeamScoredID {
		cur.Home.Score += e.Value
		s.TotalHome += e.Value
	} else {
		cur.Guest.Score += e.Value
		s.TotalGuest += e.Value
	}

	cur.LastEventSinceStart = time.Since(cur.TimeStarted)
}

//update statistics on GameStatus event
//save for concurrent call
func (s *Stats) handleGameStatusEvent(e event.GameStatus) {
	s.Lock()
	defer func() {
		s.Unlock()
		log.Println("stats updated with game status")
	}()

	gameID := e.Game.ID
	if _, ok := s.Games[gameID]; !ok {
		s.Games[gameID] = &Game{
			ID: gameID,
			Home: &TeamScore{
				Team: e.Game.Home,
			},
			Guest: &TeamScore{
				Team: e.Game.Guest,
			},
			TimeStarted: time.Now(),
		}
	}

	cur := s.Games[gameID]
	cur.Home.Team = e.Game.Home
	cur.Guest.Team = e.Game.Guest
	cur.Status = e.Status
	cur.LastEventSinceStart = time.Since(cur.TimeStarted)

	if e.Status == event.GameStatusFinished {
		cur.TimeFinished = time.Now()
	}
}

//StartUpdater starts statistic updater loop.
//It reads data from ch, tries to decode message to known event and update statistic.
//Return channel to read errors happened in current goroutine
func (s *Stats) StartUpdater(ctx context.Context, ch <-chan []byte) <-chan error {
	log.Println("start goal updater")
	errCh := make(chan error)
	go func() {
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}

				var err error

				var goalEvent event.Goal
				if err = gob.NewDecoder(bytes.NewBuffer(msg)).Decode(&goalEvent); err == nil {
					log.Println("goal message received", goalEvent)
					s.handleGoalEvent(&goalEvent)
					continue
				}

				var gameStatusEvent event.GameStatus
				if err = gob.NewDecoder(bytes.NewBuffer(msg)).Decode(&gameStatusEvent); err == nil {
					log.Println("game status message received", gameStatusEvent)
					s.handleGameStatusEvent(gameStatusEvent)
					continue
				}

				errCh <- errors.Wrap(err, "unable to decode message")
			}
		}

	}()

	return errCh
}
