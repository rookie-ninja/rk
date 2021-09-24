// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package common

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"reflect"
	"testing"
)

func TestNewActionChain(t *testing.T) {
	chain := NewActionChain()
	assert.NotNil(t, chain)
	assert.NotNil(t, chain.chain)
	assert.NotNil(t, chain.description)
	assert.NotNil(t, chain.ignoreError)
}

func TestActionChain_Add(t *testing.T) {
	chain := NewActionChain()
	desc := "ut-description"
	f := fakeAction
	chain.Add(desc, f, false)

	assert.Equal(t, desc, chain.description[0])
	assert.Equal(t, reflect.ValueOf(f).Pointer(), reflect.ValueOf(chain.chain[0]).Pointer())
	assert.False(t, chain.ignoreError[0])
}

func TestActionChain_Execute(t *testing.T) {
	// Do not ignore error
	chain := NewActionChain()
	chain.Add("ut-description", func(context *cli.Context) error {
		return errors.New("ut-error")
	}, false)

	err := chain.Execute(cli.NewContext(nil, nil, nil))
	assert.NotNil(t, err)

	// Ignore error
	chain = NewActionChain()
	chain.Add("ut-description", func(context *cli.Context) error {
		return errors.New("ut-error")
	}, true)
	err = chain.Execute(cli.NewContext(nil, nil, nil))
	assert.Nil(t, err)
}

func fakeAction(ctx *cli.Context) error {
	return nil
}
