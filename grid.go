package grid

import (
	"fmt"
	"time"
)

type Tick struct {
	OpenTime            time.Time
	Open                float64
	High                float64
	Low                 float64
	Close               float64
	Volume              float64
	CloseTime           float64
	QuoteVolume         float64
	Count               float64
	TakerBuyVolume      float64
	TakerBuyQuoteVolume float64
	Ignore              float64
}

type Side int

const (
	Unknown Side = iota
	Buy
	Sell
)

type Order struct {
	Side     Side
	Price    float64
	Quantity float64
	QuoteQty float64
	DealTime time.Time
}

// id,price,qty,quote_qty,time,is_buyer_maker
// 4819675588,71363.0,0.191,13630.333,1711929600010,false
type Trade struct {
	Id           int
	Price        float64
	Qty          float64
	QuoteQty     float64
	Time         time.Time
	IsBuyerMaker bool
}

type Strategy interface {
	Enter(price float64)
	OnTick(*Tick)
	DumpPnl()
}

type Grid struct {
	// Direction  NEUTRAL/Short/Long
	Low    float64
	High   float64
	Number int

	InitialMargin float64
	EntryPrice    float64
	Quantity      float64 // TODO: InitialMargin * Leverage / Number

	// running
	PendingOrders list[Order]
	History       list[Order]
	LastPrice     float64
	Current       int // current position in all grids
	Margin        float64
}

func NewGrid(low, high float64, number int) *Grid {
	return &Grid{
		Low:    low,
		High:   high,
		Number: number,

		InitialMargin: 5000,
		Current:       -1,
		Quantity:      816.32654,
	}
}

func (g *Grid) Enter(price float64) {
	g.EntryPrice = price
	g.LastPrice = price

	grid := g.grid()
	side := Buy
	for i := 0; i < g.Number+1; i++ {
		price := g.Low + float64(i)*grid

		if side == Buy && price > g.LastPrice {
			side = Sell
			g.Current = i
			g.PendingOrders.Add(Order{
				Side:     Unknown,
				Price:    price,
				Quantity: g.Quantity,
				QuoteQty: g.Quantity / price,
			})
			continue
		}

		g.PendingOrders.Add(Order{
			Side:     side,
			Price:    price,
			Quantity: g.Quantity,
			QuoteQty: g.Quantity / price,
		})
	}
}

func (g *Grid) Place(side Side, n int) {
	g.PendingOrders[g.Current].Side = Unknown
	g.PendingOrders[n].Side = side
}

func (g *Grid) grid() float64 {
	return (g.High - g.Low) / float64(g.Number)
}
func (g *Grid) DumpPnl() {
	target := 0.0
	pnl := 0.0
	fee := 0.0
	for i := 0; i < len(g.History); i++ {
		o := &g.History[i]

		if o.Side == Buy {
			fee += o.Quantity * 0.0005
			target += o.QuoteQty
			pnl -= o.Quantity
		} else if o.Side == Sell {
			q := o.QuoteQty * o.Price
			target -= o.QuoteQty
			pnl += q
			fee += q * 0.0002
		}
	}
	fmt.Printf("target: %4.2f\n", target)
	fmt.Printf("fee: %4.2f\n", fee)
	last := target*g.LastPrice + pnl - fee
	fmt.Printf("pnl: %4.2f\n", last)
}

func (g *Grid) OnTick(t *Tick) {
	// look up
	if g.Current < g.Number && t.Open >= g.PendingOrders[g.Current+1].Price {
		g.History.Add(g.PendingOrders[g.Current+1])
		g.PendingOrders[g.Current+1].Side = Unknown
		g.PendingOrders[g.Current].Side = Buy
		g.Current++
		g.LastPrice = t.Open
	} else if g.Current > 0 && t.Open < g.PendingOrders[g.Current-1].Price {
		g.History.Add(g.PendingOrders[g.Current-1])
		g.PendingOrders[g.Current-1].Side = Unknown
		g.PendingOrders[g.Current].Side = Sell
		g.Current--
		g.LastPrice = t.Open
	} else {
		g.LastPrice = t.Open
	}
}

func Dump(s Strategy, cht chan Tick) {
	//
	t := <-cht
	s.Enter(t.Close)

	for t := range cht {
		s.OnTick(&t)
	}

	s.DumpPnl()
}
