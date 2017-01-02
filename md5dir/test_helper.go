package md5dir

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type fixture struct {
	root string
	md5  DirMD5
}

func NewFixture(fileCount int) (*fixture, error) {
	root, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	f, err := computeDirMD5(root, fileCount)
	if err != nil {
		return nil, err
	}

	log.Debugf("Test Fixture\n%s", f)
	return &fixture{root: root, md5: f}, nil
}

func computeDirMD5(root string, fileCount int) (DirMD5, error) {
	var results DirMD5
	for i := 0; i < fileCount; i++ {
		f, err := ioutil.TempFile(root, "")
		if err != nil {
			return nil, err
		}

		b := []byte(fmt.Sprintf("Test message for file #%d", i))
		if _, err := f.Write(b); err != nil {
			return nil, err
		}

		r := &FileMD5{path: f.Name(), sum: md5.Sum(b), err: err}
		results = append(results, r)
	}

	return results, nil
}
