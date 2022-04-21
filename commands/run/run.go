package run

import (
	"errors"
	"fmt"
	"github.com/rookie-ninja/rk/commands/build"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func Run() *cli.Command {
	command := &cli.Command{
		Name:      "run",
		Usage:     "Run server build by rk build",
		UsageText: "rk run",
		Action:    runAction,
	}

	return command
}

func runAction(ctx *cli.Context) error {
	if err := common.UnmarshalBootConfig("build.yaml", common.BuildConfig); err != nil {
		return err
	}

	chain := common.NewActionChain()
	chain.Add("Clearing target folder", func(ctx *cli.Context) error {
		// 0: Not dir of where go.mod file exists
		if !common.FileExists("go.mod") {
			return errors.New("not a go directory, failed to lookup go.mod file")
		}
		return os.RemoveAll(common.BuildTarget)
	}, false)

	switch common.BuildConfig.Build.Type {
	case "go":
		chain.Add("Execute user command before", build.ExecCommandsBefore, false)
		chain.Add("Execute user script before", build.ExecScriptBefore, false)
		chain.Add("Build go file", build.BuildGoFile, false)
		chain.Add("Copy to target folder", build.CopyToTarget, false)
		chain.Add("Execute user script after", build.ExecScriptAfter, false)
		chain.Add("Execute user command after", build.ExecCommandsAfter, false)
	default:
		chain.Add("Execute user command before", build.ExecCommandsBefore, false)
		chain.Add("Execute user script before", build.ExecScriptBefore, false)
		chain.Add("Copy to target folder", build.CopyToTarget, false)
		chain.Add("Execute user script after", build.ExecScriptAfter, false)
		chain.Add("Execute user command after", build.ExecCommandsAfter, false)
	}

	chain.Add("Run application", runBinary, false)

	return chain.Execute(ctx)
}

func runBinary(ctx *cli.Context) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	os.Chdir(common.BuildTarget)
	cmd := exec.Command(fmt.Sprintf("./bin/%s", common.TargetPkgName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		os.RemoveAll(common.BuildTarget)
		return err
	}

	go func() {
		<-sig
		cmd.Process.Signal(syscall.SIGINT)
	}()

	cmd.Wait()

	return nil
}
