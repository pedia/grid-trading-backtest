package main

import (
	"fmt"
	"grid"
)

func main() {
	// d := grid.OpenCSV("data/BTCUSDT-1m-2024-04-19.csv")
	// t := grid.NewGrid(60000, 66000, 43)
	// pnl, err := grid.Run(t, d)
	// fmt.Printf(" - %6.0f %02d %8.0f %8.0f %v\n", 66000., 43, pnl.Fee, pnl.Profit, err)

	// d := grid.OpenCSV("data/BTCUSDT-1m-2024-01.csv")
	// t := grid.NewGrid(40000, up, n)

	// return
	// variables: up, number
	for up := 62000.0; up < 66000; up += 400 {
		for n := 20; n < 160; n += 3 {
			d := grid.OpenCSV("data/BTCUSDT-1m-2024-04-19.csv")

			t := grid.NewGrid(40000, up, n)
			pnl, err := grid.Run(t, d)
			// if err != nil || pnl.Profit < 10 {
			// 	continue
			// }
			_ = err
			// if pnl.Profit > 1000 {
			//            low, hi  n    seg   qty  fee    profit
			fmt.Printf(" %.0f,%.0f \t%02d \t%3.0f \t%5.0f \t%8.02f \t%8.0f \t%2d\n", t.Low, t.High,
				n, (t.High-t.Low)/float64(t.Number),
				t.Quantity,
				pnl.Fee, pnl.Profit, pnl.TradeCount)
			// }
			// if pnl.Profit > 2000 {
			// 	for _, o := range t.History {
			// 		fmt.Printf("%s %.0f %v\n", o.DealTime, o.Price, o.Side)
			// 	}
			// }

			d.Close()
		}
	}
}

func foo() {
	for up := 42000.0; up < 60000; up += 1000 {
		for n := 30; n < 100; n += 5 {
			d := grid.OpenCSV("data/BTCUSDT-1m-2024-01.csv")

			t := grid.NewGrid(40000, up, n)
			pnl, err := grid.Run(t, d)
			if pnl.Profit < 10 {
				continue
			}
			fmt.Printf(" - %6.0f %02d %8.02f %8.02f %v\n", up, n, pnl.Fee, pnl.Profit, err)

			d.Close()
		}
	}
}
