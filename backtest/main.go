package main

import "grid"

func main() {
	ct := make(chan grid.Tick, 1)
	go grid.Read("data/BTCUSDT-1m-2024-04-19.csv", ct)

	t := grid.NewGrid(60000, 66000, 30)
	grid.Dump(t, ct)
}
