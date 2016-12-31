package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestSerialWalk(t *testing.T) {
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

	actuals, err := computeSerial(root)
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	if !reflect.DeepEqual(expected, actuals) {
		t.Errorf("Mismatch result.\nExpected %+v\nBut got %+v", tmpFilesWithMD5, actuals)
	}
}

func tmpFilesWithMD5(root string, fileCount int) (map[string][md5.Size]byte, error) {
	files := make(map[string][md5.Size]byte, fileCount)
	for i := 0; i < fileCount; i++ {
		f, err := ioutil.TempFile(root, "")
		if err != nil {
			return nil, err
		}

		b := []byte(fmt.Sprintf("Test message for file #%d", i))
		if _, err := f.Write(b); err != nil {
			return nil, err
		}
		files[f.Name()] = md5.Sum(b)
	}

	return files, nil
}
