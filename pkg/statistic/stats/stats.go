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

type TeamScore struct {
	Team  *game.Team
	Score int
}

type Game struct {
	ID                  string
	Home                *TeamScore
	Guest               *TeamScore
	Status              event.Status
	TimeStarted         time.Time
	TimeFinished        time.Time
	LastEventSinceStart time.Duration
}

type Stats struct {
	sync.Mutex
	Games                 map[string]*Game
	TotalHome, TotalGuest int
}

func New() *Stats {
	return &Stats{
		Games: make(map[string]*Game),
	}
}

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

func (s *Stats) handleGameStatusEvent(e event.GameStatus) {
	s.Lock()
	defer func() {
		s.Unlock()
		log.Println("stats updated with game status")
	}()

	log.Printf("received game status change event: %s - %s : %d", e.Game.Home.Name, e.Game.Guest.Name, e.Status)

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

				var goalEvent event.Goal
				if err := gob.NewDecoder(bytes.NewBuffer(msg)).Decode(&goalEvent); err == nil {
					log.Println("goal message received", goalEvent)
					s.handleGoalEvent(&goalEvent)
					continue
				}

				var gameStatusEvent event.GameStatus
				if err := gob.NewDecoder(bytes.NewBuffer(msg)).Decode(&gameStatusEvent); err == nil {
					log.Println("game status message received", gameStatusEvent)
					s.handleGameStatusEvent(gameStatusEvent)
					continue
				}

				errCh <- errors.New("unable to deserialize message")
			}
		}

	}()

	return errCh
}
