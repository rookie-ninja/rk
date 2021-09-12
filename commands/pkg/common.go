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
	//wd, _ := os.Getwd()
	dest := ctx.Value("destination").(string)
	//if strings.HasPrefix(dest,".") {
	//	dest = path.Join(wd, strings.TrimPrefix(dest, "."))
	//} else if !path.IsAbs(dest) {
	//	dest = path.Join(wd, dest)
	//}

	return dest
}

func getPkgFromRemoteSourceFromFlags(ctx *cli.Context) bool {
	return ctx.Value("remote-source").(bool)
}
