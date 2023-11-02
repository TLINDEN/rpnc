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
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

type Calc struct {
	debug          bool
	batch          bool
	stdin          bool
	stack          *Stack
	history        []string
	completer      readline.AutoCompleter
	interpreter    *Interpreter
	Operators      *regexp.Regexp
	Space          *regexp.Regexp
	Constants      []string
	MathFunctions  []string
	BatchFunctions []string
	LuaFunctions   []string
}

// help for lua functions will be added dynamically
const Help string = `Available commands:
batch                toggle batch mode
debug                toggle debug output
dump                 display the stack contents
clear                clear the whole stack
shift                remove the last element of the stack
history              display calculation history
help|?               show this message
quit|exit|c-d|c-c    exit program

Available operators:
basic operators: + - x /

Available math functions:
sqrt                 square root
mod                  remainder of division (alias: remainder)
max                  batch mode only: max of all values
min                  batch mode only: min of all values
mean                 batch mode only: mean of all values (alias: avg)
median               batch mode only: median of all values
%                    percent
%-                   substract percent
%+                   add percent

Math operators:
^                    power`

// commands, constants and operators,  defined here to feed completion
// and our mode switch in Eval() dynamically
const (
	Commands       string = `dump reverse debug undebug clear batch shift undo help history manual exit quit`
	Operators      string = `+ - * x / ^ % %- %+`
	MathFunctions  string = `sqrt remainder`
	Constants      string = `Pi Phi Sqrt2 SqrtE SqrtPi SqrtPhi Ln2 Log2E Ln10 Log10E`
	BatchFunctions string = `median avg mean max min`
)

// That way we can add custom functions to completion
func GetCompleteCustomFunctions() func(string) []string {
	return func(line string) []string {
		completions := []string{}

		for luafunc := range LuaFuncs {
			completions = append(completions, luafunc)
		}

		completions = append(completions, strings.Split(Commands, " ")...)
		completions = append(completions, strings.Split(Operators, " ")...)
		completions = append(completions, strings.Split(MathFunctions, " ")...)
		completions = append(completions, strings.Split(Constants, " ")...)

		return completions
	}
}

func NewCalc() *Calc {
	c := Calc{stack: NewStack(), debug: false}

	c.completer = readline.NewPrefixCompleter(
		// custom lua functions
		readline.PcItemDynamic(GetCompleteCustomFunctions()),
	)

	// pre-calculate mode switching regexes
	reg := `^[`
	for _, op := range strings.Split(Operators, " ") {
		switch op {
		case "x":
			reg += op
		default:
			reg += `\` + op
		}
	}
	reg += `]$`
	c.Operators = regexp.MustCompile(reg)

	c.Space = regexp.MustCompile(`\s+`)

	// pre-calculate mode switching arrays
	c.Constants = strings.Split(Constants, " ")
	c.MathFunctions = strings.Split(MathFunctions, " ")
	c.BatchFunctions = strings.Split(BatchFunctions, " ")

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
	line = strings.TrimSpace(line)

	if line == "" {
		return
	}

	for _, item := range c.Space.Split(line, -1) {
		num, err := strconv.ParseFloat(item, 64)

		if err == nil {
			c.stack.Backup()
			c.stack.Push(num)
		} else {
			if c.Operators.MatchString(item) {
				// simple ops like + or x
				c.simple(item[0])
				continue
			}

			if contains(c.Constants, item) {
				// put the constant onto the stack
				c.stack.Backup()
				c.stack.Push(const2num(item))
				continue
			}

			if contains(c.MathFunctions, item) {
				// go builtin math function, if implemented
				c.mathfunc(item)
				continue
			}

			if contains(c.BatchFunctions, item) {
				// math functions only supported in batch mode like max or mean
				c.batchfunc(item)
				continue
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
			case "undo":
				c.stack.Restore()
			case "history":
				for _, entry := range c.history {
					fmt.Println(entry)
				}
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

	fmt.Println(c.stack.Last())

	return c.stack.Last()
}

func (c *Calc) Debug(msg string) {
	if c.debug {
		fmt.Printf("DEBUG(calc): %s\n", msg)
	}
}

// do simple calculations
func (c *Calc) simple(op byte) {
	c.stack.Backup()

	for c.stack.Len() > 1 {
		b := c.stack.Pop()
		a := c.stack.Pop()
		var x float64

		c.Debug(fmt.Sprintf("evaluating: %.2f %c %.2f", a, op, b))

		switch op {
		case '+':
			x = a + b
		case '-':
			x = a - b
		case 'x':
			fallthrough // alias for *
		case '*':
			x = a * b
		case '/':
			if b == 0 {
				fmt.Println("error: division by null!")
				return
			}
			x = a / b
		case '^':
			x = math.Pow(a, b)
		default:
			panic("invalid operator!")
		}

		c.stack.Push(x)

		c.History("%f %c %f = %f", a, op, b, x)

		if !c.batch {
			break
		}
	}

	c.Result()
}

// calc using go math lib functions
func (c *Calc) mathfunc(funcname string) {
	c.stack.Backup()

	for c.stack.Len() > 0 {
		var x float64

		switch funcname {
		case "sqrt":
			a := c.stack.Pop()
			x = math.Sqrt(a)
			c.History("sqrt(%f) = %f", a, x)

		case "mod":
			fallthrough // alias
		case "remainder":
			b := c.stack.Pop()
			a := c.stack.Pop()
			x = math.Remainder(a, b)
			c.History("remainderf(%f / %f) = %f", a, b, x)

		case "%":
			b := c.stack.Pop()
			a := c.stack.Pop()

			x = (a / 100) * b
			c.History("%f percent of %f = %f", b, a, x)

		case "%-":
			b := c.stack.Pop()
			a := c.stack.Pop()

			x = a - ((a / 100) * b)
			c.History("%f minus %f percent of %f = %f", a, b, a, x)

		case "%+":
			b := c.stack.Pop()
			a := c.stack.Pop()

			x = a + ((a / 100) * b)
			c.History("%f plus %f percent of %f = %f", a, b, a, x)

		}

		c.stack.Push(x)

		if !c.batch {
			break
		}
	}

	c.Result()
}

// execute pure batch functions, operating on the whole stack
func (c *Calc) batchfunc(funcname string) {
	if !c.batch {
		fmt.Println("error: only available in batch mode")
	}

	c.stack.Backup()
	var x float64
	count := c.stack.Len()

	switch funcname {
	case "median":
		all := []float64{}

		for c.stack.Len() > 0 {
			all = append(all, c.stack.Pop())
		}

		middle := count / 2

		x = all[middle]
		c.History("median(all)")

	case "mean":
		fallthrough // alias
	case "avg":
		var sum float64

		for c.stack.Len() > 0 {
			sum += c.stack.Pop()
		}

		x = sum / float64(count)
		c.History("avg(all)")
	case "min":
		x = c.stack.Pop() // initialize with the last one

		for c.stack.Len() > 0 {
			val := c.stack.Pop()
			if val < x {
				x = val
			}
		}

		c.History("min(all)")
	case "max":
		x = c.stack.Pop() // initialize with the last one

		for c.stack.Len() > 0 {
			val := c.stack.Pop()
			if val > x {
				x = val
			}
		}

		c.History("max(all)")
	}

	c.stack.Push(x)
	_ = c.Result()
}

func (c *Calc) luafunc(funcname string) {
	// called from calc loop
	var x float64
	var err error

	switch c.interpreter.FuncNumArgs(funcname) {
	case 1:
		x, err = c.interpreter.CallLuaFunc(funcname, []float64{c.stack.Last()})
	case 2:
		x, err = c.interpreter.CallLuaFunc(funcname, c.stack.LastTwo())
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
