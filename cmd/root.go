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

package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chzyer/readline"
	flag "github.com/spf13/pflag"
	lua "github.com/yuin/gopher-lua"
)

const VERSION string = "2.1.7"

const Usage string = `This is rpn, a reverse polish notation calculator cli.

Usage: rpn [-bdvh] [<operator>]

Options:
  -b, --batchmode       enable batch mode
  -d, --debug           enable debug mode
  -s, --stack           show last 5 items of the stack (off by default)
  -i  --intermediate    print intermediate results
  -m, --manual          show manual
  -c, --config <file>   load <file> containing LUA code
  -p, --precision <int> floating point number precision (default 2)
  -v, --version         show version
  -h, --help            show help

When <operator>  is given, batch  mode ist automatically  enabled. Use
this only when working with stdin. E.g.: echo "2 3 4 5" | rpn +

Copyright (c) 2023-2025 T.v.Dein`

func Main() int {
	calc := NewCalc()

	showversion := false
	showhelp := false
	showmanual := false
	enabledebug := false
	configfile := ""

	flag.BoolVarP(&calc.batch, "batchmode", "b", false, "batch mode")
	flag.BoolVarP(&calc.showstack, "show-stack", "s", false, "show stack")
	flag.BoolVarP(&calc.intermediate, "showin-termediate", "i", false,
		"show intermediate results")
	flag.BoolVarP(&enabledebug, "debug", "d", false, "debug mode")
	flag.BoolVarP(&showversion, "version", "v", false, "show version")
	flag.BoolVarP(&showhelp, "help", "h", false, "show usage")
	flag.BoolVarP(&showmanual, "manual", "m", false, "show manual")
	flag.StringVarP(&configfile, "config", "c",
		os.Getenv("HOME")+"/.rpn.lua", "config file (lua format)")
	flag.IntVarP(&calc.precision, "precision", "p", Precision, "floating point precision")

	flag.Parse()

	if showversion {
		fmt.Printf("This is rpn version %s\n", VERSION)

		return 0
	}

	if showhelp {
		fmt.Println(Usage)

		return 0
	}

	if enabledebug {
		calc.ToggleDebug()
	}

	if showmanual {
		man()

		return 0
	}

	// the lua state object is global, instantiate it early
	LuaInterpreter = lua.NewState(lua.Options{SkipOpenLibs: true})
	defer LuaInterpreter.Close()

	// our config file is interpreted  as lua code, only functions can
	// be defined, init() will be called by InitLua().
	if _, err := os.Stat(configfile); err == nil {
		luarunner := NewInterpreter(configfile, enabledebug)
		luarunner.InitLua()
		calc.SetInt(luarunner)

		if calc.debug {
			fmt.Println("loaded config")
		}
	} else if calc.debug {
		fmt.Println(err)
	}

	if len(flag.Args()) > 1 {
		// commandline calc operation, no readline etc needed
		// called like rpn 2 2 +
		calc.stdin = true
		if err := calc.Eval(strings.Join(flag.Args(), " ")); err != nil {
			fmt.Println(err)

			return 1
		}

		return 0
	}

	// interactive mode, need readline
	reader, err := readline.NewEx(&readline.Config{
		Prompt:            calc.Prompt(),
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
	defer func() {
		if err := reader.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	reader.CaptureExitSignal()

	if inputIsStdin() {
		// commands are  coming on stdin, however we  will still enter
		// the same loop since readline just reads fine from stdin
		calc.ToggleStdin()
	}

	for {
		// primary program repl
		line, err := reader.Readline()
		if err != nil {
			break
		}

		err = calc.Eval(line)
		if err != nil {
			fmt.Println(err)
		}

		reader.SetPrompt(calc.Prompt())
	}

	if len(flag.Args()) > 0 {
		// called like this:
		// echo 1 2 3 4 | rpn +
		// batch mode enabled automatically
		calc.batch = true
		if err = calc.Eval(flag.Args()[0]); err != nil {
			fmt.Println(err)

			return 1
		}
	}

	return 0
}

func inputIsStdin() bool {
	stat, _ := os.Stdin.Stat()

	return (stat.Mode() & os.ModeCharDevice) == 0
}

func man() {
	Pager("rpn manual page", manpage)
}
