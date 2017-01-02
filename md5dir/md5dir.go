package md5dir

import (
	"crypto/md5"
	"fmt"
)

// MD5Dir can compute the MD5 checksum of all the files in a given root folder.
type MD5Dir interface {
	// Compute computes the MD5 checksum of all the files in the root folder.
	// It returns a list of files MD5 checksum, or an error.
	Compute(root string) (DirMD5, error)
}

// DirMD5 consists of a a list of filepath and the MD5 checksum of the respective files.
type DirMD5 []*FileMD5

// String returns the string representation of m.
func (m DirMD5) String() string {
	var result string
	for _, md5Result := range m {
		result += fmt.Sprintf("%s\n", md5Result)
	}

	return result
}

// FileMD5 contains the path to a file, the file's md5 checksum. and any error that may have occurred during the processing of the file.
type FileMD5 struct {
	path string
	sum  [md5.Size]byte
	err  error
}

// String returns a string representation of m.
func (m FileMD5) String() string {
	return "File: " + m.path + ", Checksum: " + fmt.Sprintf("%x", m.sum) + ", Error: " + fmt.Sprint(m.err)
}

const (
	// 'Interrupt' warning message
	InterruptWarningMessage = "Process interrupted"
)

// InterruptWarning is an error that represents an interruption of the processing.
type InterruptWarning struct {
	msg string
}

// NewInterruptWarning returns a InterruptWarning instance with the specified reason.
func NewInterruptWarning(reason string) *InterruptWarning {
	msg := InterruptWarningMessage + ": " + reason
	return &InterruptWarning{msg: msg}
}

// Error returns the string representation of the InterruptWarning.
func (i InterruptWarning) Error() string {
	return i.msg
}
