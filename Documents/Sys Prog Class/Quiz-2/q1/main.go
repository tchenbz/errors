//interfaces
package main

import (
	"fmt"
)

//struct to implement the interface
type Values struct {
	radius float64
	height float64
}

//interface
type Cup interface {
	volume() float64
}

//implement volume()
func (v Values) volume() float64 {
	return 3.14 * v.radius * v.radius * v.height
}

//access method
func calculate(c Cup) {
	fmt.Println("Volume of cup is: ", c.volume)
}

func main() {
	val := Values{5, 8}
	calculate(val)
}