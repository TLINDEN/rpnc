/*
Copyright © 2023-2024 Thomas von Dein

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
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

type Calc struct {
	debug        bool
	batch        bool
	stdin        bool
	showstack    bool
	intermediate bool
	notdone      bool // set to true as long as there are items left in the eval loop
	precision    int

	stack        *Stack
	history      []string
	completer    readline.AutoCompleter
	interpreter  *Interpreter
	Space        *regexp.Regexp
	Comment      *regexp.Regexp
	Register     *regexp.Regexp
	Constants    []string
	LuaFunctions []string

	Funcalls      Funcalls
	BatchFuncalls Funcalls

	// different kinds of commands, displays nicer in help output
	StackCommands    Commands
	SettingsCommands Commands
	ShowCommands     Commands
	Commands         Commands

	Vars map[string]float64
}

// help for lua functions will be added dynamically
const Help string = `
Operators:
basic operators: + - x * / ^  (* is an alias of x)

Bitwise operators: and or xor < (left shift) > (right shift)

Percent functions:
%                    percent
%-                   subtract percent
%+                   add percent

Math functions (see https://pkg.go.dev/math):
mod sqrt abs acos acosh asin asinh atan atan2 atanh cbrt ceil cos cosh
erf erfc  erfcinv erfinv exp  exp2 expm1 floor  gamma ilogb j0  j1 log
log10 log1p log2 logb pow round roundtoeven sin sinh tan tanh trunc y0
y1 copysign dim hypot

Converter functions:
cm-to-inch              yards-to-meters         bytes-to-kilobytes
inch-to-cm              meters-to-yards         bytes-to-megabytes
gallons-to-liters       miles-to-kilometers     bytes-to-gigabytes
liters-to-gallons       kilometers-to-miles     bytes-to-terabytes

Batch functions:
sum                  sum of all values (alias: +)
max                  max of all values
min                  min of all values
mean                 mean of all values (alias: avg)
median               median of all values

Register variables:
>NAME                Put last stack element into variable NAME
<NAME                Retrieve variable NAME and put onto stack`

// commands, constants and operators,  defined here to feed completion
// and our mode switch in Eval() dynamically
const (
	Constants    string = `Pi Phi Sqrt2 SqrtE SqrtPi SqrtPhi Ln2 Log2E Ln10 Log10E`
	Precision    int    = 2
	ShowStackLen int    = 5
)

// That way we can add custom functions to completion
func GetCompleteCustomFunctions() func(string) []string {
	return func(line string) []string {
		completions := []string{}

		for luafunc := range LuaFuncs {
			completions = append(completions, luafunc)
		}

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

		for command := range c.SettingsCommands {
			if len(command) > 1 {
				completions = append(completions, command)
			}
		}

		for command := range c.ShowCommands {
			if len(command) > 1 {
				completions = append(completions, command)
			}
		}

		for command := range c.StackCommands {
			if len(command) > 1 {
				completions = append(completions, command)
			}
		}

		for command := range c.Commands {
			if len(command) > 1 {
				completions = append(completions, command)
			}
		}

		return completions
	}
}

func NewCalc() *Calc {
	calc := Calc{stack: NewStack(), debug: false, precision: Precision}

	calc.Funcalls = DefineFunctions()
	calc.BatchFuncalls = DefineBatchFunctions()
	calc.Vars = map[string]float64{}

	calc.completer = readline.NewPrefixCompleter(
		// custom lua functions
		readline.PcItemDynamic(GetCompleteCustomFunctions()),
		readline.PcItemDynamic(calc.GetCompleteCustomFuncalls()),
	)

	calc.Space = regexp.MustCompile(`\s+`)
	calc.Comment = regexp.MustCompile(`#.*`) // ignore everything after #
	calc.Register = regexp.MustCompile(`^([<>])([A-Z][A-Z0-9]*)`)

	// pre-calculate mode switching arrays
	calc.Constants = strings.Split(Constants, " ")

	calc.SetCommands()

	return &calc
}

// setup the interpreter, called from main(), import lua functions
func (c *Calc) SetInt(interpreter *Interpreter) {
	c.interpreter = interpreter

	for name := range LuaFuncs {
		c.LuaFunctions = append(c.LuaFunctions, name)
	}
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
	prompt := "\033[31m»\033[0m "
	batch := ""

	if c.batch {
		batch = "->batch"
	}

	debug := ""
	revision := ""

	if c.debug {
		debug = "->debug"
		revision = fmt.Sprintf("/rev%d", c.stack.rev)
	}

	return fmt.Sprintf("rpn%s%s [%d%s]%s", batch, debug, c.stack.Len(), revision, prompt)
}

// the actual work horse, evaluate a line of calc command[s]
func (c *Calc) Eval(line string) error {
	// remove surrounding whitespace and comments, if any
	line = strings.TrimSpace(c.Comment.ReplaceAllString(line, ""))

	if line == "" {
		return nil
	}

	items := c.Space.Split(line, -1)

	for pos, item := range items {
		if pos+1 < len(items) {
			c.notdone = true
		} else {
			c.notdone = false
		}

		if err := c.EvalItem(item); err != nil {
			return err
		}
	}

	if c.showstack && !c.stdin {
		dots := ""

		if c.stack.Len() > ShowStackLen {
			dots = "... "
		}

		last := c.stack.Last(ShowStackLen)

		fmt.Printf("stack: %s%s\n", dots, list2str(last))
	}

	return nil
}

func (c *Calc) EvalItem(item string) error {
	num, err := strconv.ParseFloat(item, 64)

	if err == nil {
		c.stack.Backup()
		c.stack.Push(num)

		return nil
	}

	// try time
	var hour, min int
	_, err = fmt.Sscanf(item, "%d:%d", &hour, &min)
	if err == nil {
		c.stack.Backup()
		c.stack.Push(float64(hour) + float64(min)/60)

		return nil
	}

	// try hex
	var i int
	_, err = fmt.Sscanf(item, "0x%x", &i)
	if err == nil {
		c.stack.Backup()
		c.stack.Push(float64(i))

		return nil
	}

	if contains(c.Constants, item) {
		// put the constant onto the stack
		c.stack.Backup()
		c.stack.Push(const2num(item))

		return nil
	}

	if exists(c.Funcalls, item) {
		if err := c.DoFuncall(item); err != nil {
			return Error(err.Error())
		}

		c.Result()

		return nil
	}

	if exists(c.BatchFuncalls, item) {
		if !c.batch {
			return Error("only supported in batch mode")
		}

		if err := c.DoFuncall(item); err != nil {
			return Error(err.Error())
		}

		c.Result()

		return nil
	}

	if contains(c.LuaFunctions, item) {
		// user provided custom lua functions
		c.EvalLuaFunction(item)

		return nil
	}

	regmatches := c.Register.FindStringSubmatch(item)
	if len(regmatches) == 3 {
		switch regmatches[1] {
		case ">":
			c.PutVar(regmatches[2])
		case "<":
			c.GetVar(regmatches[2])
		}

		return nil
	}

	// internal commands
	// FIXME: propagate errors
	if exists(c.Commands, item) {
		c.Commands[item].Func(c)

		return nil
	}

	if exists(c.ShowCommands, item) {
		c.ShowCommands[item].Func(c)

		return nil
	}

	if exists(c.StackCommands, item) {
		c.StackCommands[item].Func(c)

		return nil
	}

	if exists(c.SettingsCommands, item) {
		c.SettingsCommands[item].Func(c)

		return nil
	}

	switch item {
	case "?", "help":
		c.PrintHelp()

	default:
		return Error("unknown command or operator")
	}

	return nil
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
		return Error("function not defined but in completion list")
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
	funcresult := function.Func(args)

	if funcresult.Err != nil {
		// leave the stack untouched in case of any error
		return funcresult.Err
	}

	// don't forget to backup!
	c.stack.Backup()

	// "pop"
	if batch {
		// get rid of stack
		c.stack.Clear()
	} else {
		// remove operands
		c.stack.Shift(function.Expectargs)
	}

	// save result
	c.stack.Push(funcresult.Res)

	// thanks a lot
	c.SetHistory(funcname, args, funcresult.Res)

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
	// we only  print the result if it's either  a final result or
	// (if it is intermediate) if -i has been given
	if c.intermediate || !c.notdone {
		// only needed in repl
		if !c.stdin {
			fmt.Print("= ")
		}

		result := c.stack.Last()[0]
		truncated := math.Trunc(result)
		precision := c.precision

		if result == truncated {
			precision = 0
		}

		format := fmt.Sprintf("%%.%df\n", precision)
		fmt.Printf(format, result)
	}

	return c.stack.Last()[0]
}

func (c *Calc) Debug(msg string) {
	if c.debug {
		fmt.Printf("DEBUG(calc): %s\n", msg)
	}
}

func (c *Calc) EvalLuaFunction(funcname string) {
	// called from calc loop
	var luaresult float64

	var err error

	switch c.interpreter.FuncNumArgs(funcname) {
	case 0:
		fallthrough
	case 1:
		luaresult, err = c.interpreter.CallLuaFunc(funcname, c.stack.Last())
	case 2:
		luaresult, err = c.interpreter.CallLuaFunc(funcname, c.stack.Last(2))
	case -1:
		luaresult, err = c.interpreter.CallLuaFunc(funcname, c.stack.All())
	default:
		luaresult, err = 0, errors.New("invalid number of argument requested")
	}

	if err != nil {
		fmt.Println(err)

		return
	}

	c.stack.Backup()

	dopush := true

	switch c.interpreter.FuncNumArgs(funcname) {
	case 0:
		a := c.stack.Last()

		if len(a) == 1 {
			c.History("%s(%f) = %f", funcname, a, luaresult)
		}

		dopush = false
	case 1:
		a := c.stack.Pop()
		c.History("%s(%f) = %f", funcname, a, luaresult)
	case 2:
		a := c.stack.Pop()
		b := c.stack.Pop()
		c.History("%s(%f,%f) = %f", funcname, a, b, luaresult)
	case -1:
		c.stack.Clear()
		c.History("%s(*) = %f", funcname, luaresult)
	}

	if dopush {
		c.stack.Push(luaresult)
	}

	c.Result()
}

func (c *Calc) PutVar(name string) {
	last := c.stack.Last()

	if len(last) == 1 {
		c.Debug(fmt.Sprintf("register %.2f in %s", last[0], name))
		c.Vars[name] = last[0]
	} else {
		fmt.Println("empty stack")
	}
}

func (c *Calc) GetVar(name string) {
	if exists(c.Vars, name) {
		c.Debug(fmt.Sprintf("retrieve %.2f from %s", c.Vars[name], name))
		c.stack.Backup()
		c.stack.Push(c.Vars[name])
	} else {
		fmt.Println("variable doesn't exist")
	}
}

func sortcommands(hash Commands) []string {
	keys := make([]string, 0, len(hash))

	for key := range hash {
		if len(key) > 1 {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	return keys
}

func (c *Calc) PrintHelp() {
	output := "Available configuration commands:\n"

	for _, name := range sortcommands(c.SettingsCommands) {
		output += fmt.Sprintf("%-20s %s\n", name, c.SettingsCommands[name].Help)
	}

	output += "\nAvailable show commands:\n"

	for _, name := range sortcommands(c.ShowCommands) {
		output += fmt.Sprintf("%-20s %s\n", name, c.ShowCommands[name].Help)
	}

	output += "\nAvailable stack manipulation commands:\n"

	for _, name := range sortcommands(c.StackCommands) {
		output += fmt.Sprintf("%-20s %s\n", name, c.StackCommands[name].Help)
	}

	output += "\nOther commands:\n"

	for _, name := range sortcommands(c.Commands) {
		output += fmt.Sprintf("%-20s %s\n", name, c.Commands[name].Help)
	}

	output += "\n" + Help

	// append lua functions, if any
	if len(LuaFuncs) > 0 {
		output += "\nLua functions:\n"

		for name, function := range LuaFuncs {
			output += fmt.Sprintf("%-20s %s\n", name, function.help)
		}
	}

	Pager("rpn help overview", output)
}
