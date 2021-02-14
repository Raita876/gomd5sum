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

func (c *CommandLine) Check() error {
	ml := []md5.Md5{}
	for _, p := range c.Paths {
		tmpMl, err := md5.Parse(p)
		if err != nil {
			return err
		}

		ml = append(ml, tmpMl...)
	}

	crl := []md5.CheckResult{}
	for _, m := range ml {
		crl = append(crl, md5.Check(m))
	}

	for _, cr := range crl {
		fmt.Printf("%s\n", cr.Print())
	}

	if md5.HasCheckError(crl) {
		return xerrors.New("Failure Check")
	}

	return nil
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
	if c.IsCheck {
		return c.Check()
	}

	return c.Md5sum()
}

func main() {
	app := &cli.App{
		Version: version,
		Name:    name,
		Usage:   "gomd5sum",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "check",
				Aliases: []string{"c"},
				Usage:   "",
			},
		},
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
				isCheckOption(c.Bool("check")),
			)

			return cl.Exec()
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
