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

// Parallel computes the MD5 checksum of files in a folder using a two-stage pipeline.
// The first stage, sumFiles, walks the tree, digests each file in a new goroutine, and sends the results on a channel. The second stage happens in the Compute method which receives the digest values.
type Parallel struct{}

// Compute computes the md5sum of all the files in the root folder.
// It returns a list of files MD5 checksum, or an error.
func (p *Parallel) Compute(root string) (DirMD5, error) {
	var result DirMD5
	done := make(chan struct{})
	defer close(done)

	out, cerr := sumFiles(root, done)
	for fileMD5 := range out {
		result = append(result, fileMD5)
	}
	if err := <-cerr; err != nil {
		return nil, err
	}

	return result, nil
}

// sumFiles walks through all the files in the root directory, computes their respective md5sum
// and sends them back to the caller via the returned channel. All errors encountered during the file processing are ent tothe returned error channel.
// The done channel is used to stop all processing.
func sumFiles(root string, done <-chan struct{}) (<-chan *FileMD5, <-chan error) {
	out := make(chan *FileMD5)
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
		go func() {
			select {
			case <-done: // if done, skip processing
				warning := NewInterruptWarning("Received done signal")
				log.Debugln(fmt.Sprint(warning))
			default:
				b, err := ioutil.ReadFile(path)
				r := &FileMD5{path: path, sum: md5.Sum(b), err: err}
				log.Debugf("Processing file %+v", r)
				out <- r
			}

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
