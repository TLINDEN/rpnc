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
	"fmt"
	"math"
	"strings"
)

// find an item in a list, generic variant
func contains[E comparable](s []E, v E) bool {
	for _, vs := range s {
		if v == vs {
			return true
		}
	}

	return false
}

// look if a key in a map exists, generic variant
func exists[K comparable, V any](m map[K]V, v K) bool {
	if _, ok := m[v]; ok {
		return true
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

func list2str(list Numbers) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(list)), " "), "[]")
}

func Error(m string) error {
	return fmt.Errorf("Error: %s", m)
}
