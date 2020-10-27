package rk_uninstall

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

const EOE = "------------------------------------------------------------------"

// Uninstall rk-std on target hosts
func UninstallRkStdCommand() *cli.Command {
	command := &cli.Command{
		Name:      "rk-std",
		Usage:     "uninstall rk standard environment on local machine",
		UsageText: "rk uninstall rk-std",
		Action:    UninstallRkStdAction,
	}

	return command
}

func UninstallRkStdAction(ctx *cli.Context) error {
	// uninstall protobuf
	color.Magenta("[1] Uninstall protobuf")
	if err := UninstallProtobufAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// uninstall protoc-gen-go
	color.Magenta("[2] Uninstall protoc-gen-go")
	if err := UninstallProtocGenGoAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// uninstall protoc-gen-grpc-gateway
	color.Magenta("[3] Uninstall protoc-gen-grpc-gateway")
	if err := UninstallProtocGenGrpcGatewayAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// uninstall protoc-gen-doc
	color.Magenta("[4] Uninstall protoc-gen-doc")
	if err := UninstallProtocGenDocAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// uninstall protoc-gen-swagger
	color.Magenta("[5] Uninstall protoc-gen-swagger")
	if err := UninstallProtocGenSwaggerAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// uninstall mockgen
	color.Magenta("[6] Uninstall mockgen")
	if err := UninstallMockGenAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// uninstall golangci-lint
	color.Magenta("[7] Uninstall golangci-lint")
	if err := UninstallGoLangCILintAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// uninstall gocov
	color.Magenta("[8] Uninstall gocov")
	if err := UninstallGoCovAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// uninstall golint
	color.Magenta("[9] Uninstall golint")
	if err := UninstallGoLintAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// uninstall swag
	color.Magenta("[10] Uninstall swag")
	if err := UninstallSwagAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	return nil
}
