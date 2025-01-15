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

	lua "github.com/yuin/gopher-lua"
)

type Interpreter struct {
	debug  bool
	script string
}

// LuaInterpreter is the lua interpreter, instantiated in main()
var LuaInterpreter *lua.LState

// holds a user provided lua function
type LuaFunction struct {
	name    string
	help    string
	numargs int
}

// LuaFuncs must be global since init() is being called from lua which
// doesn't have access to the interpreter instance
var LuaFuncs map[string]LuaFunction

func NewInterpreter(script string, debug bool) *Interpreter {
	return &Interpreter{debug: debug, script: script}
}

// initialize the lua environment properly
func (i *Interpreter) InitLua() {
	// we only  load a subset of lua Open  modules and don't allow
	// net, system or io stuff
	for _, pair := range []struct {
		n string
		f lua.LGFunction
	}{
		{lua.LoadLibName, lua.OpenPackage},
		{lua.BaseLibName, lua.OpenBase},
		{lua.TabLibName, lua.OpenTable},
		{lua.DebugLibName, lua.OpenDebug},
		{lua.MathLibName, lua.OpenMath},
	} {
		if err := LuaInterpreter.CallByParam(lua.P{
			Fn:      LuaInterpreter.NewFunction(pair.f),
			NRet:    0,
			Protect: true,
		}, lua.LString(pair.n)); err != nil {
			panic(err)
		}
	}

	// load the lua config (which we expect to contain init() and math functions)
	if err := LuaInterpreter.DoFile(i.script); err != nil {
		panic(err)
	}

	// instantiate
	LuaFuncs = map[string]LuaFunction{}

	// that way the user can call register(...) from lua inside init()
	LuaInterpreter.SetGlobal("register", LuaInterpreter.NewFunction(register))

	// actually call init()
	if err := LuaInterpreter.CallByParam(lua.P{
		Fn:      LuaInterpreter.GetGlobal("init"),
		NRet:    0,
		Protect: true,
	}); err != nil {
		panic(err)
	}
}

func (i *Interpreter) Debug(msg string) {
	if i.debug {
		fmt.Printf("DEBUG(lua): %s\n", msg)
	}
}

func (i *Interpreter) FuncNumArgs(name string) int {
	return LuaFuncs[name].numargs
}

// Call a user provided math function registered with register().
//
// Each function has  to tell us how many args  it expects, the actual
// function call  from here  is different depending  on the  number of
// arguments. 1 uses the last item of the stack, 2 the last two and -1
// all items (which translates to batch mode)
//
// The  items  array  will  be  provided  by  calc.Eval(),  these  are
// non-popped stack  items. So  the items will  only removed  from the
// stack when the lua function execution is successful.
func (i *Interpreter) CallLuaFunc(funcname string, items []float64) (float64, error) {
	i.Debug(fmt.Sprintf("calling lua func %s() with %d args",
		funcname, LuaFuncs[funcname].numargs))

	switch LuaFuncs[funcname].numargs {
	case 0, 1:
		// 1 arg variant
		if err := LuaInterpreter.CallByParam(lua.P{
			Fn:      LuaInterpreter.GetGlobal(funcname),
			NRet:    1,
			Protect: true,
		}, lua.LNumber(items[0])); err != nil {
			return 0, fmt.Errorf("failed to exec lua func %s: %w", funcname, err)
		}
	case 2:
		// 2 arg variant
		if err := LuaInterpreter.CallByParam(lua.P{
			Fn:      LuaInterpreter.GetGlobal(funcname),
			NRet:    1,
			Protect: true,
		}, lua.LNumber(items[0]), lua.LNumber(items[1])); err != nil {
			return 0, fmt.Errorf("failed to exec lua func %s: %w", funcname, err)
		}
	case -1:
		// batch variant, use lua table as array
		table := LuaInterpreter.NewTable()

		// put the whole stack into it
		for _, item := range items {
			table.Append(lua.LNumber(item))
		}

		if err := LuaInterpreter.CallByParam(lua.P{
			Fn:      LuaInterpreter.GetGlobal(funcname),
			NRet:    1,
			Protect: true,
		}, table); err != nil {
			return 0, fmt.Errorf("failed to exec lua func %s: %w", funcname, err)
		}
	}

	// get result and cast to float64
	if res, ok := LuaInterpreter.Get(-1).(lua.LNumber); ok {
		LuaInterpreter.Pop(1)

		return float64(res), nil
	}

	return 0, errors.New("function did not return a float64")
}

// called from lua to register a math  function numargs may be 1, 2 or
// -1, it denotes the number of  items from the stack requested by the
// lua function. -1 means batch mode, that is all items
func register(lstate *lua.LState) int {
	function := lstate.ToString(1)
	numargs := lstate.ToInt(2)
	help := lstate.ToString(3)

	LuaFuncs[function] = LuaFunction{
		name:    function,
		numargs: numargs,
		help:    help,
	}

	return 1
}
