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
	"fmt"
	"strconv"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestCommentsAndWhitespace(t *testing.T) {
	calc := NewCalc()

	var tests = []struct {
		name string
		cmd  []string
		exp  float64 // last element of the stack
	}{
		{
			name: "whitespace prefix",
			cmd:  []string{"  5"},
			exp:  5.0,
		},
		{
			name: "whitespace postfix",
			cmd:  []string{"5  "},
			exp:  5.0,
		},
		{
			name: "whitespace both",
			cmd:  []string{"  5   "},
			exp:  5.0,
		},
		{
			name: "comment line w/ spaces",
			cmd:  []string{"5", "   #   19"},
			exp:  5.0,
		},
		{
			name: "comment line w/o spaces",
			cmd:  []string{"5", `#19`},
			exp:  5.0,
		},
		{
			name: "inline comment w/ spaces",
			cmd:  []string{"5   #   19"},
			exp:  5.0,
		},
		{
			name: "inline comment w/o spaces",
			cmd:  []string{"5#19"},
			exp:  5.0,
		},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s .(expect %.2f)",
			test.name, test.exp)

		t.Run(testname, func(t *testing.T) {
			for _, line := range test.cmd {
				if err := calc.Eval(line); err != nil {
					t.Error(err.Error())
				}
			}
			got := calc.stack.Last()

			if len(got) > 0 {
				if got[0] != test.exp {
					t.Errorf("parsing failed:\n+++  got: %f\n--- want: %f",
						got, test.exp)
				}
			}

			if calc.stack.Len() != 1 {
				t.Errorf("invalid stack size:\n+++  got: %d\n--- want: 1",
					calc.stack.Len())
			}
		})

		calc.stack.Clear()
	}
}

func TestCalc(t *testing.T) {
	calc := NewCalc()

	var tests = []struct {
		name  string
		cmd   string
		exp   float64
		batch bool
	}{
		// ops
		{
			name: "plus",
			cmd:  `15 15 +`,
			exp:  30,
		},
		{
			name: "power",
			cmd:  `4 2 ^`,
			exp:  16,
		},
		{
			name: "minus",
			cmd:  `100 50 -`,
			exp:  50,
		},
		{
			name: "multi",
			cmd:  `4 4 x`,
			exp:  16,
		},
		{
			name: "divide",
			cmd:  `10 2 /`,
			exp:  5,
		},
		{
			name: "percent",
			cmd:  `400 20 %`,
			exp:  80,
		},
		{
			name: "percent-minus",
			cmd:  `400 20 %-`,
			exp:  320,
		},
		{
			name: "percent-plus",
			cmd:  `400 20 %+`,
			exp:  480,
		},

		// math tests
		{
			name: "mod",
			cmd:  `9 2 mod`,
			exp:  1,
		},
		{
			name: "sqrt",
			cmd:  `16 sqrt`,
			exp:  4,
		},
		{
			name: "ceil",
			cmd:  `15.5 ceil`,
			exp:  16,
		},
		{
			name: "dim",
			cmd:  `6 4 dim`,
			exp:  2,
		},

		// constants tests
		{
			name: "pitimes2",
			cmd:  `Pi 2 *`,
			exp:  6.283185307179586,
		},
		{
			name: "pi+sqrt2",
			cmd:  `Pi Sqrt2 +`,
			exp:  4.555806215962888,
		},

		// batch tests
		{
			name:  "batch-sum",
			cmd:   `2 2 2 2 sum`,
			exp:   8,
			batch: true,
		},
		{
			name:  "batch-median",
			cmd:   `1 2 3 4 5 median`,
			exp:   3,
			batch: true,
		},
		{
			name:  "batch-mean",
			cmd:   `2 2 8 2 2 mean`,
			exp:   3.2,
			batch: true,
		},
		{
			name:  "batch-min",
			cmd:   `1 2 3 4 5 min`,
			exp:   1,
			batch: true,
		},
		{
			name:  "batch-max",
			cmd:   `1 2 3 4 5 max`,
			exp:   5,
			batch: true,
		},

		// stack tests
		{
			name: "use-vars",
			cmd:  `10 >TEN clear 5 <TEN *`,
			exp:  50,
		},
		{
			name: "reverse",
			cmd:  `100 500 reverse -`,
			exp:  400,
		},
		{
			name: "swap",
			cmd:  `2 16 swap /`,
			exp:  8,
		},
		{
			name:  "clear batch",
			cmd:   "1 1 1 1 1 clear 1 1 sum",
			exp:   2,
			batch: true,
		},
		{
			name: "undo",
			cmd:  `4 4 + undo *`,
			exp:  16,
		},

		// bit tests
		{
			name: "bit and",
			cmd:  `1 3 and`,
			exp:  1,
		},
		{
			name: "bit or",
			cmd:  `1 3 or`,
			exp:  3,
		},
		{
			name: "bit xor",
			cmd:  `1 3 xor`,
			exp:  2,
		},

		// converters
		{
			name: "inch-to-cm",
			cmd:  `111 inch-to-cm`,
			exp:  281.94,
		},
		{
			name: "gallons-to-liters",
			cmd:  `111 gallons-to-liters`,
			exp:  420.135,
		},
		{
			name: "meters-to-yards",
			cmd:  `111 meters-to-yards`,
			exp:  1.2139107611548556,
		},
		{
			name: "miles-to-kilometers",
			cmd:  `111 miles-to-kilometers`,
			exp:  178.599,
		},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("cmd-%s-expect-%.2f",
			test.name, test.exp)

		t.Run(testname, func(t *testing.T) {
			calc.batch = test.batch
			if err := calc.Eval(test.cmd); err != nil {
				t.Error(err.Error())
			}
			got := calc.Result()
			calc.stack.Clear()
			if got != test.exp {
				t.Errorf("calc failed:\n+++  got: %f\n--- want: %f",
					got, test.exp)
			}
		})
	}
}

func TestCalcLua(t *testing.T) {
	var tests = []struct {
		function string
		stack    []float64
		exp      float64
	}{
		{
			function: "lower",
			stack:    []float64{5, 6},
			exp:      5.0,
		},
		{
			function: "parallelresistance",
			stack:    []float64{100, 200, 300},
			exp:      54.54545454545455,
		},
	}

	calc := NewCalc()

	LuaInterpreter = lua.NewState(lua.Options{SkipOpenLibs: true})
	defer LuaInterpreter.Close()

	luarunner := NewInterpreter("example.lua", false)
	luarunner.InitLua()
	calc.SetInt(luarunner)

	for _, test := range tests {
		testname := fmt.Sprintf("lua-%s", test.function)

		t.Run(testname, func(t *testing.T) {
			calc.stack.Clear()
			for _, item := range test.stack {
				calc.stack.Push(item)
			}

			calc.EvalLuaFunction(test.function)

			got := calc.stack.Last()

			if calc.stack.Len() != 1 {
				t.Errorf("invalid stack size:\n+++  got: %d\n--- want: 1",
					calc.stack.Len())
			}

			if got[0] != test.exp {
				t.Errorf("lua function %s failed:\n+++  got: %f\n--- want: %f",
					test.function, got, test.exp)
			}
		})
	}
}

func FuzzEval(f *testing.F) {
	legal := []string{
		"dump",
		"showstack",
		"help",
		"Pi 31 *",
		"SqrtE Pi /",
		"55.5 yards-to-meters",
		"2 4 +",
		"7 8 batch sum",
		"7 8 %-",
		"7 8 clear",
		"7 8 /",
		"b",
		"#444",
		"<X",
		"?",
		"help",
	}

	for _, item := range legal {
		f.Add(item)
	}

	calc := NewCalc()

	var hexnum, hour, min int

	f.Fuzz(func(t *testing.T, line string) {
		t.Logf("Stack:\n%v\nLine: <%s>\n", calc.stack.All(), line)
		switch line {
		case "help", "?":
			return
		}
		if err := calc.EvalItem(line); err == nil {
			t.Logf("given: <%s>", line)
			// not corpus and empty?
			if !contains(legal, line) && len(line) > 0 {
				item := strings.TrimSpace(calc.Comment.ReplaceAllString(line, ""))
				_, hexerr := fmt.Sscanf(item, "0x%x", &hexnum)
				_, timeerr := fmt.Sscanf(item, "%d:%d", &hour, &min)
				// no comment?
				if len(item) > 0 {
					// no known command or function?
					if _, err := strconv.ParseFloat(item, 64); err != nil {
						if !contains(calc.Constants, item) &&
							!exists(calc.Funcalls, item) &&
							!exists(calc.BatchFuncalls, item) &&
							!contains(calc.LuaFunctions, item) &&
							!exists(calc.Commands, item) &&
							!exists(calc.ShowCommands, item) &&
							!exists(calc.SettingsCommands, item) &&
							!exists(calc.StackCommands, item) &&
							!calc.Register.MatchString(item) &&
							item != "?" && item != "help" &&
							hexerr != nil &&
							timeerr != nil {
							t.Errorf("Fuzzy input accepted: <%s>", line)
						}
					}
				}
			}
		}
	})
}
