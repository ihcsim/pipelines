package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestParallelMD5(t *testing.T) {
	root, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	defer os.RemoveAll(root)

	fileCount := 5
	expected, err := tmpFilesWithMD5(root, fileCount)
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	t.Run("MD5All", func(t *testing.T) {
		actual, err := MD5All(root)
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Mismatched result.\n  Expected %+v\n  But got %+v", expected, actual)
		}
	})

	t.Run("MD5All_NonExistentRoot", func(t *testing.T) {
		_, err := MD5All("non-existent")
		if err == nil {
			t.Error("Expected error didn't occur")
		}
	})

	t.Run("MD5All_EmptyRoot", func(t *testing.T) {
		emptyRoot, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}
		defer os.RemoveAll(emptyRoot)

		if _, err := MD5All(emptyRoot); err != nil {
			t.Error("Unexpected error: ", err)
		}
	})

	t.Run("sumFiles", func(t *testing.T) {
		out, cerr := sumFiles(root, nil)
		var count int
		for {
			select {
			case result := <-out:
				expectedSum, ok := expected[result.path]
				if !ok {
					t.Error("Unexpected file path: ", result.path)
				}

				if expectedSum != result.sum {
					t.Errorf("Mismatched MD5 sum. Expected %s, but got %s", expectedSum, result.sum)
				}

				count++
				if count == fileCount {
					return
				}
			case err := <-cerr:
				t.Error("Unexpected error: ", err)
			}
		}
	})

	t.Run("sumFiles_Done", func(t *testing.T) {
		done := make(chan struct{})
		close(done)

		_, cerr := sumFiles(root, done)
		if err := <-cerr; err == nil {
			t.Error("Expected error didn't occur")
		}
	})

	t.Run("sumFilesError_NonExistentRoot", func(t *testing.T) {
		done := make(chan struct{})
		_, cerr := sumFiles("nonexistent", done)
		if err := <-cerr; err == nil {
			t.Error("Expected error didn't occur")
		}
	})

	t.Run("sumFilesError_EmptyRoot", func(t *testing.T) {
		done := make(chan struct{})
		emptyRoot, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}
		defer os.RemoveAll(emptyRoot)

		_, cerr := sumFiles(emptyRoot, done)
		if err := <-cerr; err != nil {
			t.Error("Unexpected error: ", err)
		}
	})
}
