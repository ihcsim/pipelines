package md5dir

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestParallelCompute(t *testing.T) {
	if testing.Verbose() {
		log.SetLevel(log.DebugLevel)
	}

	fileCount := 5
	fixture, err := NewFixture(fileCount)
	if err != nil {
		t.Fatal("Unexpected error while setting up fixture: ", err)
	}
	defer os.RemoveAll(fixture.root)

	parallel := &Parallel{}
	t.Run("Compute", func(t *testing.T) {
		results, err := parallel.Compute(fixture.root)
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}

		for _, expected := range fixture.md5 {
			var found bool
			for _, actual := range results {
				if reflect.DeepEqual(actual, expected) {
					found = true
					break
				}
			}

			if !found {
				t.Error("Missing expected file checksum: ", expected)
			}
		}
	})

	t.Run("Compute_NonExistentRoot", func(t *testing.T) {
		_, err := parallel.Compute("non-existent")
		if err == nil {
			t.Error("Expected error didn't occur")
		}
	})

	t.Run("Compute_EmptyRoot", func(t *testing.T) {
		emptyRoot, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}
		defer os.RemoveAll(emptyRoot)

		if _, err := parallel.Compute(emptyRoot); err != nil {
			t.Error("Unexpected error: ", err)
		}
	})

	t.Run("sumFiles", func(t *testing.T) {
		out, cerr := sumFiles(fixture.root, nil)

		actuals := []*FileMD5{}
		for actual := range out {
			actuals = append(actuals, actual)
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
				t.Error("Missing expected checksum entry: ", expected)
			}
		}

		if err := <-cerr; err != nil {
			t.Fatal("Unexpected error: ", err)
		}
	})

	t.Run("sumFiles_Done", func(t *testing.T) {
		done := make(chan struct{})
		close(done)

		_, cerr := sumFiles(fixture.root, done)
		if err := <-cerr; err == nil {
			t.Error("Expected error didn't occur")
		}
	})

	t.Run("sumFiles_NonExistentRoot", func(t *testing.T) {
		done := make(chan struct{})
		_, cerr := sumFiles("nonexistent", done)
		if err := <-cerr; err == nil {
			t.Error("Expected error didn't occur")
		}
	})

	t.Run("sumFiles_EmptyRoot", func(t *testing.T) {
		emptyRoot, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}
		defer os.RemoveAll(emptyRoot)

		_, cerr := sumFiles(emptyRoot, nil)
		if err := <-cerr; err != nil {
			t.Error("Unexpected error: ", err)
		}
	})
}
