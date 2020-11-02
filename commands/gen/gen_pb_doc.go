// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_gen

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"path"
)

type genPbDocInfo struct {
	Sources *cli.StringSlice
	Imports *cli.StringSlice
	Type    *cli.StringSlice
	Output  string
	Name    string
}

var GenPbDocInfo = genPbDocInfo{
	Sources: cli.NewStringSlice(),
	Imports: cli.NewStringSlice(),
	Type:    cli.NewStringSlice(),
	Output:  "",
	Name:    "",
}

// Generate Go files from proto
func GenPbDocCommand() *cli.Command {
	command := &cli.Command{
		Name:      "pb-doc",
		Usage:     "generate documentation from proto",
		UsageText: "rk gen pb-doc -s [source] -i [import] -o [output] -t [output type]",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:        "source, s",
				Aliases:     []string{"s"},
				Destination: GenPbDocInfo.Sources,
				Required:    true,
				Usage:       "where proto file located",
			},
			&cli.StringFlag{
				Name:        "output, o",
				Aliases:     []string{"o"},
				Destination: &GenPbDocInfo.Output,
				Required:    false,
				Usage:       "documentation output type",
			},
			&cli.StringSliceFlag{
				Name:        "import, i",
				Aliases:     []string{"i"},
				Destination: GenPbDocInfo.Imports,
				Required:    false,
				Usage:       "specify files need to import",
			},
			&cli.StringFlag{
				Name:        "name, n",
				Aliases:     []string{"n"},
				Destination: &GenPbDocInfo.Name,
				Required:    false,
				Usage:       "document name",
			},
			&cli.StringSliceFlag{
				Name:        "type, t",
				Aliases:     []string{"t"},
				Destination: GenPbDocInfo.Type,
				Required:    false,
				Usage:       "specify output file type",
			},
		},
		Action: GenPbDocAction,
	}

	return command
}

func GenPbDocAction(ctx *cli.Context) error {
	err := genDoc()

	if err != nil {
		return err
	}

	return nil
}

func genDoc() error {
	// run protoc command
	color.Cyan("Generate documentation file from proto at %v to docs folder", GenPbDocInfo.Sources)

	args := make([]string, 0)
	args = append(args, "-I.")

	// add imports
	addImports(GenPbDocInfo.Imports.Value(), &args)

	wd, _ := os.Getwd()

	// add output path
	// Set with working directory
	if !path.IsAbs(GenPbDocInfo.Output) {
		if len(GenPbDocInfo.Output) < 1 {
			GenPbDocInfo.Output = "docs"
		}
		GenPbDocInfo.Output = path.Join(wd, GenPbDocInfo.Output)

		// create folder
		os.MkdirAll(GenPbDocInfo.Output, os.ModePerm)
	}

	args = append(args, fmt.Sprintf("--doc_out=%s", GenPbDocInfo.Output))

	// add doc type
	if len(GenPbDocInfo.Name) < 1 {
		GenPbDocInfo.Name = path.Base(wd)
	}

	types := GenPbDocInfo.Type.Value()
	if len(types) < 1 {
		types = append(types, "html", "markdown")
	}

	for i := range types {
		copiedArgs := make([]string, len(args))
		copy(copiedArgs, args)
		genDocWithType(types[i], copiedArgs)
	}

	color.Green("[Done]")

	return nil
}

func genDocWithType(docType string, args []string) {
	args = append(args, fmt.Sprintf("--doc_opt=%s,%s.%s", docType, GenPbDocInfo.Name, docType))
	// add sources
	addSources(GenPbDocInfo.Sources.Value(), &args)

	bytes, err := exec.Command("protoc", args...).CombinedOutput()
	if err != nil {
		color.Red("failed to build proto to documentation file\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return
	}

	if GenInfo.Debug {
		color.Blue("[stdout] %s", string(bytes))
	}
}
