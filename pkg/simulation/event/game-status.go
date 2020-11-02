package event

import "github.com/viktorminko/nba/pkg/simulation/game"

type Status int

const (
	GameStatusStarted  = 1
	GameStatusFinished = 2
)

type GameStatus struct {
	Game   *game.Game
	Status Status
}
