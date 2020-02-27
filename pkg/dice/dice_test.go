package dice

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	bigMaxInt = big.NewInt(math.MaxInt64)
)

// A Mock of the Rand interface which returns a constant value.
type mockMaxRand struct{}

// Return the MockRand as an int
func (r mockMaxRand) Intn(n int) int {
	// Always return the max value allowed
	return n - 1
}

type mockIterRand struct {
	Value int
}

func (r *mockIterRand) Intn(n int) int {
	v := r.Value % n
	r.Value++
	return v
}

func TestDice(t *testing.T) {
	t.Run("New", testDiceNew)
	t.Run("NewExt", testDiceNewExt)
	t.Run("Roll", testDiceRoll)
	t.Run("RollRand", testDiceRollRand)
	t.Run("RollResults", testDiceRollResults)
	t.Run("String", testDiceString)
	t.Run("Valiate", testDiceValidate)
}

func testDiceNew(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		d, err := New(2, 20)
		require.NoError(t, err)
		require.Equal(t, &Dice{
			Number: 2,
			Sides:  20,
		}, d)
	})
	t.Run("InvalidNumber", func(t *testing.T) {
		d, err := New(0, 20)
		require.EqualError(t, err, ErrNumberTooLow(0).Error())
		require.Nil(t, d)
	})
	t.Run("InvalidSides", func(t *testing.T) {
		d, err := New(2, 1)
		require.EqualError(t, err, ErrSidesTooLow(1).Error())
		require.Nil(t, d)
	})
}

func testDiceNewExt(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		d, err := NewExt(2, 20, 0, 0)
		require.NoError(t, err)
		require.Equal(t, &Dice{
			Number: 2,
			Sides:  20,
		}, d)
	})
	t.Run("InvalidNumber", func(t *testing.T) {
		d, err := NewExt(0, 20, 0, 0)
		require.EqualError(t, err, ErrNumberTooLow(0).Error())
		require.Nil(t, d)
	})
	t.Run("InvalidSides", func(t *testing.T) {
		d, err := NewExt(2, 1, 0, 0)
		require.EqualError(t, err, ErrSidesTooLow(1).Error())
		require.Nil(t, d)
	})
	t.Run("InvalidDropLow", func(t *testing.T) {
		d, err := NewExt(2, 20, -1, 0)
		require.EqualError(t, err, ErrDropLowTooLow(-1).Error())
		require.Nil(t, d)
	})
	t.Run("InvalidDropHigh", func(t *testing.T) {
		d, err := NewExt(2, 20, 0, -1)
		require.EqualError(t, err, ErrDropHighTooLow(-1).Error())
		require.Nil(t, d)
	})
}

func testDiceValidate(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		d := &Dice{
			Number: 2,
			Sides:  20,
		}
		require.NoError(t, d.Validate())
	})
	t.Run("InvalidNumber", func(t *testing.T) {
		d := &Dice{
			Number: 0,
			Sides:  20,
		}
		require.EqualError(t, d.Validate(), ErrNumberTooLow(0).Error())
	})
	t.Run("InvalidSides", func(t *testing.T) {
		d := &Dice{
			Number: 2,
			Sides:  1,
		}
		require.EqualError(t, d.Validate(), ErrSidesTooLow(1).Error())
	})
	t.Run("InvalidDropLow", func(t *testing.T) {
		d := &Dice{
			Number:  2,
			Sides:   20,
			DropLow: -1,
		}
		require.EqualError(t, d.Validate(), ErrDropLowTooLow(-1).Error())
	})
	t.Run("InvalidDropHigh", func(t *testing.T) {
		d := &Dice{
			Number:   2,
			Sides:    20,
			DropHigh: -1,
		}
		require.EqualError(t, d.Validate(), ErrDropHighTooLow(-1).Error())
	})
	t.Run("TooManyLowDropped", func(t *testing.T) {
		d := &Dice{
			Number:  2,
			Sides:   20,
			DropLow: 2,
		}
		require.EqualError(t, d.Validate(), ErrTooManyDropped(2, 0, 2).Error())
	})
	t.Run("TooManyTotalDropped", func(t *testing.T) {
		d := &Dice{
			Number:   2,
			Sides:    20,
			DropLow:  1,
			DropHigh: 1,
		}
		require.EqualError(t, d.Validate(), ErrTooManyDropped(1, 1, 2).Error())
	})
}

func testDiceRoll(t *testing.T) {
	d, err := New(2, 20)
	// Sanity check - make sure we get a return value.
	require.NoError(t, err)
	require.NotNil(t, d)

	v := d.Roll()

	// Validate we can fit it in an int64
	require.Equal(t, -1, v.Cmp(bigMaxInt), "Result too large")

	// Convert to an int64 and check the bounds
	// We know it must be between 2 and 40
	i64value := v.Int64()
	require.LessOrEqual(t, i64value, int64(40))
	require.GreaterOrEqual(t, i64value, int64(2))
}

func testDiceRollRand(t *testing.T) {
	d, err := New(2, 20)
	require.NoError(t, err)
	require.NotNil(t, d)

	v := d.RollRand(mockMaxRand{})
	require.Equal(t, -1, v.Cmp(bigMaxInt), "Result too large")

	i64value := v.Int64()
	require.Equal(t, int64(40), i64value)
}

func testDiceRollResults(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		d, err := New(2, 20)
		require.NoError(t, err)
		require.NotNil(t, d)

		results := d.RollResults(mockMaxRand{})
		require.Equal(t, &Results{
			Value:       big.NewInt(40),
			Raw:         []int{20, 20},
			DroppedLow:  nil,
			DroppedHigh: nil,
			Kept:        []int{20, 20},
		}, results)
	})
	t.Run("DropLow", func(t *testing.T) {
		d, err := NewExt(5, 20, 3, 0)
		require.NoError(t, err)
		require.NotNil(t, d)

		// Should generate [1, 2, 3, 4, 5] and drop [1, 2, 3]
		results := d.RollResults(&mockIterRand{})
		require.Equal(t, &Results{
			Value:       big.NewInt(9),
			Raw:         []int{1, 2, 3, 4, 5},
			DroppedLow:  []int{1, 2, 3},
			DroppedHigh: nil,
			Kept:        []int{4, 5},
		}, results)
	})
	t.Run("DropHigh", func(t *testing.T) {
		d, err := NewExt(5, 20, 0, 3)
		require.NoError(t, err)
		require.NotNil(t, d)

		// Should generate [1, 2, 3, 4, 5] and drop [3, 4, 5]
		results := d.RollResults(&mockIterRand{})
		require.Equal(t, &Results{
			Value:       big.NewInt(3),
			Raw:         []int{1, 2, 3, 4, 5},
			DroppedLow:  nil,
			DroppedHigh: []int{3, 4, 5},
			Kept:        []int{1, 2},
		}, results)
	})
	t.Run("DropBoth", func(t *testing.T) {
		d, err := NewExt(5, 20, 2, 2)
		require.NoError(t, err)
		require.NotNil(t, d)

		// Should generate [1, 2, 3, 4, 5] and drop [1, 2] and [4, 5]
		results := d.RollResults(&mockIterRand{})
		require.Equal(t, &Results{
			Value:       big.NewInt(3),
			Raw:         []int{1, 2, 3, 4, 5},
			DroppedLow:  []int{1, 2},
			DroppedHigh: []int{4, 5},
			Kept:        []int{3},
		}, results)
	})
}

func testDiceString(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		d, err := New(2, 20)
		require.NoError(t, err)
		require.NotNil(t, d)
		require.Equal(t, "2d20", d.String())
	})
	t.Run("DropLow", func(t *testing.T) {
		d, err := NewExt(2, 20, 1, 0)
		require.NoError(t, err)
		require.NotNil(t, d)
		require.Equal(t, "2d20L1", d.String())
	})
	t.Run("DropHigh", func(t *testing.T) {
		d, err := NewExt(2, 20, 0, 1)
		require.NoError(t, err)
		require.NotNil(t, d)
		require.Equal(t, "2d20H1", d.String())
	})
	t.Run("DropBoth", func(t *testing.T) {
		d, err := NewExt(3, 20, 1, 1)
		require.NoError(t, err)
		require.NotNil(t, d)
		require.Equal(t, "3d20L1H1", d.String())
	})
}
