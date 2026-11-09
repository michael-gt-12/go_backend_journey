package day20

import (
	"errors"

	"github.com/stretchr/testify/suite"
)

var ErrDivideByZero = errors.New("Division by Zero is not allowed")

func Add(a, b int) int {
	return a + b
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0.0, ErrDivideByZero
	}
	return a / b, nil
}

type Counter struct {
	n int
}

func (c *Counter) Inc() {
	c.n++
}

type CounterSuite struct {
	suite.Suite
	counter *Counter
}
