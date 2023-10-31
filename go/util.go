package main

import "math"

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
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
