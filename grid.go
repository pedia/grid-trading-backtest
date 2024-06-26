package grid

import (
	"errors"
	"io"
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
	Quote    float64
	DealTime time.Time
}

// id,price,qty,quote_qty,time,is_buyer_maker
// 4819675588,71363.0,0.191,13630.333,1711929600010,false
type Trade struct {
	Id           int
	Price        float64
	Quantity     float64
	Quote        float64
	Time         time.Time
	IsBuyerMaker bool
}

var ErrLiquidated = errors.New("liquidated")
var ErrOutOfRange = errors.New("out of range")

type Pnl struct {
	Fee        float64
	Profit     float64
	TradeCount int
}

type Strategy interface {
	Enter(price float64)
	OnTick(*Tick) error
	GetPnl() Pnl
}

type Grid struct {
	// Direction  NEUTRAL/Short/Long
	Low    float64
	High   float64
	Number int

	InitialMargin float64
	EntryPrice    float64
	Quantity      float64 // Qty Per Order
	// TODO: InitialMargin * Leverage / Number

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
		// InitialMargin * Leverage / Number
		Quantity: float64(5000*100) / float64(number),
	}
}

func (g *Grid) Enter(entry_price float64) {
	g.EntryPrice = entry_price
	g.LastPrice = entry_price

	segment := g.grid()
	side := Buy
	for i := 0; i < g.Number+1; i++ {
		price := g.Low + float64(i)*segment

		if side == Buy && price > g.LastPrice {
			side = Sell
			g.Current = i
			g.PendingOrders.Add(Order{
				Side:     Unknown,
				Price:    price,
				Quantity: g.Quantity,
				Quote:    g.Quantity / price,
			})
			continue
		}

		g.PendingOrders.Add(Order{
			Side:     side,
			Price:    price,
			Quantity: g.Quantity,
			Quote:    g.Quantity / price,
		})
	}
}

func (g *Grid) check() bool {
	// TODO:
	// https://www.binance.com/en/support/faq/how-to-calculate-liquidation-price-of-usd%E2%93%A2-m-futures-contracts-b3c689c1f50a44cabb3a84e663b81d93

	target := 0.0
	pnl := 0.0
	fee := 0.0
	for i := 0; i < len(g.History); i++ {
		o := &g.History[i]

		if o.Side == Buy {
			fee += o.Quantity * 0.0005
			target += o.Quote
			pnl -= o.Quantity
		} else if o.Side == Sell {
			q := o.Quote * o.Price
			fee += q * 0.0002
			target -= o.Quote
			pnl += q
		}
	}

	// if current value < -0.5 of initial margin maybe liquidate
	return target*g.LastPrice+pnl-fee > -g.InitialMargin*0.5
}

func (g *Grid) grid() float64 {
	return (g.High - g.Low) / float64(g.Number)
}
func (g *Grid) GetPnl() Pnl {
	target := 0.0
	pnl := 0.0
	fee := 0.0
	for i := 0; i < len(g.History); i++ {
		o := &g.History[i]

		if o.Side == Buy {
			fee += o.Quantity * 0.0005
			target += o.Quote
			pnl -= o.Quantity
		} else if o.Side == Sell {
			q := o.Quote * o.Price
			target -= o.Quote
			pnl += q
			fee += q * 0.0002
		}
	}
	return Pnl{Fee: fee, Profit: target*g.LastPrice + pnl - fee, TradeCount: len(g.History)}
}

func (g *Grid) Trade(o Order, t *Tick) {
	o.DealTime = t.OpenTime
	g.History.Add(o)
}

func (g *Grid) OnTick(t *Tick) error {
	// look up
	if g.Current >= 0 && g.Current < g.Number && t.Open >= g.PendingOrders[g.Current+1].Price {
		g.Trade(g.PendingOrders[g.Current+1], t)
		g.PendingOrders[g.Current+1].Side = Unknown
		g.PendingOrders[g.Current].Side = Buy
		g.Current++
		g.LastPrice = t.Open
	} else if g.Current > 0 && t.Open < g.PendingOrders[g.Current-1].Price {
		g.Trade(g.PendingOrders[g.Current-1], t)
		g.PendingOrders[g.Current-1].Side = Unknown
		g.PendingOrders[g.Current].Side = Sell
		g.Current--
		g.LastPrice = t.Open
	} else {
		g.LastPrice = t.Open
	}

	if !g.check() {
		return ErrLiquidated
	}
	// fmt.Printf("%6.0f %2d %6.0f\n", t.Open, len(g.History), g.GetPnl().Profit)
	return nil
}

func Run(s Strategy, d DataBase) (Pnl, error) {
	t, err := d.FetchOne()
	if err != nil {
		return Pnl{}, err
	}
	s.Enter(t.Close)

	for {
		t, err = d.FetchOne()
		if err == io.EOF {
			err = nil
			break
		}

		if err = s.OnTick(&t); err != nil {
			break
		}
	}

	return s.GetPnl(), err
}
