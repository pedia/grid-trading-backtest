package grid

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func ensure_int64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return i
	}
	panic(fmt.Errorf("%s to int failed %s", s, err))
}
func ensure_float(s string) float64 {
	i, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return i
	}
	panic(fmt.Errorf("%s to float failed %s", s, err))
}

func OrderFrom(r []string) Tick {
	return Tick{
		time.UnixMilli(ensure_int64(r[0])),
		ensure_float(r[1]),
		ensure_float(r[2]),
		ensure_float(r[3]),
		ensure_float(r[4]),
		ensure_float(r[5]),
		ensure_float(r[6]),
		ensure_float(r[7]),
		ensure_float(r[8]),
		ensure_float(r[9]),
		ensure_float(r[10]),
		ensure_float(r[11]),
	}
}

// read csv file to chan
// read all records is better?
func Read(filename string, co chan Tick) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(f)
	reader.ReuseRecord = true

	// skip header
	reader.Read()

	for {

		r, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		co <- OrderFrom(r)
	}

	close(co)
}

type DataBase interface {
	FetchOne() (Tick, error)
	Close()
}

type csv_reader struct {
	reader *csv.Reader
}

func OpenCSV(filename string) *csv_reader {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(f)
	reader.ReuseRecord = true

	// skip header
	reader.Read()

	return &csv_reader{reader: reader}
}

// io.EOF
func (c *csv_reader) FetchOne() (Tick, error) {
	r, err := c.reader.Read()
	if err != nil {
		return Tick{}, err
	}
	return OrderFrom(r), nil
}
func (c *csv_reader) Close() {
	// close(c.reader)
}
