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
	"math"
)

type R struct {
	Res float64
	Err error
}

type Numbers []float64

type Function func(Numbers) R

// every function we  are able to call must be  of type Funcall, which
// needs to  specify how  many numbers  it expects  and the  actual go
// function to be executed.
//
// The function  has to take  a float slice  as argument and  return a
// float and  an error object. The  float slice is guaranteed  to have
// the expected number of arguments.
//
// However, Lua functions are handled differently, see interpreter.go.
type Funcall struct {
	Expectargs int // -1 means batch only mode, you'll get the whole stack as arg
	Func       Function
}

// will hold all hard coded functions and operators
type Funcalls map[string]*Funcall

// convenience function,  create a  new Funcall object,  if expectargs
// was not specified, 2 is assumed.
func NewFuncall(function Function, expectargs ...int) *Funcall {
	expect := 2

	if len(expectargs) > 0 {
		expect = expectargs[0]
	}

	return &Funcall{
		Expectargs: expect,
		Func:       function,
	}
}

// Convenience function, create new result
func NewR(n float64, e error) R {
	return R{Res: n, Err: e}
}

// the actual functions, called once during initialization.
func DefineFunctions() Funcalls {
	f := map[string]*Funcall{
		// simple operators, they all expect 2 args
		"+": NewFuncall(
			func(arg Numbers) R {
				return NewR(arg[0]+arg[1], nil)
			},
		),

		"-": NewFuncall(
			func(arg Numbers) R {
				return NewR(arg[0]-arg[1], nil)
			},
		),

		"x": NewFuncall(
			func(arg Numbers) R {
				return NewR(arg[0]*arg[1], nil)
			},
		),

		"/": NewFuncall(
			func(arg Numbers) R {
				if arg[1] == 0 {
					return NewR(0, errors.New("division by null"))
				}

				return NewR(arg[0]/arg[1], nil)
			},
		),

		"^": NewFuncall(
			func(arg Numbers) R {
				return NewR(math.Pow(arg[0], arg[1]), nil)
			},
		),

		"%": NewFuncall(
			func(arg Numbers) R {
				return NewR((arg[0]/100)*arg[1], nil)
			},
		),

		"%-": NewFuncall(
			func(arg Numbers) R {
				return NewR(arg[0]-((arg[0]/100)*arg[1]), nil)
			},
		),

		"%+": NewFuncall(
			func(arg Numbers) R {
				return NewR(arg[0]+((arg[0]/100)*arg[1]), nil)
			},
		),

		"mod": NewFuncall(
			func(arg Numbers) R {
				return NewR(math.Remainder(arg[0], arg[1]), nil)
			},
		),

		"sqrt": NewFuncall(
			func(arg Numbers) R {
				return NewR(math.Sqrt(arg[0]), nil)
			},
			1),
	}

	// aliases
	f["*"] = f["x"]
	f["mod"] = f["remainder"]

	return f
}

func DefineBatchFunctions() Funcalls {
	f := map[string]*Funcall{
		"median": NewFuncall(
			func(args Numbers) R {
				middle := len(args) / 2
				return NewR(args[middle], nil)
			},
			-1),

		"mean": NewFuncall(
			func(args Numbers) R {
				var sum float64
				for _, item := range args {
					sum += item
				}
				return NewR(sum/float64(len(args)), nil)
			},
			-1),

		"min": NewFuncall(
			func(args Numbers) R {
				var min float64
				min, args = args[0], args[1:]
				for _, item := range args {
					if item < min {
						min = item
					}
				}
				return NewR(min, nil)
			},
			-1),

		"max": NewFuncall(
			func(args Numbers) R {
				var max float64
				max, args = args[0], args[1:]
				for _, item := range args {
					if item > max {
						max = item
					}
				}
				return NewR(max, nil)
			},
			-1),

		"sum": NewFuncall(
			func(args Numbers) R {
				var sum float64
				for _, item := range args {
					sum += item
				}
				return NewR(sum, nil)
			},
			-1),
	}

	// aliases
	f["+"] = f["sum"]
	f["avg"] = f["mean"]

	return f
}
