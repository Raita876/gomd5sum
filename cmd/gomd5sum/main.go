package main

import (
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
	c.IsBin = bool(ib)
}

type isTextOption bool

func (it isTextOption) apply(c *CommandLine) {
	c.IsText = bool(it)
}

type isCheckOption bool

func (ic isCheckOption) apply(c *CommandLine) {
	c.IsCheck = bool(ic)
}

func (c *CommandLine) Set(opts ...Option) {
	for _, o := range opts {
		o.apply(c)
	}
}

func (c *CommandLine) Check() error {
	crl, err := md5.Check(c.Paths)
	if err != nil {
		return err
	}

	crl.Print()

	if crl.HasError() {
		return xerrors.New("Failure Check")
	}

	return nil
}

func (c *CommandLine) Md5sum() error {
	hrl := md5.Md5sum(c.Paths)

	hrl.Print()

	if hrl.HasError() {
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
