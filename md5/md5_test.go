package md5

import (
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
		got, err := Hash(tt.r)
		if err != nil {
			t.Error(err)
		}

		if diff := cmp.Diff(got, tt.want); diff != "" {
			t.Errorf("Stdout missmatch (-got +want):\n%s", diff)
		}
	}
}
