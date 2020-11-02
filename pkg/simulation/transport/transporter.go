package transport

import (
	"io"
)

type Transporter interface {
	Transport(r io.Reader) error
}
