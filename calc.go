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
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

type Calc struct {
	debug        bool
	batch        bool
	stdin        bool
	showstack    bool
	stack        *Stack
	history      []string
	completer    readline.AutoCompleter
	interpreter  *Interpreter
	Space        *regexp.Regexp
	Comment      *regexp.Regexp
	Constants    []string
	LuaFunctions []string

	Funcalls      Funcalls
	BatchFuncalls Funcalls
}

// help for lua functions will be added dynamically
const Help string = `Available commands:
batch                toggle batch mode
debug                toggle debug output
show                 show the last 5 items of the stack
dump                 display the stack contents
clear                clear the whole stack
shift                remove the last element of the stack
reverse              reverse the stack elements
swap                 exchange the last two elements
history              display calculation history
help|?               show this message
quit|exit|c-d|c-c    exit program

Operators:
basic operators: + - x * / ^  (* is an alias of x)

Percent functions:
%                    percent
%-                   substract percent
%+                   add percent

Math functions (see https://pkg.go.dev/math):
mod sqrt abs acos acosh asin asinh atan atan2 atanh cbrt ceil cos cosh
erf erfc  erfcinv erfinv exp  exp2 expm1 floor  gamma ilogb j0  j1 log
log10 log1p log2 logb pow round roundtoeven sin sinh tan tanh trunc y0
y1 copysign dim hypot

Batch functions:
sum                  sum of all values (alias: +)
max                  max of all values
min                  min of all values
mean                 mean of all values (alias: avg)
median               median of all values`

// commands, constants and operators,  defined here to feed completion
// and our mode switch in Eval() dynamically
const (
	Commands  string = `dump reverse debug undebug clear batch shift undo help history manual exit quit swap show`
	Constants string = `Pi Phi Sqrt2 SqrtE SqrtPi SqrtPhi Ln2 Log2E Ln10 Log10E`
)

// That way we can add custom functions to completion
func GetCompleteCustomFunctions() func(string) []string {
	return func(line string) []string {
		completions := []string{}

		for luafunc := range LuaFuncs {
			completions = append(completions, luafunc)
		}

		completions = append(completions, strings.Split(Commands, " ")...)
		completions = append(completions, strings.Split(Constants, " ")...)

		return completions
	}
}

func (c *Calc) GetCompleteCustomFuncalls() func(string) []string {
	return func(line string) []string {
		completions := []string{}

		for function := range c.Funcalls {
			completions = append(completions, function)
		}

		for function := range c.BatchFuncalls {
			completions = append(completions, function)
		}

		return completions
	}

}

func NewCalc() *Calc {
	c := Calc{stack: NewStack(), debug: false}

	c.Funcalls = DefineFunctions()
	c.BatchFuncalls = DefineBatchFunctions()

	c.completer = readline.NewPrefixCompleter(
		// custom lua functions
		readline.PcItemDynamic(GetCompleteCustomFunctions()),
		readline.PcItemDynamic(c.GetCompleteCustomFuncalls()),
	)

	c.Space = regexp.MustCompile(`\s+`)
	c.Comment = regexp.MustCompile(`#.*`) // ignore everything after #

	// pre-calculate mode switching arrays
	c.Constants = strings.Split(Constants, " ")

	for name := range LuaFuncs {
		c.LuaFunctions = append(c.LuaFunctions, name)
	}

	return &c
}

// setup the interpreter, called from main()
func (c *Calc) SetInt(I *Interpreter) {
	c.interpreter = I
}

func (c *Calc) ToggleDebug() {
	c.debug = !c.debug
	c.stack.ToggleDebug()
	fmt.Printf("debugging set to %t\n", c.debug)
}

func (c *Calc) ToggleBatch() {
	c.batch = !c.batch
	fmt.Printf("batchmode set to %t\n", c.batch)
}

func (c *Calc) ToggleStdin() {
	c.stdin = !c.stdin
}

func (c *Calc) ToggleShow() {
	c.showstack = !c.showstack
}

func (c *Calc) Prompt() string {
	p := "\033[31m»\033[0m "
	b := ""

	if c.batch {
		b = "->batch"
	}

	d := ""
	v := ""

	if c.debug {
		d = "->debug"
		v = fmt.Sprintf("/rev%d", c.stack.rev)
	}

	return fmt.Sprintf("rpn%s%s [%d%s]%s", b, d, c.stack.Len(), v, p)
}

// the actual work horse, evaluate a line of calc command[s]
func (c *Calc) Eval(line string) {
	// remove surrounding whitespace and comments, if any
	line = strings.TrimSpace(c.Comment.ReplaceAllString(line, ""))

	if line == "" {
		return
	}

	for _, item := range c.Space.Split(line, -1) {
		num, err := strconv.ParseFloat(item, 64)

		if err == nil {
			c.stack.Backup()
			c.stack.Push(num)
		} else {
			if contains(c.Constants, item) {
				// put the constant onto the stack
				c.stack.Backup()
				c.stack.Push(const2num(item))
				continue
			}

			if _, ok := c.Funcalls[item]; ok {
				if err := c.DoFuncall(item); err != nil {
					fmt.Println(err)
				} else {
					c.Result()
				}
				continue
			}

			if c.batch {
				if _, ok := c.BatchFuncalls[item]; ok {
					if err := c.DoFuncall(item); err != nil {
						fmt.Println(err)
					} else {
						c.Result()
					}
					continue
				}
			} else {
				if _, ok := c.BatchFuncalls[item]; ok {
					fmt.Println("only supported in batch mode")
					continue
				}
			}

			if contains(c.LuaFunctions, item) {
				// user provided custom lua functions
				c.luafunc(item)
				continue
			}

			// management commands
			switch item {
			case "?":
				fallthrough
			case "help":
				fmt.Println(Help)
				fmt.Println("Lua functions:")
				for name, function := range LuaFuncs {
					fmt.Printf("%-20s %s\n", name, function.help)
				}
			case "dump":
				c.stack.Dump()
			case "debug":
				c.ToggleDebug()
			case "undebug":
				c.debug = false
			case "batch":
				c.ToggleBatch()
			case "clear":
				c.stack.Backup()
				c.stack.Clear()
			case "shift":
				c.stack.Backup()
				c.stack.Shift()
			case "reverse":
				c.stack.Backup()
				c.stack.Reverse()
			case "swap":
				if c.stack.Len() < 2 {
					fmt.Println("stack too small, can't swap")
				} else {
					c.stack.Backup()
					c.stack.Swap()
				}
			case "undo":
				c.stack.Restore()
			case "history":
				for _, entry := range c.history {
					fmt.Println(entry)
				}
			case "show":
				c.ToggleShow()
			case "exit":
				fallthrough
			case "quit":
				os.Exit(0)
			case "manual":
				man()
			default:
				fmt.Println("unknown command or operator!")
			}
		}
	}

	if c.showstack && !c.stdin {
		dots := ""

		if c.stack.Len() > 5 {
			dots = "... "
		}
		last := c.stack.Last(5)
		fmt.Printf("stack: %s%s\n", dots, list2str(last))
	}
}

// Execute a math function, check if it is defined just in case
func (c *Calc) DoFuncall(funcname string) error {
	var function *Funcall
	if c.batch {
		function = c.BatchFuncalls[funcname]
	} else {
		function = c.Funcalls[funcname]
	}

	if function == nil {
		panic("function not defined but in completion list")
	}

	var args Numbers
	batch := false

	if function.Expectargs == -1 {
		// batch mode, but always < stack len, so check first
		args = c.stack.All()
		batch = true
	} else {
		//  this is way better behavior than just using 0 in place of
		// non-existing stack items
		if c.stack.Len() < function.Expectargs {
			return errors.New("stack doesn't provide enough arguments")
		}

		args = c.stack.Last(function.Expectargs)
	}

	c.Debug(fmt.Sprintf("calling %s with args: %v", funcname, args))

	// the  actual lambda call, so  to say. We provide  a slice of
	// the requested size, fetched  from the stack (but not popped
	// yet!)
	R := function.Func(args)

	if R.Err != nil {
		// leave the stack untouched in case of any error
		return R.Err
	}

	if batch {
		// get rid of stack
		c.stack.Clear()
	} else {
		// remove operands
		c.stack.Shift(function.Expectargs)
	}

	// save result
	c.stack.Push(R.Res)

	// thanks a lot
	c.SetHistory(funcname, args, R.Res)
	return nil
}

// we need to add a history entry for each operation
func (c *Calc) SetHistory(op string, args Numbers, res float64) {
	c.History("%s %s -> %f", list2str(args), op, res)
}

// just a textual representation of math operations, viewable with the
// history command
func (c *Calc) History(format string, args ...any) {
	c.history = append(c.history, fmt.Sprintf(format, args...))
}

// print the result
func (c *Calc) Result() float64 {
	if !c.stdin {
		fmt.Print("= ")
	}

	fmt.Println(c.stack.Last()[0])

	return c.stack.Last()[0]
}

func (c *Calc) Debug(msg string) {
	if c.debug {
		fmt.Printf("DEBUG(calc): %s\n", msg)
	}
}

func (c *Calc) luafunc(funcname string) {
	// called from calc loop
	var x float64
	var err error

	switch c.interpreter.FuncNumArgs(funcname) {
	case 1:
		x, err = c.interpreter.CallLuaFunc(funcname, c.stack.Last())
	case 2:
		x, err = c.interpreter.CallLuaFunc(funcname, c.stack.Last(2))
	case -1:
		x, err = c.interpreter.CallLuaFunc(funcname, c.stack.All())
	default:
		x, err = 0, errors.New("invalid number of argument requested")
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	c.stack.Backup()

	switch c.interpreter.FuncNumArgs(funcname) {
	case 1:
		a := c.stack.Pop()
		c.History("%s(%f) = %f", funcname, a, x)
	case 2:
		a := c.stack.Pop()
		b := c.stack.Pop()
		c.History("%s(%f,%f) = %f", funcname, a, b, x)
	case -1:
		c.stack.Clear()
		c.History("%s(*) = %f", funcname, x)
	}

	c.stack.Push(x)

	c.Result()
}
