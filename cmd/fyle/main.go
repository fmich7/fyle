package main

import (
	"github.com/fmich7/fyle/pkg/cli"
	"github.com/spf13/afero"
)

func main() {
	cli := cli.NewCliClient(afero.NewOsFs())
	cli.Execute()
}
