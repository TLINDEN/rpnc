package main

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
	flag "github.com/spf13/pflag"
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

Copyright (c) 2023 T.v.Dein
`

func main() {
	calc := NewCalc()

	showversion := false
	showhelp := false
	enabledebug := false

	flag.BoolVarP(&calc.batch, "batchmode", "b", false, "batch mode")
	flag.BoolVarP(&enabledebug, "debug", "d", false, "debug mode")
	flag.BoolVarP(&showversion, "version", "v", false, "show version")
	flag.BoolVarP(&showhelp, "help", "h", false, "show usage")
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
