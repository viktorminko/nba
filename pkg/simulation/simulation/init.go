package simulation

import (
	"encoding/json"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/viktorminko/nba/pkg/simulation/game"
	"io"
	"io/ioutil"
	"math/rand"
	"time"
)

type playerWithTeam struct {
	Name string `json:"name"`
	Team string `json:"team"`
}

//Init reads team info from reader and prepares games
//by randomizing team pairs
func Init(r io.Reader) ([]*game.Game, error) {
	playersWithTeams, err := readTeams(r)
	if err != nil {
		return nil, errors.Wrap(err, "read teams")
	}

	return buildPairs(buildTeams(playersWithTeams)), nil

}

func buildTeams(playersWithTeams []playerWithTeam) []*game.Team {
	teams := make(map[string]*game.Team)
	for _, pl := range playersWithTeams {
		_, ok := teams[pl.Team]

		if !ok {
			teams[pl.Team] = &game.Team{
				ID:   uuid.NewV4().String(),
				Name: pl.Team,
				Players: []*game.Player{
					{
						ID:   uuid.NewV4().String(),
						Name: pl.Name,
					},
				},
			}
			continue
		}
	}

	res := make([]*game.Team, 0, len(teams))
	for _, t := range teams {
		res = append(res, t)
	}

	return res
}

func buildPairs(teams []*game.Team) []*game.Game {
	//shuffle teams
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(teams), func(i, j int) { teams[i], teams[j] = teams[j], teams[i] })

	var res []*game.Game
	for i := 0; i+1 < len(teams); i += 2 {
		res = append(res, &game.Game{
			ID:    uuid.NewV4().String(),
			Home:  teams[i],
			Guest: teams[i+1],
		})
	}

	return res
}

func readTeams(r io.Reader) ([]playerWithTeam, error) {
	bts, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "read data")
	}

	var data []playerWithTeam
	if err := json.Unmarshal(bts, &data); err != nil {
		return nil, errors.Wrap(err, "unmarshall json")
	}

	return data, nil
}
