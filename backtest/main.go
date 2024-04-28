package main

import (
	"fmt"
	"grid"
)

func main() {
	for up := 40000.0; up < 60000; up += 1000 {
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
