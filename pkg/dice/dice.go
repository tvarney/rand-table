package dice

import (
	"fmt"
	"math/big"
	"sort"
)

// Dice is a struct representing a set of similar dice which may be rolled.
type Dice struct {
	Number   int
	Sides    int
	DropLow  int
	DropHigh int
}

// Results is the set of results from a dice roll.
type Results struct {
	Value       *big.Int
	DroppedLow  []int
	DroppedHigh []int
	Kept        []int
	Raw         []int
}

// New creates a new simple Dice instance which consists of just a number of
// dice.
func New(num, sides int) (*Dice, error) {
	return NewExt(num, sides, 0, 0)
}

// NewExt create a new Dice instance which will roll a given number of dice
// and drop some low and/or high dice.
func NewExt(num, sides, droplow, drophigh int) (*Dice, error) {
	d := &Dice{
		Number:   num,
		Sides:    sides,
		DropLow:  droplow,
		DropHigh: drophigh,
	}
	err := d.Validate()
	if err != nil {
		return nil, err
	}
	return d, nil
}

// String returns a string representation of this Dice instance.
func (d *Dice) String() string {
	if d.DropLow > 0 {
		if d.DropHigh > 0 {
			return fmt.Sprintf("%dd%dL%dH%d", d.Number, d.Sides, d.DropLow, d.DropHigh)
		}
		return fmt.Sprintf("%dd%dL%d", d.Number, d.Sides, d.DropLow)
	}
	if d.DropHigh > 0 {
		return fmt.Sprintf("%dd%dH%d", d.Number, d.Sides, d.DropHigh)
	}
	return fmt.Sprintf("%dd%d", d.Number, d.Sides)
}

// Validate checks that a Dice instance is valid.
//
// If the instance is not valid, an error describing the problem is returned.
// The error returned is the first error encountered.
//
// This is called automatically by the New() and NewExt() methods - if the
// instance was deserialized from json or yaml, then you should call this
// method manually to ensure calling one of the Roll*() methods won't cause a
// panic.
func (d *Dice) Validate() error {
	if d.Number < 1 {
		return ErrNumberTooLow(d.Number)
	}
	if d.Sides <= 1 {
		return ErrSidesTooLow(d.Sides)
	}
	if d.DropLow < 0 {
		return ErrDropLowTooLow(d.DropLow)
	}
	if d.DropHigh < 0 {
		return ErrDropHighTooLow(d.DropHigh)
	}
	if d.DropLow+d.DropHigh >= d.Number {
		return ErrTooManyDropped(d.DropLow, d.DropHigh, d.Number)
	}

	return nil
}

// Roll simulates a dice roll as specified by the Dice instance.
//
// This function is equivalent to calling `d.RollRand(StdRand)`.
func (d *Dice) Roll() *big.Int {
	return d.RollRand(StdRand)
}

// RollRand simulates a dice roll as specified by the Dice instance with the
// given Rand implementation.
func (d *Dice) RollRand(r Rand) *big.Int {
	return d.RollResults(r).Value
}

// RollResults simulates a dice roll as specified by the Dice instance with
// the given Rand implementation and returns detailed results.
//
// This function assumes that the Dice instance is valid - if it isn't, the
// roll may cause a panic (e.g. if you have negatives somewhere).
func (d *Dice) RollResults(r Rand) *Results {
	results := &Results{
		Value: big.NewInt(0),
		Raw:   make([]int, 0, d.Number),
	}

	for i := 0; i < d.Number; i++ {
		results.Raw = append(results.Raw, r.Intn(d.Sides)+1)
	}

	sort.Ints(results.Raw)

	if d.DropLow > 0 {
		results.DroppedLow = results.Raw[:d.DropLow]
	}
	if d.DropHigh > 0 {
		results.DroppedHigh = results.Raw[d.Number-d.DropHigh:]
	}
	results.Kept = results.Raw[d.DropLow : d.Number-d.DropHigh]

	value := results.Value
	for _, v := range results.Kept {
		value = value.Add(value, big.NewInt(int64(v)))
	}
	results.Value = value

	return results
}
