// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package clear

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClear(t *testing.T) {
	assert.NotNil(t, Clear())
}

func TestClearAction(t *testing.T) {
	assert.Nil(t, clearAction(nil))
}
