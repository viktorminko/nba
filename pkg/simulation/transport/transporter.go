package transport

import (
	"io"
)

//Transporter interface specifies interface for transporting data from reader
type Transporter interface {
	Transport(r io.Reader) error
}
