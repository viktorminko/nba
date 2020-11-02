package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

type mockTransport struct {
	data []string
}

func (t *mockTransport) Transport(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	t.data = append(t.data, string(b))
	return nil
}

func TestStart(t *testing.T) {
	eventsTr := &mockTransport{}
	assert.NoError(t, Start(
		context.Background(),
		//nolint
		strings.NewReader("[\n  {\n    \"name\": \"Sun Yue\",\n    \"team\": \"Lakers\",\n    \"points\": 6\n  },\n  {\n    \"name\": \"Arron Afflalo\",\n    \"team\": \"Rockets\",\n    \"points\": 724\n  },\n  {\n    \"name\": \"Alexis Ajinca\",\n    \"team\": \"Bobcats\",\n    \"points\": 10\n  },\n  {\n    \"name\": \"LaMarcus Aldridge\",\n    \"team\": \"Trailblazers\",\n    \"points\": 1393\n  },\n  {\n    \"name\": \"Joe Alexander\",\n    \"team\": \"Bulls\",\n    \"points\": 4\n  },\n  {\n    \"name\": \"Malik Allen\",\n    \"team\": \"Rockets\",\n    \"points\": 105\n  },\n  {\n    \"name\": \"Ray Allen\",\n    \"team\": \"Celtics\",\n    \"points\": 1304\n  },\n  {\n    \"name\": \"Tony Allen\",\n    \"team\": \"Celtics\",\n    \"points\": 330\n  },\n  {\n    \"name\": \"Rafer Alston\",\n    \"team\": \"Heat\",\n    \"points\": 165\n  },\n  {\n    \"name\": \"Rafer Alston\",\n    \"team\": \"Nets\",\n    \"points\": 262\n  }]"),
		eventsTr,
		4*time.Second,
		1*time.Second,
	))

	assert.NotEmpty(t, eventsTr.data)
}
