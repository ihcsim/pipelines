package main

import (
	"crypto/md5"
	"io/ioutil"
	"os"
	"path/filepath"
)

// computeSerial reads all the files in the folder root and returns a map
// from file path to the MD5 sum of the file's contents.  If the directory walk
// fails or any read operation fails, MD5All returns an error.
func computeSerial(root string) (map[string][md5.Size]byte, error) {
	md5sums := make(map[string][md5.Size]byte)
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

		md5sums[path] = md5.Sum(b)
		return nil
	}

	if err := filepath.Walk(root, computeMD5); err != nil {
		return nil, err
	}
	return md5sums, nil
}
