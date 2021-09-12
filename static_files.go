// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

// Package rk is a static files which aims to include assets files into pkger
package rk

import "github.com/markbates/pkger"

func init() {
	pkger.Include("/commands/pkg/assets/")
}
