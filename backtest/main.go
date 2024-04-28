package main

import (
	"fmt"
	"grid"
)

func main() {
	for n := 30; n < 100; n += 5 {
		d := grid.OpenCSV("data/BTCUSDT-1m-2024-01.csv")

		t := grid.NewGrid(40000, 43000, n)
		pnl, err := grid.Run(t, d)
		fmt.Printf("%02d %8.02f %8.02f %v\n", n, pnl.Fee, pnl.Profit, err)

		d.Close()
	}
}
