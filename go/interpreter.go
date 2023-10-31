package main

import (
	"errors"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// LUA interpreter, instanciated in main()
var L *lua.LState

var LuaFuncs map[string]int

// FIXME: add 2nd var with help string
// called from lua to register a 1 arg math function
func RegisterFuncOneArg(L *lua.LState) int {
	function := L.ToString(1)
	LuaFuncs[function] = 1
	return 1
}

// called from lua to register a 1 arg math function
func RegisterFuncTwoArg(L *lua.LState) int {
	function := L.ToString(1)
	LuaFuncs[function] = 2
	return 1
}

func InitLua(L *lua.LState) {
	LuaFuncs = map[string]int{}
	L.SetGlobal("RegisterFuncOneArg", L.NewFunction(RegisterFuncOneArg))
	L.SetGlobal("RegisterFuncTwoArg", L.NewFunction(RegisterFuncTwoArg))

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("init"),
		NRet:    0,
		Protect: true,
	}); err != nil {
		panic(err)
	}
}

func CallLuaFunc(L *lua.LState, funcname string, a float64, b float64) (float64, error) {
	if LuaFuncs[funcname] == 1 {
		// 1 arg variant
		if err := L.CallByParam(lua.P{
			Fn:      L.GetGlobal(funcname),
			NRet:    1,
			Protect: true,
		}, lua.LNumber(a)); err != nil {
			fmt.Println(err)
			return 0, err
		}
	} else {
		// 2 arg variant
		if err := L.CallByParam(lua.P{
			Fn:      L.GetGlobal(funcname),
			NRet:    1,
			Protect: true,
		}, lua.LNumber(a), lua.LNumber(b)); err != nil {
			return 0, err
		}
	}

	// get result and cast to float64
	if res, ok := L.Get(-1).(lua.LNumber); ok {
		L.Pop(1)
		return float64(res), nil
	}

	return 0, errors.New("function did not return a float64!")
}
