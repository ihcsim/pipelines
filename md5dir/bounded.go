package md5dir

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	log "github.com/Sirupsen/logrus"
)

// BoundedParallel can compute the MD5 checksum of files in a folder using a three-stage pipeline.
// The first stage, walkFiles, emits the paths of regular files in the tree. The middle stage starts a fixed number of digester goroutines that receive file names from paths and send results on its returned channel. The final stage receives all the results from c then checks the error from errc
type BoundedParallel struct {
	bound int
}

// Compute computes the MD5 checksum of all the files in the root folder.
// It returns a list of files MD5 checksum, or an error.
func (b *BoundedParallel) Compute(root string) (DirMD5, error) {
	var result DirMD5
	done := make(chan struct{})
	defer close(done)

	paths, cerr := b.walkFiles(root, done)
	for r := range b.sumFiles(paths, done) {
		result = append(result, r)
	}

	if err := <-cerr; err != nil {
		return nil, err
	}

	return result, nil
}

// walkFiles walks the file tree rooted at root, and sends the path of each file to the returned channel. All errors encountered during the file processing are sent to the returned error channel.
// The number of concurrent processing is bounded by b.bound.
func (b *BoundedParallel) walkFiles(root string, done <-chan struct{}) (<-chan string, <-chan error) {
	paths := make(chan string, b.bound)
	cerr := make(chan error, 1)

	var wait sync.WaitGroup
	computeMD5 := func(path string, info os.FileInfo, err error) error {
		select {
		case <-done:
			warning := NewInterruptWarning("Received done signal")
			log.Debugln(fmt.Sprint(warning))
			return warning
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

// sumFiles calculate the MD5 sum of the files it receives from paths. It sends the result to
// the returned channel.
func (b *BoundedParallel) sumFiles(paths <-chan string, done <-chan struct{}) <-chan *FileMD5 {
	out := make(chan *FileMD5)

	var wait sync.WaitGroup
	go func() {
		for path := range paths {
			select {
			case <-done:
				warning := NewInterruptWarning("Received done signal")
				log.Debugln(fmt.Sprint(warning))
				return
			default:
			}

			wait.Add(1)
			go func(filepath string) {
				data, err := ioutil.ReadFile(filepath)
				out <- &FileMD5{path: filepath, sum: md5.Sum(data), err: err}
				wait.Done()
			}(path)
		}
		wait.Wait()
		close(out)
	}()

	return out
}
