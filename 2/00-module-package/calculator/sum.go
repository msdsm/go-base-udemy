package calculator

import "fmt"

var offset float64 = 1
var Offset float64 = 1

func Sum(a float64, b float64) float64 {
	fmt.Println(multiply(a, b))
	return a + b + offset
}
