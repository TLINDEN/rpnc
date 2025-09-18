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

package cmd

import (
	"container/list"
	"fmt"
	"sync"
)

// The stack uses a linked  list provided by container/list as storage
// and works after the LIFO principle (last in first out). Most of the
// work is  being done  in the  linked list,  but we  add a  couple of
// cenvenient functions,  so that the  user doesn't have to  cope with
// list directly.

type Stack struct {
	linklist  list.List
	backup    list.List
	debug     bool
	rev       int
	backuprev int
	mutex     sync.Mutex
}

// FIXME: maybe use a separate stack  object for backup so that it has
// its own revision etc
func NewStack() *Stack {
	return &Stack{
		linklist:  list.List{},
		backup:    list.List{},
		rev:       0,
		backuprev: 0,
	}
}

func (s *Stack) Debug(msg string) {
	if s.debug {
		fmt.Printf("DEBUG(%03d): %s\n", s.rev, msg)
	}
}

func (s *Stack) ToggleDebug() {
	s.debug = !s.debug
}

func (s *Stack) Bump() {
	s.rev++
}

// append an item to the stack
func (s *Stack) Push(item float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Debug(fmt.Sprintf("     push to stack: %.2f", item))

	s.Bump()
	s.linklist.PushBack(item)
}

// remove and return an item from the stack
func (s *Stack) Pop() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.linklist.Len() == 0 {
		return 0
	}

	tail := s.linklist.Back()
	val := tail.Value
	s.linklist.Remove(tail)

	s.Debug(fmt.Sprintf(" remove from stack: %.2f", val))

	s.Bump()

	return val.(float64)
}

// just remove the last item, do not return it
func (s *Stack) Shift(num ...int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	count := 1

	if len(num) > 0 {
		count = num[0]
	}

	if s.linklist.Len() == 0 {
		return
	}

	for i := 0; i < count; i++ {
		tail := s.linklist.Back()
		s.linklist.Remove(tail)
		s.Debug(fmt.Sprintf("remove from stack: %.2f", tail.Value))
	}
}

func (s *Stack) Swap() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.linklist.Len() < 2 {
		return
	}

	prevA := s.linklist.Back()
	s.linklist.Remove(prevA)

	prevB := s.linklist.Back()
	s.linklist.Remove(prevB)

	s.Debug(fmt.Sprintf("swapping %.2f with %.2f", prevB.Value, prevA.Value))

	s.linklist.PushBack(prevA.Value)
	s.linklist.PushBack(prevB.Value)
}

// Return the last num items from the stack w/o modifying it.
func (s *Stack) Last(num ...int) []float64 {
	items := []float64{}
	stacklen := s.Len()
	count := 1

	if len(num) > 0 {
		count = num[0]
	}

	for e := s.linklist.Front(); e != nil; e = e.Next() {
		if stacklen <= count {
			items = append(items, e.Value.(float64))
		}
		stacklen--
	}

	return items
}

// Return all elements of the stack without modifying it.
func (s *Stack) All() []float64 {
	items := []float64{}

	for e := s.linklist.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value.(float64))
	}

	return items
}

// dump the stack to stdout, including backup if debug is enabled
func (s *Stack) Dump() {
	fmt.Printf("Stack revision %d (%p):\n", s.rev, &s.linklist)

	for e := s.linklist.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

	if s.debug {
		fmt.Printf("Backup stack revision %d (%p):\n", s.backuprev, &s.backup)

		for e := s.backup.Front(); e != nil; e = e.Next() {
			fmt.Println(e.Value)
		}
	}
}

func (s *Stack) Clear() {
	s.Debug("clearing stack")

	s.linklist = list.List{}
}

func (s *Stack) Len() int {
	return s.linklist.Len()
}

func (s *Stack) Backup() {
	// we need clean the list and restore it from scratch each time we
	// make a backup, because the elements in list.List{} are pointers
	// and lead to unexpected  results. The methid here works reliably
	// at least.
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Debug(fmt.Sprintf("backing up %d items from rev %d",
		s.linklist.Len(), s.rev))

	s.backup = list.List{}
	for e := s.linklist.Front(); e != nil; e = e.Next() {
		s.backup.PushBack(e.Value.(float64))
	}

	s.backuprev = s.rev
}

func (s *Stack) Restore() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.rev == 0 {
		fmt.Println("error: stack is empty.")

		return
	}

	s.Debug(fmt.Sprintf("restoring stack to revision %d", s.backuprev))

	s.rev = s.backuprev

	s.linklist = list.List{}
	for e := s.backup.Front(); e != nil; e = e.Next() {
		s.linklist.PushBack(e.Value.(float64))
	}
}

func (s *Stack) Reverse() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	items := []float64{}

	for e := s.linklist.Front(); e != nil; e = e.Next() {
		tail := s.linklist.Back()
		items = append(items, tail.Value.(float64))
		s.linklist.Remove(tail)
	}

	for i := len(items) - 1; i >= 0; i-- {
		s.linklist.PushFront(items[i])
	}
}
