package main

import (
	"fmt"
	"math"
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
	history   []string
	completer readline.AutoCompleter
}

const Help string = `Available commands:
batch   enable batch mode
debug   enable debug output
dump    display the stack contents
clear   clear the whole stack
shift   remove the last element of the stack
history display calculation history
help    show this message

Available operators:
basic operators: + - * /

Math operators:
^       power`

func NewCalc() *Calc {
	c := Calc{stack: NewStack(), debug: false}

	c.completer = readline.NewPrefixCompleter(
		// commands
		readline.PcItem("dump"),
		readline.PcItem("reverse"),
		readline.PcItem("debug"),
		readline.PcItem("clear"),
		readline.PcItem("batch"),
		readline.PcItem("shift"),
		readline.PcItem("undo"),
		readline.PcItem("help"),
		readline.PcItem("history"),

		// ops
		readline.PcItem("+"),
		readline.PcItem("-"),
		readline.PcItem("*"),
		readline.PcItem("/"),
		readline.PcItem("^"),
		readline.PcItem("%"),
		readline.PcItem("%-"),
		readline.PcItem("%+"),

		// constants
		readline.PcItem("Pi"),
		readline.PcItem("Phi"),
		readline.PcItem("Sqrt2"),
		readline.PcItem("SqrtE"),
		readline.PcItem("SqrtPi"),
		readline.PcItem("SqrtPhi"),
		readline.PcItem("Ln2"),
		readline.PcItem("Log2E"),
		readline.PcItem("Ln10"),
		readline.PcItem("Log10E"),

		// math functions
		readline.PcItem("sqrt"),
		readline.PcItem("remainder"),
		readline.PcItem("avg"),
		readline.PcItem("median"),
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
	simple := regexp.MustCompile(`^[\+\-\*\/]$`)
	constants := []string{"E", "Pi", "Phi", "Sqrt2", "SqrtE", "SqrtPi",
		"SqrtPhi", "Ln2", "Log2E", "Ln10", "Log10E"}
	functions := []string{"sqrt", "remainder", "%", "%-", "%+"}
	batch := []string{"median", "avg"}

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

			if contains(constants, item) {
				c.stack.Backup()
				c.stack.Push(const2num(item))
				continue
			}

			if contains(functions, item) {
				c.mathfunc(item)
				continue
			}

			if contains(batch, item) {
				c.batchfunc(item)
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
			case "reverse":
				c.stack.Backup()
				c.stack.Reverse()
			case "undo":
				c.stack.Restore()
			case "history":
				for _, entry := range c.history {
					fmt.Println(entry)
				}
			case "^":
				c.exp()
			default:
				fmt.Println("unknown command or operator!")
			}
		}
	}
}

func (c *Calc) History(format string, args ...any) {
	c.history = append(c.history, fmt.Sprintf(format, args...))
}

func (c *Calc) Result() float64 {
	if !c.stdin {
		fmt.Print("= ")
	}

	fmt.Println(c.stack.Last())

	return c.stack.Last()
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

		c.History("%f %c %f = %f", a, op, b, x)

		if !c.batch {
			break
		}
	}

	_ = c.Result()
}

func (c *Calc) exp() {
	c.stack.Backup()

	for c.stack.Len() > 1 {
		b := c.stack.Pop()
		a := c.stack.Pop()
		x := math.Pow(a, b)

		c.stack.Push(x)

		c.History("%f ^ %f = %f", a, b, x)

		if !c.batch {
			break
		}
	}

	_ = c.Result()
}

func (c *Calc) mathfunc(funcname string) {
	// FIXME: split into 2 funcs, one  working with 1 the other with 2
	// args, saving Pop calls
	c.stack.Backup()

	for c.stack.Len() > 0 {
		var x float64

		switch funcname {
		case "sqrt":
			a := c.stack.Pop()
			x = math.Sqrt(a)
			c.History("sqrt(%f) = %f", a, x)

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

	_ = c.Result()
}

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

	case "avg":
		var sum float64

		for c.stack.Len() > 0 {
			sum += c.stack.Pop()
		}

		x = sum / float64(count)
		c.History("avg(all)")
	}

	c.stack.Push(x)
	_ = c.Result()
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func const2num(name string) float64 {
	switch name {
	case "Pi":
		return math.Pi
	case "Phi":
		return math.Phi
	case "Sqrt2":
		return math.Sqrt2
	case "SqrtE":
		return math.SqrtE
	case "SqrtPi":
		return math.SqrtPi
	case "SqrtPhi":
		return math.SqrtPhi
	case "Ln2":
		return math.Ln2
	case "Log2E":
		return math.Log2E
	case "Ln10":
		return math.Ln10
	case "Log10E":
		return math.Log10E
	default:
		return 0
	}
}
