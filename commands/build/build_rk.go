package rk_build

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/commands/gen"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
)

type BuildGoConfig struct {
	AppName string   `yaml:"appName"`
	BuildGo *BuildGo `yaml:"build"`
	Env map[string]string `yaml:"env"`
}
type BuildGo struct {
	GoOs     string   `yaml:"goos"`
	GoArch   string   `yaml:"goarch"`
}

type buildGoInfo struct {
	AppName      string
	GoOs         string
	GoArch       string
	ConfigFilePath string
}

var BuildGoInfo = buildGoInfo{
	AppName:      "",
	GoOs:         "",
	GoArch:       "",
	ConfigFilePath: "",
}

// Run go build command
func BuildGoCommand() *cli.Command {
	command := &cli.Command{
		Name:      "go",
		Usage:     "build golang program",
		UsageText: "rk build go -os [goos] -arch [goarch] -f [file]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config file, f",
				Aliases:     []string{"f"},
				Destination: &BuildGoInfo.ConfigFilePath,
				Required:    false,
				Usage:       "config file path, flag will override config file",
			},
			&cli.StringFlag{
				Name:        "os",
				Destination: &BuildGoInfo.GoOs,
				Required:    false,
				Usage:       "set goos",
			},
			&cli.StringFlag{
				Name:        "arch",
				Destination: &BuildGoInfo.GoArch,
				Required:    false,
				Usage:       "set goarch",
			},

		},
		Action: BuildGoAction,
	}

	return command
}

func BuildGoAction(ctx *cli.Context) error {
	// for rk style, we will do the following steps
	// 1: read config
	config := readBuildGoConfig()

	// 2: Check before scripts locates at build/before.sh
	if err := runScript("build/before.sh"); err != nil {
		return nil
	}

	if !isValidPath("build/override.sh") {
		// 3: generate files based on proto
		if err := compileProto(ctx); err != nil {
			return err
		}

		// 4: build go files
		if err := buildGoFiles(); err != nil {
			return err
		}

		// 5: create target folder
		if err := createTargetFolder(); err != nil {
			return err
		}

		// 6: copy to target
		if err := copyFiles(config); err != nil {
			return err
		}
	} else {
		if err := runScript("build/override.sh"); err != nil {
			return nil
		}
	}

	if err := runScript("build/after.sh"); err != nil {
		return nil
	}

	return nil
}

func compileProto(ctx *cli.Context) error {
	cmd := rk_gen.GenPbGoCommand()
	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "config, f",
			Value: BuildGoInfo.ConfigFilePath,
			Destination: &rk_gen.GenPbGoInfo.ConfigFilePath,
		},
	}
	err := cmd.Run(ctx)

	if err != nil {
		return err
	}

	return nil
}

func buildGoFiles() error {
	color.Cyan("Exporting GOOS=%s", BuildGoInfo.GoOs)
	exec.Command(fmt.Sprintf("GOOS=%s", BuildGoInfo.GoOs)).CombinedOutput()
	color.Green("[Done]")

	color.Cyan("Exporting GOARCH=%s", BuildGoInfo.GoArch)
	exec.Command(fmt.Sprintf("GOARCH=%s", BuildGoInfo.GoArch)).CombinedOutput()
	color.Green("[Done]")

	color.Cyan("Compile go files locate at internal/main.go")
	args := []string{
		"build",
		"-o",
		BuildGoInfo.AppName,
		"internal/main.go",
	}

	bytes, err := exec.Command("go", args...).CombinedOutput()
	if err != nil {
		color.Red("failed to compile go file\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	if BuildInfo.Debug {
		color.Blue("[stdout] %s", string(bytes))
	}

	color.Green("[Done]")

	if err != nil {
		return err
	}

	return nil
}

func createTargetFolder() error {
	color.Cyan("Clean and create target folder")

	err := os.RemoveAll("target")
	basePath := fmt.Sprintf("target/%s", BuildGoInfo.AppName)
	err = os.MkdirAll(basePath, os.ModePerm)
	err = os.MkdirAll(basePath+"/bin", os.ModePerm)
	err = os.MkdirAll(basePath+"/scripts", os.ModePerm)
	err = os.MkdirAll(basePath+"/log", os.ModePerm)
	err = os.MkdirAll(basePath+"/configs", os.ModePerm)
	err = os.MkdirAll(basePath+"/api", os.ModePerm)
	err = os.MkdirAll(basePath+"/assets", os.ModePerm)

	if err != nil {
		return err
		color.Red("[Failed]")
	}
	color.Green("[Done]")
	return nil
}

func runScript(path string) error {
	if isValidPath(path) {
		color.Cyan("Run %s script", path)

		bytes, err := exec.Command("/bin/sh", path).CombinedOutput()
		if err != nil {
			color.Red("[stderr] %s", string(bytes))
			color.Red("%v", err)
			return err
		}

		if len(bytes) > 0 {
			color.Blue("[stdout] %s", string(bytes))
		}
		return nil
	}

	return nil
}

func copyFiles(config *BuildGoConfig) error {
	basePath := fmt.Sprintf("target/%s", BuildGoInfo.AppName)

	// 5.1: copy configs/
	color.Cyan("Copy files in configs folder to target/%s/configs", BuildGoInfo.AppName)
	bytes, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("cp -r configs/* %s", basePath+"/configs")).CombinedOutput()

	if len(bytes) > 0 {
		color.Blue("[stdout] %s", string(bytes))
	}

	if err != nil {
		color.Red("[Failed]")
		return err
	}
	color.Green("[Done]")

	// 5.2: copy api/
	color.Cyan("Copy swagger json files to target/%s/api", BuildGoInfo.AppName)
	bytes, err = exec.Command("/bin/sh", "-c", fmt.Sprintf("cp -r api/* %s", basePath+"/api")).CombinedOutput()
	if len(bytes) > 0 {
		color.Blue("[stdout] %s", string(bytes))
	}
	if err != nil {
		color.Red("[Failed]")
		return err
	}

	// remove .proto and go files
	exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -rf %s", basePath+"/api/*/*.proto")).CombinedOutput()
	exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -rf %s", basePath+"/api/*/*.go")).CombinedOutput()

	color.Green("[Done]")

	// 5.3: copy bin/
	color.Cyan("Copy binary file to target/%s/bin", BuildGoInfo.AppName)
	bytes, err = exec.Command("mv", "./"+BuildGoInfo.AppName, basePath+"/bin").CombinedOutput()
	if len(bytes) > 0 {
		color.Blue("[stdout] %s", string(bytes))
	}
	if err != nil {
		color.Red("[Failed]")
		return err
	}

	color.Green("[Done]")

	// 5.4: copy assets/
	color.Cyan("Copy files in assets folder to target/%s/assets", BuildGoInfo.AppName)
	bytes, err = exec.Command("/bin/sh", "-c", fmt.Sprintf("cp -r assets/* %s", basePath+"/assets")).CombinedOutput()
	if len(bytes) > 0 {
		color.Blue("[stdout] %s", string(bytes))
	}
	if err != nil {
		color.Red("[Failed]")
		return err
	}

	color.Green("[Done]")

	// 5.5: copy scripts/
	color.Cyan("Copy files in assets folder to target/%s/scripts", BuildGoInfo.AppName)
	bytes, err = exec.Command("/bin/sh", "-c", fmt.Sprintf("cp -r scripts/* %s", basePath+"/scripts")).CombinedOutput()
	if len(bytes) > 0 {
		color.Blue("[stdout] %s", string(bytes))
	}
	if err != nil {
		color.Red("[Failed]")
		return err
	}

	exec.Command("/bin/sh", "-c", fmt.Sprintf("chmod +x %s", basePath+"/scripts/*")).CombinedOutput()

	color.Green("[Done]")

	// 5.6 Export .env file
	color.Cyan("Export env as target/%s/.env", BuildGoInfo.AppName)
	// create .env
	os.RemoveAll(basePath+"/.env")
	file, err := os.Create(basePath+"/.env")
	if err != nil {
		color.Red("[stderr] %v", err)
		color.Red("[Failed]")
		return err
	}
	if config != nil && config.Env != nil {
		for k, v := range config.Env {
			_, err := file.WriteString(k+"="+v+"\n")
			if err != nil {
				color.Red("[stderr] %v", err)
				color.Red("[Failed]")
				return err
			}
		}
	}

	file.WriteString("appName="+BuildGoInfo.AppName)
	file.Sync()

	color.Green("[Done]")
	return nil
}

func readBuildGoConfig() *BuildGoConfig {
	var config *BuildGoConfig
	if len(BuildGoInfo.ConfigFilePath) > 0 {
		config = fromConfigFileForGo()
	}

	if config == nil {
		return nil
	}

	// for goos
	BuildGoInfo.GoOs = runtime.GOOS
	if len(config.BuildGo.GoOs) > 0 {
		BuildGoInfo.GoOs = config.BuildGo.GoOs
	}

	// for goarch
	BuildGoInfo.GoArch = runtime.GOARCH
	if len(config.BuildGo.GoArch) > 0 {
		BuildGoInfo.GoArch = config.BuildGo.GoArch
	}

	// app name
	if len(config.AppName) > 0 {
		BuildGoInfo.AppName = config.AppName
	}

	return config
}

func fromConfigFileForGo() *BuildGoConfig {
	path := BuildGoInfo.ConfigFilePath
	if !isValidPath(path) {
		return nil
	}

	bytes, ext, err := readFile(path)

	if err != nil {
		return nil
	}

	if ext == ".yaml" || ext == ".yml" {
		// unmarshal yaml
		config := &BuildGoConfig{}
		if err = yaml.Unmarshal(bytes, config); err != nil {
			println(fmt.Sprintf("%v", err))
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