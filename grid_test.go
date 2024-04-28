package grid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrid(t *testing.T) {
	is := assert.New(t)

	g := NewGrid(60000.0, 66000, 48)
	g.Enter(63548.6)

	equal_order := func(o Order, side Side, price float64) {
		is.Equal(side, o.Side)
		is.Equal(price, o.Price)
	}
	equal_order(g.PendingOrders[0], Buy, g.Low)
	equal_order(g.PendingOrders[g.Current-1], Buy, 63500.00)
	equal_order(g.PendingOrders[g.Current], Sell, 63750.00)
	equal_order(g.PendingOrders[len(g.PendingOrders)-1], Sell, g.High)
}
