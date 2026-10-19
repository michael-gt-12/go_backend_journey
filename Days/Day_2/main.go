package day2

import (
	"fmt"
)

// Pointers
// pointers hold the memory address of a value

// Structs
// It is a collection of fields
type Vertex struct {
	X int
	Y int
}

// struct literals
var (
	v1 = Vertex{1, 2}  // has type Vertex
	v2 = Vertex{X: 1}  // Y:0 is implicit
	v3 = Vertex{}      // X:0 and Y:0
	p  = &Vertex{1, 2} // has type *Vertex
)

// helper function to print slices
func printSlice(s string, x []int) {
	fmt.Printf("%s len=%d cap=%d %v\n",
		s, len(x), cap(x), x)
}

func main() {
	i, j := 42, 2701
	// for i
	// p holds the address of i , which means it points to i
	p := &i
	fmt.Println(p)
	// here we are accessing the value that is being pointed at by p , which is in this case i, it is called dereferencing or indirecting
	fmt.Println(*p)
	// here we are reassigning using indirecting
	*p = 21
	fmt.Println(i)

	// for j
	// v holds the address of j which means it points to j
	v := &j
	// accessing j through its memory address that is placed inside v
	*v = *v / 37
	// printting v shows the memory address
	fmt.Println(v)
	fmt.Println(j)
	// Structs
	fmt.Println(Vertex{1, 2})
	fmt.Println(Vertex{1, 2}.X)
	fmt.Println(Vertex{1, 2}.Y)
	// structs and pointers
	ver := Vertex{10, 20}
	poin := &ver
	(*poin).X = 11
	fmt.Println(ver, (*poin))
	// a short hand for accessing the elements of a struct using the pointer of the struct
	poin.X = 1e9
	poin.Y = 17
	fmt.Println(*poin)
	// struct literals
	fmt.Println(v1, p, v2, v3)
	// Arrays
	// in Go arrays can not be resized
	var a [2]string
	a[0] = "Hello"
	a[1] = "World"
	fmt.Println(a[0], a[1])
	fmt.Println(a)

	primes := [6]int{2, 3, 5, 7, 11, 13}
	fmt.Println(primes)
	// slices
	// is a dynamically-sized, flexible view into the elements of an array. In practice, slices are much more common than arrays
	// changing the elements of a slice changes the elements of the underlying array, and other slices which share the underlying array will also be changed
	var s1 []int = primes[1:4]
	s1[0] = 10
	fmt.Println(s1)
	fmt.Println(primes)

	// Array literals and slice literals
	q := []int{2, 3, 5, 7, 11, 13}
	fmt.Println(q)

	r := []bool{true, false, true, true, false, true}
	fmt.Println(r)

	s := []struct {
		i int
		b bool
	}{
		{2, true},
		{3, false},
		{5, true},
		{7, true},
		{11, false},
		{13, true},
	}
	fmt.Println(s)
	//  creating slices with make
	//  this is how you create dynamically-sized arrays

	e := make([]int, 5)
	printSlice("e", e)

	b := make([]int, 0, 5)
	printSlice("b", b)

	c := b[:2]
	printSlice("c", c)

	d := c[2:5]
	printSlice("d", d)
}
