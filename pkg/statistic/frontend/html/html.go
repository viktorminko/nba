package frontend

import (
	"github.com/pkg/errors"
	"github.com/viktorminko/nba/pkg/statistic/stats"
	"html/template"
	"io"
	"sort"
	"strconv"
	"time"
)

//ViewFunc is a wrapper that implements Displayer interface
type ViewFunc func(w io.Writer, st *stats.Stats) error

func (v ViewFunc) Display(w io.Writer, st *stats.Stats) error {
	return v(w, st)
}

type game struct {
	Home                string
	Guest               string
	Score               string
	TimeStarted         string
	TimeFinished        string
	LastEventSinceStart string
}

type viewData struct {
	Games      []game
	TotalScore string
	Debug      string
}

//New returns a function to render statistic into html using
//template in statsTemplatePath
func New(statsTemplatePath string) (ViewFunc, error) {
	tpl, err := template.ParseFiles(statsTemplatePath)
	if err != nil {
		return nil, errors.Wrap(err, "parse template file")
	}

	return ViewFunc(func(w io.Writer, st *stats.Stats) error {
		var data viewData

		for _, v := range st.Games {
			var started, finished, lastEvent string
			if !v.TimeStarted.IsZero() {
				started = v.TimeStarted.Format(time.StampMilli)
				lastEvent = v.LastEventSinceStart.Round(time.Millisecond).String()
			}
			if !v.TimeFinished.IsZero() {
				finished = v.TimeFinished.Format(time.StampMilli)
			}
			data.Games = append(data.Games, game{
				Home:                v.Home.Team.Name,
				Guest:               v.Guest.Team.Name,
				Score:               strconv.Itoa(v.Home.Score) + " : " + strconv.Itoa(v.Guest.Score),
				TimeStarted:         started,
				TimeFinished:        finished,
				LastEventSinceStart: lastEvent,
			})
		}

		data.TotalScore = strconv.Itoa(st.TotalHome) + " : " + strconv.Itoa(st.TotalGuest)

		//sort by home team
		sort.Slice(data.Games, func(i, j int) bool {
			return data.Games[i].Home < data.Games[j].Home
		})

		return tpl.Execute(w, data)
	}), nil
}
