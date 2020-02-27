package dice

import (
	"fmt"
)

// ErrNumberTooLow returns an error indicating that the number of dice to roll
// is too low.
func ErrNumberTooLow(num int) error {
	return fmt.Errorf("number of dice is too low: %d", num)
}

// ErrSidesTooLow returns an error indicating that the number of sides on the
// dice is too low.
func ErrSidesTooLow(num int) error {
	return fmt.Errorf("number of sides is too low: %d", num)
}

// ErrDropLowTooLow returns an error indicating that the number of low dice
// to drop is less than zero.
func ErrDropLowTooLow(num int) error {
	return fmt.Errorf("number of low dice to drop must be positive: %d", num)
}

// ErrDropHighTooLow returns an error indicating that the number of high dice
// to drop is less than zero.
func ErrDropHighTooLow(num int) error {
	return fmt.Errorf("number of high dice to drop must be positive: %d", num)
}

// ErrTooManyDropped returns an error indicating that too many dice in a dice
// specification would be dropped to return a result.
func ErrTooManyDropped(low, high, num int) error {
	return fmt.Errorf("too many dice dropped: %d + %d >= %d", low, high, num)
}
