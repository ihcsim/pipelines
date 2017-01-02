package main

import (
	"crypto/md5"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// MD5AllBounded calculates the md5sum of all the files in the root folder.
// It returns a map of filepaths and the MD5 sum of the files.
func MD5AllBounded(root string, bound int) (map[string][md5.Size]byte, error) {
	md5sums := make(map[string][md5.Size]byte)
	done := make(chan struct{})
	defer close(done)

	paths, cerr := walkFiles(root, done)
	for r := range boundedRead(paths, done, bound) {
		md5sums[r.path] = r.sum
	}

	if err := <-cerr; err != nil {
		return nil, err
	}

	return md5sums, nil
}

// walkFiles walks the file tree rooted at root, and returns the name of each file
// to the returned string channel. Any errors encountered during the tree walk are returned
// to the returned error channel.
func walkFiles(root string, done <-chan struct{}) (<-chan string, <-chan error) {
	paths, cerr := make(chan string), make(chan error)

	var wait sync.WaitGroup
	computeMD5 := func(path string, info os.FileInfo, err error) error {
		select {
		case <-done:
			return errors.New("Cancelled processing")
		default:
		}

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		wait.Add(1)
		paths <- path
		wait.Done()
		return nil
	}

	go func() {
		if err := filepath.Walk(root, computeMD5); err != nil {
			cerr <- err
		}

		wait.Wait()
		close(paths)
		close(cerr)
	}()
	return paths, cerr
}

// boundedRead calculates the MD5 sum of the files it receives from paths. It sends the result to
// the returned channel. Bound enforces an upper bound on the number of simultaneous processing.
func boundedRead(paths <-chan string, done <-chan struct{}, bound int) chan result {
	out := make(chan result)
	boundedPath := make(chan string, bound)

	var wait sync.WaitGroup
	go func() {
		for path := range paths {
			select {
			case <-done:
				return
			default:
			}

			boundedPath <- path // blocks when the number of queued items exceeded bound
			wait.Add(1)
			go func() {
				filepath := <-boundedPath
				data, err := ioutil.ReadFile(filepath)
				out <- result{filepath, md5.Sum(data), err}
				wait.Done()
			}()
		}
		wait.Wait()
		close(out)
	}()

	return out
}
