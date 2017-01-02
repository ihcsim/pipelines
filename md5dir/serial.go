package md5dir

import (
	"crypto/md5"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Serial computes the MD5 checksum of files in a folder in sequential lexical order.
type Serial struct{}

// Compute computes the MD5 checksum of all the files in the root folder.
// It returns a list of files MD5 checksum, or an error.
func (s *Serial) Compute(root string) (DirMD5, error) {
	var m DirMD5
	computeMD5 := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		m = append(m, &FileMD5{path: path, sum: md5.Sum(b), err: err})
		return nil
	}

	if err := filepath.Walk(root, computeMD5); err != nil {
		return nil, err
	}
	return m, nil
}
