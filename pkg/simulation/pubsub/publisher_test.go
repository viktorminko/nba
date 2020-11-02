package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/viktorminko/nba/pkg/simulation/event"
	"github.com/viktorminko/nba/pkg/simulation/transport"
	"io"
	"io/ioutil"
	"sync"
	"testing"
	"time"
)

type mockTrans func(r io.Reader) error

func (tr mockTrans) Transport(r io.Reader) error {
	return tr(r)
}

func TestStart(t *testing.T) {
	for _, tt := range []struct {
		name       string
		getCtx     func() (context.Context, context.CancelFunc)
		getDataCh  func() <-chan []byte
		trans      transport.Transporter
		expError   bool
		errHandler func(errCh <-chan error)
	}{
		{
			name:     "nil transport",
			trans:    nil,
			expError: true,
		},

		{
			name: "transport error",
			getCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			getDataCh: func() <-chan []byte {
				ch := make(chan []byte)
				go func() {
					b, err := json.Marshal(event.Goal{})
					assert.NoError(t, err)
					ch <- b
				}()
				time.Sleep(100 * time.Millisecond)
				return ch
			},
			trans: mockTrans(func(r io.Reader) error {
				return errors.New("transport error")
			}),
			expError: false,
			errHandler: func(errCh <-chan error) {
				assert.Contains(t, (<-errCh).Error(), "transport error")
			},
		},

		{
			name: "valid call",
			getCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			getDataCh: func() <-chan []byte {
				ch := make(chan []byte)
				go func() {
					b, err := json.Marshal(event.Goal{Value: 4})
					assert.NoError(t, err)
					ch <- b
				}()
				time.Sleep(100 * time.Millisecond)
				return ch
			},
			trans: mockTrans(func(r io.Reader) error {
				b, err := ioutil.ReadAll(r)
				assert.NoError(t, err)

				g := event.Goal{}
				assert.NoError(t, json.Unmarshal(b, &g))

				assert.Equal(t, event.Goal{Value: 4}, g)
				return nil
			}),
			expError: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var dataCh <-chan []byte
			if tt.getDataCh != nil {
				dataCh = tt.getDataCh()
			}

			var ctx context.Context
			var cancel context.CancelFunc

			if tt.getCtx != nil {
				ctx, cancel = tt.getCtx()
				defer cancel()
			}

			var wg sync.WaitGroup
			errCh, err := StartPub(ctx, &wg, dataCh, tt.trans)
			if tt.expError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			if tt.errHandler != nil {
				tt.errHandler(errCh)
			}

		})
	}
}
