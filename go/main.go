package main

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
	flag "github.com/spf13/pflag"
	lua "github.com/yuin/gopher-lua"
)

const VERSION string = "0.0.1"

const Usage string = `This is rpn, a reverse polish notation calculator cli.

Usage: rpn [-bdvh] [<operator>]

Options:
  -b, --batchmode   enable batch mode
  -d, --debug       enable debug mode
  -v, --version     show version
  -h, --help        show help

When <operator>  is given, batch  mode ist automatically  enabled. Use
this only when working with stdin. E.g.: echo "2 3 4 5" | rpn +

Copyright (c) 2023 T.v.Dein`

func main() {
	calc := NewCalc()

	showversion := false
	showhelp := false
	enabledebug := false
	configfile := ""

	flag.BoolVarP(&calc.batch, "batchmode", "b", false, "batch mode")
	flag.BoolVarP(&enabledebug, "debug", "d", false, "debug mode")
	flag.BoolVarP(&showversion, "version", "v", false, "show version")
	flag.BoolVarP(&showhelp, "help", "h", false, "show usage")
	flag.StringVarP(&configfile, "config", "c", os.Getenv("HOME")+"/.rpn.lua", "config file (lua format)")

	flag.Parse()

	if showversion {
		fmt.Printf("This is rpn version %s\n", VERSION)
		return
	}

	if showhelp {
		fmt.Println(Usage)
		return
	}

	if enabledebug {
		calc.ToggleDebug()
	}

	if _, err := os.Stat(configfile); err == nil {
		// FIXME: put into interpreter.go, probably with its own obj
		// then just Interpreter.Init(configfile) should suffice
		L = lua.NewState(lua.Options{SkipOpenLibs: true})
		defer L.Close()

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
			if err := L.CallByParam(lua.P{
				Fn:      L.NewFunction(pair.f),
				NRet:    0,
				Protect: true,
			}, lua.LString(pair.n)); err != nil {
				panic(err)
			}
		}

		if err := L.DoFile(configfile); err != nil {
			panic(err)
		}

		InitLua(L)
		calc.L = L
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:            "\033[31m»\033[0m ",
		HistoryFile:       os.Getenv("HOME") + "/.rpn-history",
		HistoryLimit:      500,
		AutoComplete:      calc.completer,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	})

	if err != nil {
		panic(err)
	}
	defer rl.Close()
	rl.CaptureExitSignal()

	if inputIsStdin() {
		calc.ToggleStdin()
	}

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		calc.Eval(line)
	}

	if len(flag.Args()) > 0 {
		calc.Eval(flag.Args()[0])
	}
}

func inputIsStdin() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
