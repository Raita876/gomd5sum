package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"

	"gomd5sum/md5"
)

var (
	version string
	name    string
)

type CommandLine struct {
	Paths []string
	IsBin bool
}

type Option interface {
	apply(*CommandLine)
}

type isBinOption bool

func (ib isBinOption) apply(c *CommandLine) {
	c.IsBin = true
}

type HashResult struct {
	Path  string
	Md5   string
	Error error
}

func Hash(path string) HashResult {
	f, err := os.Open(path)
	if err != nil {
		return HashResult{Error: err}
	}

	md5, err := md5.Hash(f)
	if err != nil {
		return HashResult{Error: err}
	}

	return HashResult{
		Path:  path,
		Md5:   md5,
		Error: err,
	}
}

func HasHashError(hrl []HashResult) bool {
	for _, hr := range hrl {
		if hr.Error != nil {
			return true
		}
	}

	return false
}

func (c *CommandLine) Set(opts ...Option) {
	for _, o := range opts {
		o.apply(c)
	}
}

func (c *CommandLine) Exec() error {
	hrl := []HashResult{}

	for _, p := range c.Paths {
		hrl = append(hrl, Hash(p))
	}

	for _, hr := range hrl {
		if hr.Error == nil {
			fmt.Printf("%s  %s\n", hr.Md5, hr.Path)
		} else {
			fmt.Printf("%s: %s\n", name, hr.Error)
		}
	}

	if HasHashError(hrl) {
		return xerrors.New("Failure exec command")
	}

	return nil
}

func main() {
	app := &cli.App{
		Version: version,
		Name:    name,
		Usage:   "gomd5sum",
		Action: func(c *cli.Context) error {
			paths := []string{}
			for i := 0; i < c.Args().Len(); i++ {
				paths = append(paths, c.Args().Get(i))
			}

			cl := &CommandLine{
				Paths: paths,
			}

			cl.Set(
				isBinOption(false),
			)

			return cl.Exec()
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
