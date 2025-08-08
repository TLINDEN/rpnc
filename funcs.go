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

type Result struct {
	Res float64
	Err error
}

type Numbers []float64

type Function func(Numbers) Result

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
func NewResult(n float64, e error) Result {
	return Result{Res: n, Err: e}
}

// the actual functions, called once during initialization.
func DefineFunctions() Funcalls {
	funcmap := map[string]*Funcall{
		// simple operators, they all expect 2 args
		"+": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]+arg[1], nil)
			},
		),

		"-": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]-arg[1], nil)
			},
		),

		"x": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]*arg[1], nil)
			},
		),

		"/": NewFuncall(
			func(arg Numbers) Result {
				if arg[1] == 0 {
					return NewResult(0, errors.New("division by null"))
				}

				return NewResult(arg[0]/arg[1], nil)
			},
		),

		"^": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Pow(arg[0], arg[1]), nil)
			},
		),

		"%": NewFuncall(
			func(arg Numbers) Result {
				return NewResult((arg[0]/100)*arg[1], nil)
			},
		),

		"%-": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]-((arg[0]/100)*arg[1]), nil)
			},
		),

		"%+": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]+((arg[0]/100)*arg[1]), nil)
			},
		),

		"mod": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Remainder(arg[0], arg[1]), nil)
			},
		),

		"sqrt": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Sqrt(arg[0]), nil)
			},
			1),

		"abs": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Abs(arg[0]), nil)
			},
			1),

		"acos": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Acos(arg[0]), nil)
			},
			1),

		"acosh": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Acosh(arg[0]), nil)
			},
			1),

		"asin": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Asin(arg[0]), nil)
			},
			1),

		"asinh": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Asinh(arg[0]), nil)
			},
			1),

		"atan": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Atan(arg[0]), nil)
			},
			1),

		"atan2": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Atan2(arg[0], arg[1]), nil)
			},
			2),

		"atanh": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Atanh(arg[0]), nil)
			},
			1),

		"cbrt": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Cbrt(arg[0]), nil)
			},
			1),

		"ceil": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Ceil(arg[0]), nil)
			},
			1),

		"cos": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Cos(arg[0]), nil)
			},
			1),

		"cosh": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Cosh(arg[0]), nil)
			},
			1),

		"erf": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Erf(arg[0]), nil)
			},
			1),

		"erfc": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Erfc(arg[0]), nil)
			},
			1),

		"erfcinv": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Erfcinv(arg[0]), nil)
			},
			1),

		"erfinv": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Erfinv(arg[0]), nil)
			},
			1),

		"exp": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Exp(arg[0]), nil)
			},
			1),

		"exp2": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Exp2(arg[0]), nil)
			},
			1),

		"expm1": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Expm1(arg[0]), nil)
			},
			1),

		"floor": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Floor(arg[0]), nil)
			},
			1),

		"gamma": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Gamma(arg[0]), nil)
			},
			1),

		"ilogb": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(float64(math.Ilogb(arg[0])), nil)
			},
			1),

		"j0": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.J0(arg[0]), nil)
			},
			1),

		"j1": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.J1(arg[0]), nil)
			},
			1),

		"log": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Log(arg[0]), nil)
			},
			1),

		"log10": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Log10(arg[0]), nil)
			},
			1),

		"log1p": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Log1p(arg[0]), nil)
			},
			1),

		"log2": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Log2(arg[0]), nil)
			},
			1),

		"logb": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Logb(arg[0]), nil)
			},
			1),

		"pow": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Pow(arg[0], arg[1]), nil)
			},
			2),

		"round": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Round(arg[0]), nil)
			},
			1),

		"roundtoeven": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.RoundToEven(arg[0]), nil)
			},
			1),

		"sin": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Sin(arg[0]), nil)
			},
			1),

		"sinh": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Sinh(arg[0]), nil)
			},
			1),

		"tan": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Tan(arg[0]), nil)
			},
			1),

		"tanh": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Tanh(arg[0]), nil)
			},
			1),

		"trunc": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Trunc(arg[0]), nil)
			},
			1),

		"y0": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Y0(arg[0]), nil)
			},
			1),

		"y1": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Y1(arg[0]), nil)
			},
			1),

		"copysign": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Copysign(arg[0], arg[1]), nil)
			},
			2),

		"dim": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Dim(arg[0], arg[1]), nil)
			},
			2),

		"hypot": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(math.Hypot(arg[0], arg[1]), nil)
			},
			2),

		// converters of all kinds
		"cm-to-inch": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]/2.54, nil)
			},
			1),

		"inch-to-cm": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]*2.54, nil)
			},
			1),

		"gallons-to-liters": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]*3.785, nil)
			},
			1),

		"liters-to-gallons": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]/3.785, nil)
			},
			1),

		"yards-to-meters": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]*91.44, nil)
			},
			1),

		"meters-to-yards": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]/91.44, nil)
			},
			1),

		"miles-to-kilometers": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]*1.609, nil)
			},
			1),

		"kilometers-to-miles": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]/1.609, nil)
			},
			1),

		"bytes-to-kilobytes": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]/1024, nil)
			},
			1),

		"bytes-to-megabytes": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]/1024/1024, nil)
			},
			1),

		"bytes-to-gigabytes": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]/1024/1024/1024, nil)
			},
			1),

		"bytes-to-terabytes": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(arg[0]/1024/1024/1024/1024, nil)
			},
			1),

		"or": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(float64(int(arg[0])|int(arg[1])), nil)
			},
			2),

		"and": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(float64(int(arg[0])&int(arg[1])), nil)
			},
			2),

		"xor": NewFuncall(
			func(arg Numbers) Result {
				return NewResult(float64(int(arg[0])^int(arg[1])), nil)
			},
			2),

		"<": NewFuncall(
			func(arg Numbers) Result {
				// Shift by negative number provibited, so check it.
				// Note that we check against uint64 overflow as well here
				if arg[1] < 0 || uint64(arg[1]) > math.MaxInt64 {
					return NewResult(0, errors.New("negative shift amount"))
				}

				return NewResult(float64(int(arg[0])<<int(arg[1])), nil)
			},
			2),

		">": NewFuncall(
			func(arg Numbers) Result {
				if arg[1] < 0 || uint64(arg[1]) > math.MaxInt64 {
					return NewResult(0, errors.New("negative shift amount"))
				}

				return NewResult(float64(int(arg[0])>>int(arg[1])), nil)
			},
			2),
	}

	// aliases
	funcmap["*"] = funcmap["x"]
	funcmap["remainder"] = funcmap["mod"]

	return funcmap
}

func DefineBatchFunctions() Funcalls {
	funcmap := map[string]*Funcall{
		"median": NewFuncall(
			func(args Numbers) Result {
				middle := len(args) / 2

				return NewResult(args[middle], nil)
			},
			-1),

		"mean": NewFuncall(
			func(args Numbers) Result {
				var sum float64
				for _, item := range args {
					sum += item
				}

				return NewResult(sum/float64(len(args)), nil)
			},
			-1),

		"min": NewFuncall(
			func(args Numbers) Result {
				var min float64
				min, args = args[0], args[1:]
				for _, item := range args {
					if item < min {
						min = item
					}
				}

				return NewResult(min, nil)
			},
			-1),

		"max": NewFuncall(
			func(args Numbers) Result {
				var max float64
				max, args = args[0], args[1:]
				for _, item := range args {
					if item > max {
						max = item
					}
				}

				return NewResult(max, nil)
			},
			-1),

		"sum": NewFuncall(
			func(args Numbers) Result {
				var sum float64
				for _, item := range args {
					sum += item
				}

				return NewResult(sum, nil)
			},
			-1),
	}

	// aliases
	funcmap["+"] = funcmap["sum"]
	funcmap["avg"] = funcmap["mean"]

	return funcmap
}
