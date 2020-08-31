package rk_gen

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

type ProtoConfig struct {
	Proto *Proto `yaml:"proto"`
}
type Doc struct {
	Output string   `yaml:"output"`
	Name   string   `yaml:"name"`
	Type   []string `yaml:"type"`
}
type Proto struct {
	Sources []string `yaml:"source"`
	Imports []string `yaml:"import"`
	Doc     *Doc     `yaml:"doc"`
}

type genPbGoInfo struct {
	ConfigFilePath string
	Sources        *cli.StringSlice
	Imports        *cli.StringSlice
}

var GenPbGoInfo = genPbGoInfo{
	ConfigFilePath: "",
	Sources:        cli.NewStringSlice(),
	Imports:        cli.NewStringSlice(),
}

// Generate Go files from proto
func GenPbGoCommand() *cli.Command {
	command := &cli.Command{
		Name:      "pb-go",
		Usage:     "generate go & doc file from proto",
		UsageText: "rk gen pb-go -s [source] -i [import] -f [file]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config, f",
				Aliases:     []string{"f"},
				Destination: &GenPbGoInfo.ConfigFilePath,
				Required:    false,
				Usage:       "config file path, flag will override config file",
			},
			&cli.StringSliceFlag{
				Name:        "sources, s",
				Aliases:     []string{"s"},
				Destination: GenPbGoInfo.Sources,
				Required:    false,
				Usage:       "where proto files located",
			},
			&cli.StringSliceFlag{
				Name:        "imports, i",
				Aliases:     []string{"i"},
				Destination: GenPbGoInfo.Imports,
				Required:    false,
				Usage:       "specify files need to import",
			},
		},
		Action: GenPbGoAction,
	}

	return command
}

func GenPbGoAction(ctx *cli.Context) error {
	config := readProtoConfig()

	err := genGo()
	err = genGateway()
	err = genSwagger()

	if config != nil && config.Proto.Doc != nil {
		err = genDoc()
	}

	if err != nil {
		return err
	}

	return nil
}

func genGo() error {
	color.Cyan("Generate go file from proto at %v to the same path", GenPbGoInfo.Sources)

	args := make([]string, 0)
	args = append(args, "-I.")

	// add imports
	addImports(GenPbGoInfo.Imports.Value(), &args)

	// add plugins
	args = append(args, "--go_out=plugins=grpc:.")

	// add output path
	args = append(args, "--go_opt=paths=source_relative")

	// add source
	if err := addSources(GenPbGoInfo.Sources.Value(), &args); err != nil {
		return err
	}

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
	color.Cyan("Generate gRpc gateway from proto at %v to the same path", GenPbGoInfo.Sources)

	args := make([]string, 0)
	args = append(args, "-I.")
	// add imports
	addImports(GenPbGoInfo.Imports.Value(), &args)

	// add output path
	args = append(args, "--grpc-gateway_out=logtostderr=true,paths=source_relative:.")

	if err := addSources(GenPbGoInfo.Sources.Value(), &args); err != nil {
		return err
	}

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
	color.Cyan("Generate swagger json file from proto at %v to the same path", GenPbGoInfo.Sources)

	args := make([]string, 0)
	args = append(args, "-I.")
	// add imports
	addImports(GenPbGoInfo.Imports.Value(), &args)

	// add output path
	args = append(args, "--grpc-gateway_out=logtostderr=true,paths=source_relative:.",
		"--swagger_out=logtostderr=true:.")

	if err := addSources(GenPbGoInfo.Sources.Value(), &args); err != nil {
		return err
	}

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

func addImports(imports []string, args *[]string) {
	for i := range imports {
		element := imports[i]
		*args = append(*args, "-I", element)
	}
}

func addSources(sources []string, args *[]string) error {
	for i := range sources {
		element := sources[i]

		// list files if path is composed with wildcards
		files, err := filepath.Glob(element)
		if err != nil {
			color.Red("failed to list proto files with path:%s to go file\n[err] %v", element, err)
			return err
		}

		for j := range files {
			*args = append(*args, files[j])
		}
	}

	return nil
}

func readProtoConfig() *ProtoConfig {
	var config *ProtoConfig



	if len(GenPbGoInfo.ConfigFilePath) > 0 {
		config = fromConfigFile()
	}

	if config == nil {
		return nil
	}

	for i := range config.Proto.Sources {
		GenPbGoInfo.Sources.Set(config.Proto.Sources[i])
	}

	for i := range config.Proto.Imports {
		GenPbGoInfo.Imports.Set(config.Proto.Imports[i])
	}

	// we are ready to generate docs
	if config.Proto.Doc != nil {
		for i := range config.Proto.Imports {
			GenPbDocInfo.Imports.Set(config.Proto.Imports[i])
		}

		for i := range config.Proto.Sources {
			GenPbDocInfo.Sources.Set(config.Proto.Sources[i])
		}

		for i := range config.Proto.Doc.Type {
			GenPbDocInfo.Type.Set(config.Proto.Doc.Type[i])
		}

		GenPbDocInfo.Output = config.Proto.Doc.Output
		GenPbDocInfo.Name = config.Proto.Doc.Name
	}

	return config
}

func fromConfigFile() *ProtoConfig {
	path := GenPbGoInfo.ConfigFilePath
	if !isValidPath(path) {
		return nil
	}

	bytes, ext, err := readFile(path)

	if err != nil {
		return nil
	}

	if ext == ".yaml" || ext == ".yml" {
		// unmarshal yaml
		config := &ProtoConfig{}
		if err = yaml.Unmarshal(bytes, config); err != nil {
			return nil
		}

		return config
	}

	return nil
}

func isValidPath(fp string) bool {
	// Check if file already exists
	if _, err := os.Stat(fp); err == nil {
		return true
	}

	return false
}

func readFile(filePath string) ([]byte, string, error) {
	if !path.IsAbs(filePath) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, "", err
		}
		filePath = path.Join(wd, filePath)
	}

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, "", err
	}

	ext := path.Ext(filePath)

	return bytes, ext, nil
}
