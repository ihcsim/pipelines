package main

import (
	"crypto/md5"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// result contains the path to a file and the file's md5sum.
// Any errors that occurred during the processing of the file are captured as err.
type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

// MD5All calculates the md5sum of all the files in the root folder.
// It returns a map of filepaths and the MD5 sum of the files.
func MD5All(root string) (map[string][md5.Size]byte, error) {
	md5sums := make(map[string][md5.Size]byte)
	done := make(chan struct{})
	defer close(done)

	out, cerr := sumFiles(root, done)
LOOP:
	for {
		select {
		case r, ok := <-out:
			if !ok {
				break LOOP
			}
			md5sums[r.path] = r.sum
		case err := <-cerr:
			if err != nil {
				return nil, err
			}
		}
	}

	return md5sums, nil
}

// sumFiles walks through all the files in the root directory, calculates their respective md5sum
// and sends them back to the caller via a result channel. Any errors encountered during the file
// processing are sent back via an error channel.
// The done channel is used to store all processing.
func sumFiles(root string, done <-chan struct{}) (<-chan result, <-chan error) {
	out := make(chan result)
	cerr := make(chan error)

	var wait sync.WaitGroup
	computeMD5 := func(path string, info os.FileInfo, err error) error {
		select {
		case <-done:
			return errors.New("Canceled processing")
		default:
		}

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		wait.Add(1)
		go func() {
			b, err := ioutil.ReadFile(path)
			out <- result{path: path, sum: md5.Sum(b), err: err}
			wait.Done()
		}()

		return nil
	}

	go func() {
		if err := filepath.Walk(root, computeMD5); err != nil {
			cerr <- err
		}

		wait.Wait()
		close(out)
		close(cerr)
	}()
	return out, cerr
}
