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
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

func (c *Calc) SetSettingsCommands() Commands {
	return Commands{
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
	}
}

func (c *Calc) SetShowCommands() Commands {
	return Commands{
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

		"hex": NewCommand(
			"show last stack item in hex form (converted to int)",
			func(c *Calc) {
				if c.stack.Len() > 0 {
					fmt.Printf("0x%x\n", int(c.stack.Last()[0]))
				}
			},
		),
	}
}

func (c *Calc) SetStackCommands() Commands {
	return Commands{
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
			CommandSwap,
		),

		"undo": NewCommand(
			"undo last operation",
			func(c *Calc) {
				c.stack.Restore()
			},
		),

		"dup": NewCommand(
			"duplicate last stack item",
			CommandDup,
		),

		"edit": NewCommand(
			"edit the stack interactively",
			CommandEdit,
		),
	}
}

// define all management (that is: non calculation) commands
func (c *Calc) SetCommands() {
	c.SettingsCommands = c.SetSettingsCommands()
	c.ShowCommands = c.SetShowCommands()
	c.StackCommands = c.SetStackCommands()

	// general commands
	c.Commands = Commands{
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
	c.Commands["quit"] = c.Commands["exit"]

	c.SettingsCommands["d"] = c.SettingsCommands["debug"]
	c.SettingsCommands["b"] = c.SettingsCommands["batch"]
	c.SettingsCommands["s"] = c.SettingsCommands["showstack"]

	c.ShowCommands["h"] = c.ShowCommands["history"]
	c.ShowCommands["p"] = c.ShowCommands["dump"]
	c.ShowCommands["v"] = c.ShowCommands["vars"]

	c.StackCommands["c"] = c.StackCommands["clear"]
	c.StackCommands["u"] = c.StackCommands["undo"]
}

// added to the command map:
func CommandSwap(c *Calc) {
	if c.stack.Len() < 2 {
		fmt.Println("stack too small, can't swap")
	} else {
		c.stack.Backup()
		c.stack.Swap()
	}
}

func CommandDup(c *Calc) {
	item := c.stack.Last()
	if len(item) == 1 {
		c.stack.Backup()
		c.stack.Push(item[0])
	} else {
		fmt.Println("stack empty")
	}
}

func CommandEdit(calc *Calc) {
	if calc.stack.Len() == 0 {
		fmt.Println("empty stack")

		return
	}

	calc.stack.Backup()

	// put the stack contents into a tmp file
	tmp, err := os.CreateTemp("", "stack")
	if err != nil {
		fmt.Println(err)

		return
	}

	defer os.Remove(tmp.Name())

	comment := `# add or remove numbers as you wish.
# each number must be on its own line.
# numbers must be floating point formatted.
`
	_, err = tmp.WriteString(comment)

	if err != nil {
		fmt.Println(err)

		return
	}

	for _, item := range calc.stack.All() {
		_, err = fmt.Fprintf(tmp, "%f\n", item)
		if err != nil {
			fmt.Println(err)

			return
		}
	}

	tmp.Close()

	// determine which editor to use
	editor := "vi"

	enveditor, present := os.LookupEnv("EDITOR")
	if present {
		if editor != "" {
			if _, err := os.Stat(editor); err == nil {
				editor = enveditor
			}
		}
	}

	// execute editor with our tmp file containing current stack
	cmd := exec.Command(editor, tmp.Name())

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("could not run editor command: ", err)

		return
	}

	// read the file back in
	modified, err := os.Open(tmp.Name())
	if err != nil {
		fmt.Println("Error opening file:", err)

		return
	}
	defer modified.Close()

	// reset the stack
	calc.stack.Clear()

	// and put the new contents (if legit) back onto the stack
	scanner := bufio.NewScanner(modified)
	for scanner.Scan() {
		line := strings.TrimSpace(calc.Comment.ReplaceAllString(scanner.Text(), ""))
		if line == "" {
			continue
		}

		num, err := strconv.ParseFloat(line, 64)
		if err != nil {
			fmt.Printf("%s is not a floating point number!\n", line)

			continue
		}

		calc.stack.Push(num)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from file:", err)
	}
}
