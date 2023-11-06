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
	"testing"
)

func TestPush(t *testing.T) {
	t.Run("push", func(t *testing.T) {
		s := NewStack()
		s.Push(5)

		if s.linklist.Back().Value != 5.0 {
			t.Errorf("push failed:\n+++  got: %f\n--- want: %f",
				s.linklist.Back().Value, 5.0)
		}
	})
}

func TestPop(t *testing.T) {
	t.Run("pop", func(t *testing.T) {
		s := NewStack()
		s.Push(5)
		got := s.Pop()

		if got != 5.0 {
			t.Errorf("pop failed:\n+++  got: %f\n--- want: %f",
				got, 5.0)
		}

		if s.Len() != 0 {
			t.Errorf("stack not empty after pop()")
		}
	})
}

func TestPops(t *testing.T) {
	t.Run("pops", func(t *testing.T) {
		s := NewStack()
		s.Push(5)
		s.Push(5)
		s.Push(5)
		s.Pop()

		if s.Len() != 2 {
			t.Errorf("stack len not correct after pop:\n+++  got: %d\n--- want: %d",
				s.Len(), 2)
		}
	})
}

func TestShift(t *testing.T) {
	t.Run("shift", func(t *testing.T) {
		s := NewStack()
		s.Shift()

		if s.Len() != 0 {
			t.Errorf("stack not empty after shift()")
		}
	})
}

func TestClear(t *testing.T) {
	t.Run("clear", func(t *testing.T) {
		s := NewStack()
		s.Push(5)
		s.Push(5)
		s.Push(5)
		s.Clear()

		if s.Len() != 0 {
			t.Errorf("stack not empty after clear()")
		}
	})
}

func TestLast(t *testing.T) {
	t.Run("last", func(t *testing.T) {
		s := NewStack()
		s.Push(5)
		got := s.Last()

		if len(got) != 1 {
			t.Errorf("last failed:\n+++  got: %d elements\n--- want: %d elements",
				len(got), 1)
		}

		if got[0] != 5.0 {
			t.Errorf("last failed:\n+++  got: %f\n--- want: %f",
				got, 5.0)
		}

		if s.Len() != 1 {
			t.Errorf("stack modified after last()")
		}
	})
}

func TestAll(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		s := NewStack()
		list := []float64{2, 4, 6, 8}

		for _, item := range list {
			s.Push(item)
		}

		got := s.All()

		if len(got) != len(list) {
			t.Errorf("all failed:\n+++  got: %d elements\n--- want: %d elements",
				len(got), len(list))
		}

		for i := 1; i < len(list); i++ {
			if got[i] != list[i] {
				t.Errorf("all failed (element %d):\n+++  got: %f\n--- want: %f",
					i, got[i], list[i])
			}
		}

		if s.Len() != len(list) {
			t.Errorf("stack modified after last()")
		}
	})
}

func TestBackupRestore(t *testing.T) {
	t.Run("shift", func(t *testing.T) {
		s := NewStack()
		s.Push(5)
		s.Backup()
		s.Clear()
		s.Restore()

		if s.Len() != 1 {
			t.Errorf("stack not correctly restored()")
		}

		a := s.Pop()
		if a != 5.0 {
			t.Errorf("stack not identical to old revision:\n+++  got: %f\n--- want: %f",
				a, 5.0)
		}
	})
}

func TestReverse(t *testing.T) {
	t.Run("reverse", func(t *testing.T) {
		s := NewStack()
		list := []float64{2, 4, 6}
		reverse := []float64{6, 4, 2}

		for _, item := range list {
			s.Push(item)
		}

		s.Reverse()

		got := s.All()

		if len(got) != len(list) {
			t.Errorf("all failed:\n+++  got: %d elements\n--- want: %d elements",
				len(got), len(list))
		}

		for i := 1; i < len(reverse); i++ {
			if got[i] != reverse[i] {
				t.Errorf("reverse failed (element %d):\n+++  got: %f\n--- want: %f",
					i, got[i], list[i])
			}
		}
	})
}
