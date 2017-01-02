package main

import (
	"crypto/md5"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestMD5AllBounded(t *testing.T) {
	root, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	defer os.RemoveAll(root)

	fileCount := 15
	expected, err := tmpFilesWithMD5(root, fileCount)
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	bound := 10
	t.Run("MD5AllBounded", func(t *testing.T) {
		actuals, err := MD5AllBounded(root, bound)
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}

		if !reflect.DeepEqual(expected, actuals) {
			t.Errorf("Mismatched results.\n  Expected %+v\nBut got %+v")
		}
	})

	t.Run("WalkFiles", func(t *testing.T) {
		paths, cerr := walkFiles(root, nil)
	LOOP:
		for {
			select {
			case path, ok := <-paths:
				if !ok {
					break LOOP
				}

				if _, ok := expected[path]; !ok {
					t.Error("Unexpected filepath entry: ", path)
				}
			case err := <-cerr:
				if err != nil {
					t.Error("Unexpected error: ", err)
				}
			}
		}
	})

	t.Run("WalkFiles_Done", func(t *testing.T) {
		done := make(chan struct{})
		close(done)

		_, cerr := walkFiles(root, done)
		if err := <-cerr; err == nil {
			t.Error("Expected error didn't occur")
		}
	})

	t.Run("WalkFiles_NonExistentRoot", func(t *testing.T) {
		_, cerr := walkFiles("non-existent", nil)
		if err := <-cerr; err == nil {
			t.Error("Expected error didn't occur")
		}
	})

	t.Run("WalkFiles_EmptyRoot", func(t *testing.T) {
		emptyRoot, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}

		_, cerr := walkFiles(emptyRoot, nil)
		if err := <-cerr; err != nil {
			t.Error("Unexpected error: ", err)
		}
	})

	t.Run("BoundedRead", func(t *testing.T) {
		paths := make(chan string)
		go func() {
			for filepath := range expected {
				paths <- filepath
			}
			close(paths)
		}()

		t.Run("Read", func(t *testing.T) {
			out := boundedRead(paths, nil, bound)
			for actual := range out {
				expectedSum, ok := expected[actual.path]
				if !ok {
					t.Error("Unexpected file entry: ", actual.path)
				}

				if expectedSum != actual.sum {
					t.Errorf("Mismatched MD5 sum. Expected %s, but got %s", string(expectedSum[:md5.Size]), string(actual.sum[:md5.Size]))
				}
			}
		})

		t.Run("Done", func(t *testing.T) {
			done := make(chan struct{})
			close(done)

			out := boundedRead(paths, done, bound)
			if _, ok := <-out; ok {
				t.Error("Expected out channel to be closed")
			}
		})
	})
}
