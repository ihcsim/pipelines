package md5dir

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestBoundedParallelCompute(t *testing.T) {
	if testing.Verbose() {
		log.SetLevel(log.DebugLevel)
	}

	fileCount := 20
	fixture, err := NewFixture(fileCount)
	if err != nil {
		t.Fatal("Unexpected error while setting up fixture: ", err)
	}
	defer os.RemoveAll(fixture.root)

	bounded := &BoundedParallel{bound: 10}
	t.Run("Compute", func(t *testing.T) {
		results, err := bounded.Compute(fixture.root)
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

	t.Run("walkFiles", func(t *testing.T) {
		paths, cerr := bounded.walkFiles(fixture.root, nil)

		actuals := []string{}
		for path := range paths {
			actuals = append(actuals, path)
		}

		for _, expected := range fixture.md5 {
			var found bool
			for _, actual := range actuals {
				if actual == expected.path {
					found = true
					break
				}
			}

			if !found {
				t.Error("Expected file entry is missing: ", expected)
			}
		}

		if err := <-cerr; err != nil {
			t.Error("Unexpected error: ", err)
		}
	})

	t.Run("walkFiles_Done", func(t *testing.T) {
		done := make(chan struct{})
		close(done)

		_, cerr := bounded.walkFiles(fixture.root, done)
		if err := <-cerr; err == nil {
			t.Error("Expected error didn't occur")
		}
	})

	t.Run("walkFiles_NonExistentRoot", func(t *testing.T) {
		_, cerr := bounded.walkFiles("non-existent", nil)
		if err := <-cerr; err == nil {
			t.Error("Expected error didn't occur")
		}
	})

	t.Run("walkFiles_EmptyRoot", func(t *testing.T) {
		emptyRoot, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}

		_, cerr := bounded.walkFiles(emptyRoot, nil)
		if err := <-cerr; err != nil {
			t.Error("Unexpected error: ", err)
		}
	})

	t.Run("sumFiles", func(t *testing.T) {
		paths := make(chan string, len(fixture.md5))
		for _, fileMD5 := range fixture.md5 {
			paths <- fileMD5.path
		}
		close(paths)

		t.Run("", func(t *testing.T) {
			results := bounded.sumFiles(paths, nil)

			actuals := []*FileMD5{}
			for r := range results {
				actuals = append(actuals, r)
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
					t.Error("Missing expected file entry: ", expected)
				}
			}
		})

		t.Run("sumFiles_Done", func(t *testing.T) {
			done := make(chan struct{})
			close(done)

			out := bounded.sumFiles(paths, done)
			if _, ok := <-out; ok {
				t.Error("Expected out channel to be closed")
			}
		})
	})
}
