/*
Copyright © 2023 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"os"
)

type CommandFunction func(*Calc)

type Command struct {
	Help string
	Func CommandFunction
}

type Commands map[string]*Command

func NewCommand(help string, function CommandFunction) *Command {
	return &Command{
		Help: help,
		Func: function,
	}
}

// define all management (that is: non calculation) commands
func (c *Calc) SetCommands() {
	f := Commands{
		// Toggles
		"debug": NewCommand(
			"toggle debugging",
			func(c *Calc) {
				c.ToggleDebug()
			},
		),

		"nodebug": NewCommand(
			"disable debugging",
			func(c *Calc) {
				c.debug = false
				c.stack.debug = false
			},
		),

		"batch": NewCommand(
			"toggle batch mode",
			func(c *Calc) {
				c.ToggleBatch()
			},
		),

		"nobatch": NewCommand(
			"disable batch mode",
			func(c *Calc) {
				c.batch = false
			},
		),

		"showstack": NewCommand(
			"toggle show last 5 items of the stack",
			func(c *Calc) {
				c.ToggleShow()
			},
		),

		"noshowstack": NewCommand(
			"disable display of the stack",
			func(c *Calc) {
				c.showstack = false
			},
		),

		// Display commands
		"dump": NewCommand(
			"display the stack contents",
			func(c *Calc) {
				c.stack.Dump()
			},
		),

		"history": NewCommand(
			"display calculation history",
			func(c *Calc) {
				for _, entry := range c.history {
					fmt.Println(entry)
				}
			},
		),

		"vars": NewCommand(
			"show list of variables",
			func(c *Calc) {
				if len(c.Vars) > 0 {
					fmt.Printf("%-20s     %s\n", "VARIABLE", "VALUE")
					for k, v := range c.Vars {
						fmt.Printf("%-20s  -> %.2f\n", k, v)
					}
				} else {
					fmt.Println("no vars registered")
				}
			},
		),

		// stack manipulation commands
		"clear": NewCommand(
			"clear the whole stack",
			func(c *Calc) {
				c.stack.Backup()
				c.stack.Clear()
			},
		),

		"shift": NewCommand(
			"remove the last element of the stack",
			func(c *Calc) {
				c.stack.Backup()
				c.stack.Shift()
			},
		),

		"reverse": NewCommand(
			"reverse the stack elements",
			func(c *Calc) {
				c.stack.Backup()
				c.stack.Reverse()
			},
		),

		"swap": NewCommand(
			"exchange the last two elements",
			func(c *Calc) {
				if c.stack.Len() < 2 {
					fmt.Println("stack too small, can't swap")
				} else {
					c.stack.Backup()
					c.stack.Swap()
				}
			},
		),

		"undo": NewCommand(
			"undo last operation",
			func(c *Calc) {
				c.stack.Restore()
			},
		),

		"dup": NewCommand(
			"duplicate last stack item",
			func(c *Calc) {
				item := c.stack.Last()
				if len(item) == 1 {
					c.stack.Backup()
					c.stack.Push(item[0])
				} else {
					fmt.Println("stack empty")
				}
			},
		),

		// general commands
		"exit": NewCommand(
			"exit program",
			func(c *Calc) {
				os.Exit(0)
			},
		),

		"manual": NewCommand(
			"show manual",
			func(c *Calc) {
				man()
			},
		),
	}

	// aliases
	f["quit"] = f["exit"]
	f["undebug"] = f["nodebug"]
	f["show"] = f["showstack"]

	c.Commands = f
}
