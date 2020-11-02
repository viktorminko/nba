package event

import "github.com/viktorminko/nba/pkg/simulation/game"

type Status int

const (
	GameStatusStarted  = 1
	GameStatusFinished = 2
)

//Status of current game e.g. started/finished
type GameStatus struct {
	Game   *game.Game
	Status Status
}
