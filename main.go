package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	version string
	name    string
)

type CommandLine struct {
	Paths []string
}

func (c *CommandLine) Exec() error {
	for _, p := range c.Paths {
		f, err := os.Open(p)
		if err != nil {
			return err
		}

		md5, err := Hash(f)
		if err != nil {
			return err
		}

		fmt.Printf("%s  %s\n", md5, p)
	}

	return nil
}

func Hash(r io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
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

			cl := CommandLine{
				Paths: paths,
			}

			return cl.Exec()
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
