package md5

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

type HashResult struct {
	Md5   Md5
	Error error
}

type Md5 struct {
	Path  string
	Value string
}

func (hr *HashResult) Print() string {
	if hr.Error != nil {
		return fmt.Sprintf("md5sum: %s", hr.Error)
	}

	return fmt.Sprintf("%s  %s", hr.Md5.Value, hr.Md5.Path)
}

func Hash(path string) HashResult {
	f, err := os.Open(path)
	if err != nil {
		return HashResult{Error: err}
	}

	md5, err := hash(f)
	if err != nil {
		return HashResult{Error: err}
	}

	return HashResult{
		Md5: Md5{
			Path:  path,
			Value: md5,
		},
		Error: err,
	}
}

func hash(r io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func HasHashError(hrl []HashResult) bool {
	for _, hr := range hrl {
		if hr.Error != nil {
			return true
		}
	}

	return false
}
