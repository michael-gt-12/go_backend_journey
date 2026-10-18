package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
)

func add(num1, num2 int) int {
	return num1 + num2
}

// function can return any number of values
func swap(a, b string) (string, string) {
	return b, a
}

// the return type of this function is called naked, because the return clause is empty and the returned values are named on the top
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return
}

var i, j int = 1, 2

var (
	ToBe   bool       = false
	MaxInt uint64     = 1<<64 - 1
	z      complex128 = cmplx.Sqrt(-5 + 12i)
)

func main() {
	fmt.Println("my number is: ", rand.Intn(10))
	// exported packages start with capital letters and unexported packages start with small letter, unexported packages are not accessible outside the package
	fmt.Println("5 added to 7: ", add(5, 7))
	a, b := swap("hello", "world")
	fmt.Println(a, b)
	fmt.Println(split(17))
	// if the variable is initialized the type of the variable can be omitted as the variable takes the type of the initializer
	var c, python, java = true, false, "no!"
	fmt.Println(i, j, c, python, java)
	// inside a function we can use := so declare variables but only inisde a function
	name := "wakanda"
	fmt.Println(name)
	// Tyoe conversion
	var x, y int = 3, 4
	var f float64 = math.Sqrt(float64(x*x + y*y))
	var z uint = uint(f)
	fmt.Println(x, y, z, f)

	// for loop is go's for and while loop
	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Println(sum)
	// for as a while loop
	sm := 1
	for sm < 1000 {
		sm += sm
	}
	fmt.Println(sm)

	// Fizz Buzz exercise
	var n int

	fmt.Print("Enter a number: ")
	if _, err := fmt.Scan(&n); err != nil {
		fmt.Println("Please enter a valid whole number.")
		return
	}
	if n < 1 {
		fmt.Println("Please enter a number greater than 0.")
		return
	}

	for i := 1; i <= n; i++ {
		switch {
		case i%15 == 0:
			fmt.Println("FizzBuzz")
		case i%3 == 0:
			fmt.Println("Fizz")
		case i%5 == 0:
			fmt.Println("Buzz")
		default:
			fmt.Println(i)
		}
	}
}
