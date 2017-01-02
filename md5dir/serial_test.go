package md5dir

import (
	"os"
	"reflect"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestSerialCompute(t *testing.T) {
	if testing.Verbose() {
		log.SetLevel(log.DebugLevel)
	}

	fileCount := 5
	fixture, err := NewFixture(fileCount)
	if err != nil {
		t.Fatal("Unexpected error while setting up fixture: ", err)
	}
	defer os.RemoveAll(fixture.root)

	serial := &Serial{}
	actuals, err := serial.Compute(fixture.root)
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	for _, expected := range fixture.md5 {
		var found bool
		for _, actual := range actuals {
			if reflect.DeepEqual(actual, expected) {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Missing expected file checksum.\nExpected to be present:\n%+v\n", expected)
		}
	}
}
