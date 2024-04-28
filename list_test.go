package grid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	is := assert.New(t)

	a := List[int](0, 1)
	a.Add(3)
	is.Equal([]int{0, 1, 3}, a.Items())
	a.Insert(2, 2)
	is.Equal([]int{0, 1, 2, 3}, a.Items())
	a.Insert(4, 4)
	is.Equal([]int{0, 1, 2, 3, 4}, a.Items())

	a.Add(5, 6)
	a.Insert(7, 7)
	is.Equal([]int{0, 1, 2, 3, 4, 5, 6, 7}, a.Items())

	a.RemoveAt(1)
	is.Equal([]int{0, 2, 3, 4, 5, 6, 7}, a.Items())
	a.RemoveAt(1)
	is.Equal([]int{0, 3, 4, 5, 6, 7}, a.Items())
	a.RemoveAt(0)
	is.Equal([]int{3, 4, 5, 6, 7}, a.Items())
	a.RemoveAt(4)
	is.Equal([]int{3, 4, 5, 6}, a.Items())
}

func BenchmarkListInsert(b *testing.B) {
	a := List[int](1, 2)

	b.Run("Add", func(b *testing.B) {
		a.Add(3)
	})

	a = List[int](1, 2)
	b.Run("Insert", func(b *testing.B) {
		a.Insert(1, 4)
	})
}
