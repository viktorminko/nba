package frontend

import (
	"github.com/viktorminko/nba/pkg/statistic/stats"
	"io"
)

//Displayer interface specifies method to write statistics data to io.Writer
type Displayer interface {
	Display(w io.Writer, st *stats.Stats) error
}
