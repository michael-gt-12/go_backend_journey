package day20

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestDivide_AssertVsRequire(t *testing.T) {
	got, err := Divide(10, 0)
	// stops the test if there is an error
	require.NoError(t, err)
	// asserts that the result should be 5
	assert.Equal(t, 5, got)
	// asserts that the result should be different from Zero
	assert.NotZero(t, got)

}

// Table driven test used to Test different cases in one function rather than repting a function for different cases
func TestDivide_Table(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    int
		wantErr error
	}{
		{name: "normal division", a: 10, b: 2, want: 5},
		{name: "rounds down", a: 7, b: 2, want: 3},
		{name: "zero numerator", a: 0, b: 5, want: 0},
		{name: "divide by zero", a: 5, b: 0, wantErr: ErrDivideByZero},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Divide(test.a, test.b)
			if test.wantErr != nil {
				assert.ErrorIs(t, err, test.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})

	}
}

// test suites are used to gropu related tests in one place
// runs before EACH test method
func (s *CounterSuite) SetupTest() {
	s.counter = &Counter{}
}

func (s *CounterSuite) TestStartsAtZero() {
	s.Equal(0, s.counter.n)
}

func (s *CounterSuite) TestInc() {
	s.counter.Inc()
	s.Equal(1, s.counter.n)
}

// proves SetupTest gave a fresh counter: it is not 2 here
func (s *CounterSuite) TestFreshStateEachTime() {
	s.counter.Inc()
	s.Equal(1, s.counter.n)
}

// the entry point go test finds
func TestCounterSuite(t *testing.T) {
	suite.Run(t, new(CounterSuite))
}
