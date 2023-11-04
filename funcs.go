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
	"fmt"
	"math"
	"strings"
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
	}

	// aliases
	f["*"] = f["x"]

	return f
}

func list2str(list Numbers) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(list)), " "), "[]")
}

// we need to add a history entry for each operation
func (c *Calc) SetHistory(op string, args Numbers) {
	c.History("%s %s", list2str(args))
}

// Execute a math function, check if it is defined just in case
//
// FIXME:  add a  loop over  DoFuncall() for  non-batch-only functions
// like + or *
//
// FIXME: use R{} as well? or even everywhere, while we're at it?
func (c *Calc) DoFuncall(funcname string) (float64, error) {
	if function, ok := c.Functions[funcname]; ok {
		args := Numbers{}
		batch := false

		if function.Expectargs == -1 {
			// batch mode, but always < stack len, so check first
			args = c.stack.All()
			batch = true
		} else {
			//  this is way better behavior than just using 0 in place of
			// non-existing stack items
			if c.stack.Len() < function.Expectargs {
				return -1, errors.New("stack doesn't provide enough arguments")
			}

			args = c.stack.Last(function.Expectargs)
		}

		// the  actual lambda call, so  to say. We provide  a slice of
		// the requested size, fetched  from the stack (but not popped
		// yet!)
		R := function.Func(args)

		if R.Err != nil {
			// leave the stack untouched in case of any error
			return R.Res, R.Err
		}

		if batch {
			// get rid of stack
			c.stack.Clear()
		} else {
			// remove operands
			c.stack.Shift(function.Expectargs)
		}

		// save result
		c.stack.Push(R.Res)

		// thanks a lot
		c.SetHistory(funcname, args)
		return R.Res, nil
	}

	// should not happen, if it does: programmer fault!
	return -1, errors.New("heck, no such function")
}
