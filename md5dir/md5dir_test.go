package md5dir

import (
	"crypto/md5"
	"fmt"
	"testing"
)

func TestDirMD5String(t *testing.T) {
	var m DirMD5

	expected := [10]*FileMD5{}
	for i := 0; i < 10; i++ {
		b := []byte("This is test message " + fmt.Sprintf("%s", i))
		sum := md5.Sum(b)
		md5Result := &FileMD5{path: "file-01", sum: sum, err: nil}
		m = append(m, md5Result)
		expected[i] = md5Result
	}

	for index, actual := range m {
		if fmt.Sprint(actual) != fmt.Sprint(expected[index]) {
			t.Error("Mismatch results. Expected %s, but got %s", fmt.Sprint(expected[index]), fmt.Sprint(actual))
		}
	}
}

func TestFileMD5String(t *testing.T) {
	b := []byte("This is a test message")
	sum := md5.Sum(b)
	m := FileMD5{path: "file-01", sum: sum, err: nil}

	expected := "File: file-01, Checksum: " + fmt.Sprintf("%x", sum) + ", Error: <nil>"
	actual := fmt.Sprint(m)
	if actual != expected {
		t.Errorf("Mismatch result. Expected %s, but got %s", expected, actual)
	}
}
