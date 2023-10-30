package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

type Calc struct {
	debug     bool
	batch     bool
	stdin     bool
	stack     *Stack
	completer readline.AutoCompleter
}

const Help string = `Available commands:
batch   enable batch mode
debug   enable debug output
dump    display the stack contents
clear   clear the whole stack
shift   remove the last element of the stack
help    show this message

Available operators:
basic operators: + - * /
`

func NewCalc() *Calc {
	c := Calc{stack: NewStack(), debug: false}

	c.completer = readline.NewPrefixCompleter(
		readline.PcItem("dump"),
		readline.PcItem("debug"),
		readline.PcItem("clear"),
		readline.PcItem("batch"),
		readline.PcItem("shift"),
		readline.PcItem("undo"),
		readline.PcItem("help"),
		readline.PcItem("+"),
		readline.PcItem("-"),
		readline.PcItem("*"),
		readline.PcItem("/"),
	)

	return &c
}

func (c *Calc) ToggleDebug() {
	c.debug = !c.debug
	c.stack.ToggleDebug()
}

func (c *Calc) ToggleBatch() {
	c.batch = !c.batch
}

func (c *Calc) ToggleStdin() {
	c.stdin = !c.stdin
}

func (c *Calc) Eval(line string) {
	line = strings.TrimSpace(line)
	space := regexp.MustCompile(`\s+`)
	simple := regexp.MustCompile(`[\+\-\*\/]`)

	if line == "" {
		return
	}

	for _, item := range space.Split(line, -1) {
		num, err := strconv.ParseFloat(item, 64)

		if err == nil {
			c.stack.Backup()
			c.stack.Push(num)
		} else {
			if simple.MatchString(line) {
				c.simple(item[0])
				continue
			}

			switch item {
			case "help":
				fmt.Println(Help)
			case "dump":
				c.stack.Dump()
			case "debug":
				c.ToggleDebug()
			case "batch":
				c.ToggleBatch()
			case "clear":
				c.stack.Backup()
				c.stack.Clear()
			case "shift":
				c.stack.Backup()
				c.stack.Shift()
			case "undo":
				c.stack.Restore()
			default:
				fmt.Println("unknown command or operator!")
			}
		}
	}
}

func (c *Calc) Result() {
	if !c.stdin {
		fmt.Print("= ")
	}

	fmt.Println(c.stack.Last())
}

func (c *Calc) simple(op byte) {
	c.stack.Backup()

	for c.stack.Len() > 1 {
		b := c.stack.Pop()
		a := c.stack.Pop()
		var x float64

		if c.debug {
			fmt.Printf("DEBUG:        evaluating: %.2f %c %.2f\n", a, op, b)
		}

		switch op {
		case '+':
			x = a + b
		case '-':
			x = a - b
		case '*':
			x = a * b
		case '/':
			if b == 0 {
				fmt.Println("error: division by null!")
				return
			}
			x = a / b
		default:
			panic("invalid operator!")
		}

		c.stack.Push(x)

		if !c.batch {
			break
		}
	}

	c.Result()
}
