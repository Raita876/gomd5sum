package md5

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/xerrors"
)

type HashResult struct {
	Md5   Md5
	Error error
}

type HashResults []HashResult

type CheckResult struct {
	Md5   Md5
	Error error
}

type CheckResults []CheckResult

type Md5 struct {
	Path  string
	Value string
}

func (crl CheckResults) Print() {
	for _, cr := range crl {
		fmt.Printf("%s\n", cr.print())
	}
}

func (cr *CheckResult) print() string {
	if cr.Error != nil {
		return fmt.Sprintf("%s: FAILED", cr.Md5.Path)
	}

	return fmt.Sprintf("%s: OK", cr.Md5.Path)
}

func (hrl HashResults) Print() {
	for _, hr := range hrl {
		fmt.Printf("%s\n", hr.print())
	}
}

func (hr *HashResult) print() string {
	if hr.Error != nil {
		return fmt.Sprintf("md5sum: %s", hr.Error)
	}

	return fmt.Sprintf("%s  %s", hr.Md5.Value, hr.Md5.Path)
}

func Parse(path string) ([]Md5, error) {
	f, err := os.Open(path)
	if err != nil {
		return []Md5{}, err
	}

	return parse(f)
}

func parse(r io.Reader) ([]Md5, error) {
	ml := []Md5{}

	sc := bufio.NewScanner(r)

	for sc.Scan() {
		sl := strings.Split(sc.Text(), "  ")
		ml = append(ml, Md5{
			Path:  sl[1],
			Value: sl[0],
		})
	}

	if err := sc.Err(); err != nil {
		return []Md5{}, err
	}

	return ml, nil
}

func Md5sum(paths []string) HashResults {
	hrl := HashResults{}
	for _, p := range paths {
		hrl = append(hrl, md5sum(p))
	}

	return hrl
}

func md5sum(path string) HashResult {
	f, err := os.Open(path)
	if err != nil {
		return HashResult{
			Md5: Md5{
				Path: path,
			},
			Error: err,
		}
	}

	md5, err := hash(f)
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

func Check(paths []string) (CheckResults, error) {
	ml := []Md5{}
	for _, p := range paths {
		tmpMl, err := Parse(p)
		if err != nil {
			return CheckResults{}, err
		}

		ml = append(ml, tmpMl...)
	}

	crl := CheckResults{}
	for _, m := range ml {
		crl = append(crl, check(m))
	}

	return crl, nil
}

func check(m Md5) CheckResult {
	f, err := os.Open(m.Path)
	if err != nil {
		return CheckResult{
			Md5:   m,
			Error: err,
		}
	}

	h, err := hash(f)
	if err != nil {
		return CheckResult{
			Md5:   m,
			Error: err,
		}
	}

	if h != m.Value {
		return CheckResult{
			Md5:   m,
			Error: xerrors.New("Checksum did not match"),
		}
	}

	return CheckResult{
		Md5:   m,
		Error: nil,
	}
}

func (hrl HashResults) HasError() bool {
	for _, hr := range hrl {
		if hr.Error != nil {
			return true
		}
	}

	return false
}

func (crl CheckResults) HasError() bool {
	for _, hr := range crl {
		if hr.Error != nil {
			return true
		}
	}

	return false
}
