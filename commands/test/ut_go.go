package rk_test

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

type UtGoConfig struct {
	Ut *Ut `yaml:"ut"`
}
type Ut struct {
	Output string `yaml:"output"`
}

type utGoInfo struct {
	Output         string
	ConfigFilePath string
}

var UtGoInfo = utGoInfo{
	Output:         "",
	ConfigFilePath: "",
}

// Generate Go files from proto
func UtGoCommand() *cli.Command {
	command := &cli.Command{
		Name:      "ut-go",
		Usage:     "run unit test in golang",
		UsageText: "rk test ut-go -o [output] -f [file]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config file, f",
				Aliases:     []string{"f"},
				Destination: &UtGoInfo.ConfigFilePath,
				Required:    false,
				Usage:       "config file path, flag will override config file",
			},
			&cli.StringFlag{
				Name:        "output, o",
				Aliases:     []string{"o"},
				Destination: &UtGoInfo.Output,
				Required:    false,
				Usage:       "test result output type",
			},
		},
		Action: UtGoAction,
	}

	return command
}

func UtGoAction(ctx *cli.Context) error {
	readUtGoConfig()

	color.Cyan("Run unit-test for golang to output file at %s", UtGoInfo.Output)
	err := utGo()
	err = utGoDoc()

	if err != nil {
		return err
	}

	return nil
}

func utGo() error {
	color.Cyan("Run unit-test for golang")
	args := make([]string, 0)

	args = append(args, "test", "./...", "-v", "short", "-coverprofile=rk_ut.out")

	bytes, err := exec.Command("go", args...).CombinedOutput()
	if err != nil {
		color.Red("failed to run go test\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	if TestInfo.Debug {
		color.Blue("[stdout] %s", string(bytes))
	}

	color.Green("[Done]")

	return nil
}

func utGoDoc() error {
	color.Cyan("Convert to go unit-test doc %s", UtGoInfo.Output)
	args := make([]string, 0)

	// validate if output is a file or not
	if len(UtGoInfo.Output) < 1 {
		UtGoInfo.Output = "docs"
	}

	if !strings.Contains(UtGoInfo.Output, ".") {
		UtGoInfo.Output = path.Join(UtGoInfo.Output, "rk_ut.html")
	}

	args = append(args, "tool", "cover", "-html=rk_ut.out", "-o", UtGoInfo.Output)

	bytes, err := exec.Command("go", args...).CombinedOutput()
	if err != nil {
		color.Red("failed to generate ut docs\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))

		os.Remove("rk_ut.out")
		return err
	}

	if TestInfo.Debug {
		color.Blue("[stdout] %s", string(bytes))
	}

	os.Remove("rk_ut.out")

	color.Green("[Done]")

	return nil
}

func readUtGoConfig() *UtGoConfig {
	var config *UtGoConfig
	if len(UtGoInfo.ConfigFilePath) > 0 {
		config = fromConfigFile()
	}

	if config == nil {
		return nil
	}

	if config.Ut != nil {
		UtGoInfo.Output = config.Ut.Output
	}

	return config
}

func fromConfigFile() *UtGoConfig {
	path := UtGoInfo.ConfigFilePath
	if !isValidPath(path) {
		return nil
	}

	bytes, ext, err := readFile(path)

	if err != nil {
		return nil
	}

	if ext == ".yaml" || ext == ".yml" {
		// unmarshal yaml
		config := &UtGoConfig{}
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
