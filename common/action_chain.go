package rk_common

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

// Simple action chain which will print to console as below:
//
// [1/5] Action name 1
// [2/5] Action name 2
type actionChain struct {
	chain       []cli.ActionFunc
	description []string
	ignoreError []bool
}

func NewActionChain() *actionChain {
	return &actionChain{
		chain:       make([]cli.ActionFunc, 0),
		description: make([]string, 0),
		ignoreError: make([]bool, 0),
	}
}

func (c *actionChain) Add(description string, action cli.ActionFunc, ignoreError bool) {
	c.chain = append(c.chain, action)
	c.description = append(c.description, description)
	c.ignoreError = append(c.ignoreError, ignoreError)
}

func (c *actionChain) Execute(ctx *cli.Context) error {
	total := len(c.chain)

	for i := range c.chain {
		prefix := fmt.Sprintf("[%d/%d]", i+1, total)
		color.Green(fmt.Sprintf("%s %s", prefix, c.description[i]))

		if err := c.chain[i](ctx); err != nil {
			if !c.ignoreError[i] {
				event := GetEventV2(ctx)
				event.AddErr(err)
				event.SetResCode("Failed")
				return ErrorV2(err)
			}
		}
		color.Green("------------------------------------[OK]")
	}
	return nil
}
