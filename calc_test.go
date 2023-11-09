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

	for _, tt := range tests {
		testname := fmt.Sprintf("%s .(expect %.2f)",
			tt.name, tt.exp)

		t.Run(testname, func(t *testing.T) {
			for _, line := range tt.cmd {
				calc.Eval(line)
			}
			got := calc.stack.Last()

			if len(got) > 0 {
				if got[0] != tt.exp {
					t.Errorf("parsing failed:\n+++  got: %f\n--- want: %f",
						got, tt.exp)
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
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("cmd-%s-expect-%.2f",
			tt.name, tt.exp)

		t.Run(testname, func(t *testing.T) {
			calc.batch = tt.batch
			calc.Eval(tt.cmd)
			got := calc.Result()
			calc.stack.Clear()
			if got != tt.exp {
				t.Errorf("calc failed:\n+++  got: %f\n--- want: %f",
					got, tt.exp)
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
	L = lua.NewState(lua.Options{SkipOpenLibs: true})
	defer L.Close()

	luarunner := NewInterpreter("example.lua", false)
	luarunner.InitLua()
	calc.SetInt(luarunner)

	for _, tt := range tests {
		testname := fmt.Sprintf("lua-%s", tt.function)

		t.Run(testname, func(t *testing.T) {
			calc.stack.Clear()
			for _, item := range tt.stack {
				calc.stack.Push(item)
			}

			calc.EvalLuaFunction(tt.function)

			got := calc.stack.Last()

			if calc.stack.Len() != 1 {
				t.Errorf("invalid stack size:\n+++  got: %d\n--- want: 1",
					calc.stack.Len())
			}

			if got[0] != tt.exp {
				t.Errorf("lua function %s failed:\n+++  got: %f\n--- want: %f",
					tt.function, got, tt.exp)
			}
		})
	}
}
