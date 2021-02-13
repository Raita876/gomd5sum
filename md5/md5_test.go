package md5

import (
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/xerrors"
)

func TestPrint(t *testing.T) {
	tests := []struct {
		hr   HashResult
		want string
	}{
		{
			HashResult{
				Md5: Md5{
					Path:  "hoge.txt",
					Value: "eed2d071c1baccf9632f7c61118a1e4a",
				},
				Error: nil,
			},
			"eed2d071c1baccf9632f7c61118a1e4a  hoge.txt",
		},
		{
			HashResult{
				Md5: Md5{
					Path:  "fuga.txt",
					Value: "1347633cdf7cdcb2168d61093630d5ae",
				},
				Error: xerrors.New("open fuga.txt: no such file or directory"),
			},
			"md5sum: open fuga.txt: no such file or directory",
		},
	}

	for _, tt := range tests {
		got := tt.hr.Print()
		if diff := cmp.Diff(got, tt.want); diff != "" {
			t.Errorf("Mismatch (-got +want):\n%s", diff)
		}
	}
}

func TestHash(t *testing.T) {
	tests := []struct {
		r    io.Reader
		want string
	}{
		{
			strings.NewReader("Hello World!"),
			"ed076287532e86365e841e92bfc50d8c",
		},
	}

	for _, tt := range tests {
		got, err := hash(tt.r)
		if err != nil {
			t.Error(err)
		}

		if diff := cmp.Diff(got, tt.want); diff != "" {
			t.Errorf("Mismatch (-got +want):\n%s", diff)
		}
	}
}

func TestHasHashError(t *testing.T) {
	tests := []struct {
		hrl  []HashResult
		want bool
	}{
		{
			[]HashResult{
				{Error: nil},
				{Error: nil},
				{Error: nil},
			},
			false,
		},
		{
			[]HashResult{
				{Error: xerrors.New("open hoge.txt: no such file or directory")},
				{Error: nil},
				{Error: nil},
			},
			true,
		},
	}

	for _, tt := range tests {
		got := HasHashError(tt.hrl)
		if got != tt.want {
			t.Errorf("Mismatch (got=%t, want=%t)", got, tt.want)
		}
	}
}
