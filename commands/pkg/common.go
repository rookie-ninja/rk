// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package pkg

import (
	"github.com/urfave/cli/v2"
	"strings"
)

func getPkgTypeFromFlags(ctx *cli.Context) string {
	return strings.ToLower(ctx.Value("type").(string))
}

func getPkgNameFromFlags(ctx *cli.Context) string {
	return ctx.Value("name").(string)
}

func getPkgDestinationFromFlags(ctx *cli.Context) string {
	return ctx.Value("destination").(string)
}

func getPkgFromRemoteSourceFromFlags(ctx *cli.Context) bool {
	return ctx.Value("remote-source").(bool)
}
