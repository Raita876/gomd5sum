package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"

	md5 "gomd5sum/md5"
)

var (
	version string
	name    string
)

type CommandLine struct {
	Paths   []string
	IsBin   bool
	IsText  bool
	IsCheck bool
}

type Option interface {
	apply(*CommandLine)
}

type isBinOption bool

func (ib isBinOption) apply(c *CommandLine) {
	c.IsBin = true
}

type isTextOption bool

func (it isTextOption) apply(c *CommandLine) {
	c.IsText = true
}

type isCheckOption bool

func (ic isCheckOption) apply(c *CommandLine) {
	c.IsCheck = true
}

func (c *CommandLine) Set(opts ...Option) {
	for _, o := range opts {
		o.apply(c)
	}
}

func (c *CommandLine) Md5sum() error {
	hrl := []md5.HashResult{}
	for _, p := range c.Paths {
		hrl = append(hrl, md5.Hash(p))
	}

	for _, hr := range hrl {
		fmt.Printf("%s\n", hr.Print())
	}

	if md5.HasHashError(hrl) {
		return xerrors.New("Failure Md5sum")
	}

	return nil
}

func (c *CommandLine) Exec() error {
	return c.Md5sum()
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
				isTextOption(false),
				isCheckOption(false),
			)

			return cl.Exec()
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
