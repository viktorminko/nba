package simulation

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func readDataFile(path string) (*os.File, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}

	return file, nil
}

func Test_Init(t *testing.T) {
	file, err := readDataFile("../../../players.json")
	if err != nil {
		log.Fatalf("read data file: %#v", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Fatal("close file", cerr)
		}
	}()

	games, err := Init(file)
	assert.NoError(t, err)

	log.Printf("games: %#v", games)
}
