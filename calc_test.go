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
)

func TestCalc(t *testing.T) {
	calc := NewCalc()

	var tests = []struct {
		name  string
		cmd   string
		exp   float64
		batch bool
	}{
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
