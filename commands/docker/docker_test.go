// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package docker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDocker(t *testing.T) {
	assert.NotNil(t, Docker())
}
