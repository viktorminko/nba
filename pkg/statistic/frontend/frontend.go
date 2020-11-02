package frontend

import (
	"github.com/viktorminko/nba/pkg/statistic/stats"
	"io"
)

type Displayer interface {
	Display(w io.Writer, st *stats.Stats) error
}
