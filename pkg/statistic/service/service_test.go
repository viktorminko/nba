package service

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viktorminko/nba/pkg/simulation/event"
	"github.com/viktorminko/nba/pkg/simulation/game"
	"github.com/viktorminko/nba/pkg/statistic/stats"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type mockSubscriber struct {
	dataCh chan []byte
}

func (s *mockSubscriber) Subscribe() <-chan []byte {
	return s.dataCh
}

type mockDisplayer func(w io.Writer, st *stats.Stats) error

func (fn mockDisplayer) Display(w io.Writer, st *stats.Stats) error {
	return fn(w, st)
}

func TestStart(t *testing.T) {
	eventSub := &mockSubscriber{
		dataCh: make(chan []byte),
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		assert.NoError(t, Start(
			ctx,
			eventSub,
			8080,
			mockDisplayer(func(w io.Writer, st *stats.Stats) error {
				exp := stats.New()
				exp.TotalGuest = 3
				exp.Games = map[string]*stats.Game{
					"game1": {
						ID: "game1",
						Home: &stats.TeamScore{
							Team: &game.Team{
								ID: "home",
							},
							Score: 0,
						},
						Guest: &stats.TeamScore{
							Team: &game.Team{
								ID: "guest",
							},
							Score: 3,
						},
						Status:              event.GameStatusStarted,
						LastEventSinceStart: st.Games["game1"].LastEventSinceStart,
					},
				}

				assert.Equal(t, exp, st)
				_, err := w.Write([]byte(fmt.Sprintf("%#v", st)))
				assert.NoError(t, err)
				return nil
			}),
		))
	}()

	var bts bytes.Buffer

	//Send goal message
	err := gob.NewEncoder(&bts).Encode(event.Goal{
		GameID:       "game1",
		TeamScoredID: "guest",
		HomeTeamID:   "home",
		GuestTeamID:  "guest",
		Value:        3,
	})
	assert.NoError(t, err)

	eventSub.dataCh <- bts.Bytes()
	time.Sleep(100 * time.Millisecond)

	//Send gameStatus message
	var bts1 bytes.Buffer
	err = gob.NewEncoder(&bts1).Encode(event.GameStatus{
		Game: &game.Game{
			ID: "game1",
			Home: &game.Team{
				ID: "home",
			},
			Guest: &game.Team{
				ID: "guest",
			},
		},
		Status: event.GameStatusStarted,
	})
	assert.NoError(t, err)

	eventSub.dataCh <- bts1.Bytes()
	time.Sleep(100 * time.Millisecond)

	res, err := http.Get("http://0.0.0.0:8080/")
	assert.NoError(t, err)

	bt, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	defer func() {
		cerr := res.Body.Close()
		assert.NoError(t, cerr)
	}()

	assert.NotEmpty(t, bt)

	cancel()
	time.Sleep(100 * time.Millisecond)
	_, err = http.Get("http://0.0.0.0:8080/")
	assert.Error(t, err)
}
