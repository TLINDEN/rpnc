package main

import (
	"container/list"
	"fmt"
	"sync"
)

type Stack struct {
	dll       list.List
	backup    list.List
	debug     bool
	rev       int
	backuprev int
	mutex     sync.Mutex
}

func NewStack() *Stack {
	return &Stack{dll: list.List{}, backup: list.List{}, rev: 0, backuprev: 0}
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

func (s *Stack) Push(x float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Debug(fmt.Sprintf("     push to stack: %.2f", x))

	s.Bump()
	s.dll.PushBack(x)
}

func (s *Stack) Pop() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.dll.Len() == 0 {
		return 0
	}

	tail := s.dll.Back()
	val := tail.Value
	s.dll.Remove(tail)

	s.Debug(fmt.Sprintf("remove from stack: %.2f", val))

	s.Bump()
	return val.(float64)
}

func (s *Stack) Shift() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.dll.Len() == 0 {
		return
	}

	tail := s.dll.Back()
	s.dll.Remove(tail)

	s.Debug(fmt.Sprintf("remove from stack: %.2f", tail.Value))
}

func (s *Stack) Last() float64 {
	if s.dll.Back() == nil {
		return 0
	}

	return s.dll.Back().Value.(float64)
}

func (s *Stack) Dump() {
	fmt.Printf("Stack revision %d (%p):\n", s.rev, &s.dll)
	for e := s.dll.Front(); e != nil; e = e.Next() {
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
	s.Debug("DEBUG: clearing stack")

	s.dll = list.List{}
}

func (s *Stack) Len() int {
	return s.dll.Len()
}

func (s *Stack) Backup() {
	// we need clean the list and restore it from scratch each time we
	// make a backup, because the elements in list.List{} are pointers
	// and lead to unexpected  results. The methid here works reliably
	// at least.
	s.backup = list.List{}
	for e := s.dll.Front(); e != nil; e = e.Next() {
		s.backup.PushBack(e.Value)
	}
	s.backuprev = s.rev
}

func (s *Stack) Restore() {
	if s.rev == 0 {
		fmt.Println("error: stack is empty.")
		return
	}

	s.Debug(fmt.Sprintf("restoring stack to revision %d", s.backuprev))

	s.rev = s.backuprev
	s.dll = s.backup
}

func (s *Stack) Reverse() {
	newstack := list.List{}

	for e := s.dll.Front(); e != nil; e = e.Next() {
		newstack.PushFront(e.Value)
	}

	s.dll = newstack
}
