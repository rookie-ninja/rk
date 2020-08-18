package rk_gen

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"os/exec"
)

type genPbGoInfo struct {
	Source string
	Imports *cli.StringSlice
}

var GenPbGoInfo = genPbGoInfo{
	Source: "",
	Imports: cli.NewStringSlice(),
}

// Generate Go files from proto
func GenPbGoCommand() *cli.Command{
	command := &cli.Command{
		Name: "pb-go",
		Usage: "generate go file from proto",
		UsageText: "rk gen pb-go -s [source] -i [import]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "source, s",
				Aliases:     []string{"s"},
				Destination: &GenPbGoInfo.Source,
				Required: 	 true,
				Usage:       "where proto file located",
			},
			&cli.StringSliceFlag{
				Name:        "imports, i",
				Aliases:     []string{"i"},
				Destination: GenPbGoInfo.Imports,
				Required: 	 false,
				Usage:       "specify files need to import",
			},
		},
		Action: GenPbGoAction,
	}

	return command
}

func GenPbGoAction(ctx *cli.Context) error {
	err := genGo()
	err = genGateway()
	err = genSwagger()

	if err != nil {
		return err
	}

	return nil
}

func genGo() error {
	// run protoc command
	color.Cyan("Generate go file from proto at %s to the same path", GenPbGoInfo.Source)

	// protoc -I. -I third-party/googleapis --go_out=plugins=grpc:. --go_opt=paths=source_relative api/*.proto
	args := make([]string, 0)
	args = append(args, "-I.")
	// add imports
	for i := range GenPbGoInfo.Imports.Value() {
		element := GenPbGoInfo.Imports.Value()[i]
		args = append(args, "-I", element)
	}

	// add plugins
	args = append(args, "--go_out=plugins=grpc:.")

	// add output path
	args = append(args, "--go_opt=paths=source_relative")

	args = append(args, GenPbGoInfo.Source)

	bytes, err := exec.Command("protoc", args...).CombinedOutput()
	if err != nil {
		color.Red("failed to build proto to go file\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	if GenInfo.Debug {
		color.Blue("[stdout] %s", string(bytes))
	}

	color.Green("[Done]")

	return nil
}

func genGateway() error {
	// run protoc command
	color.Cyan("Generate gRpc gateway from proto at %s to the same path", GenPbGoInfo.Source)

	//protoc -I. -I third-party/googleapis --grpc-gateway_out=logtostderr=true,paths=source_relative:. api/*.proto
	args := make([]string, 0)
	args = append(args, "-I.")
	// add imports
	for i := range GenPbGoInfo.Imports.Value() {
		element := GenPbGoInfo.Imports.Value()[i]
		args = append(args, "-I", element)
	}

	// add output path
	args = append(args, "--grpc-gateway_out=logtostderr=true,paths=source_relative:.")

	args = append(args, GenPbGoInfo.Source)

	bytes, err := exec.Command("protoc", args...).CombinedOutput()
	if err != nil {
		color.Red("failed to build proto to gateway file\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	if GenInfo.Debug {
		color.Blue("[stdout] %s", string(bytes))
	}

	color.Green("[Done]")

	return nil
}

func genSwagger() error {
	// run protoc command
	color.Cyan("Generate swagger json file from proto at %s to the same path", GenPbGoInfo.Source)

	// protoc -I. -I third-party/googleapis --grpc-gateway_out=logtostderr=true,paths=source_relative:. --swagger_out=logtostderr=true:. api/*.proto
	args := make([]string, 0)
	args = append(args, "-I.")
	// add imports
	for i := range GenPbGoInfo.Imports.Value() {
		element := GenPbGoInfo.Imports.Value()[i]
		args = append(args, "-I", element)
	}

	// add output path
	args = append(args, "--grpc-gateway_out=logtostderr=true,paths=source_relative:.", "--swagger_out=logtostderr=true:.")

	args = append(args, GenPbGoInfo.Source)

	bytes, err := exec.Command("protoc", args...).CombinedOutput()
	if err != nil {
		color.Red("failed to build proto to swagger file\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	if GenInfo.Debug {
		color.Blue("[stdout] %s", string(bytes))
	}

	color.Green("[Done]")

	return nil
}