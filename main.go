package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type StockData struct {
	Date  time.Time
	Close float64
}

func main() {
	// Open the CSV file
	file, err := os.Open("ksa_us_d.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	// Parse data into StockData slice
	var data []StockData
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		if len(record) < 5 {
			fmt.Printf("Invalid row %d: %v\n", i+1, record)
			continue
		}
		// Parse Date (assuming format "YYYY-MM-DD")
		date, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			fmt.Printf("Error parsing date in row %d: %v\n", i+1, err)
			continue
		}
		// Parse Close price
		closePrice, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			fmt.Printf("Error parsing close price in row %d: %v\n", i+1, err)
			continue
		}
		data = append(data, StockData{Date: date, Close: closePrice})
	}

	// Calculate 20-day SMA and detect crossovers
	const window = 20
	for i := 0; i < len(data); i++ {
		if i < window-1 {
			continue // Need at least 20 days for SMA
		}
		// Calculate SMA
		sum := 0.0
		for j := i - window + 1; j <= i; j++ {
			sum += data[j].Close
		}
		sma := sum / float64(window)

		// Detect crossovers
		signal := ""
		if i > window-1 { // Ensure we can compare with previous day
			currentClose := data[i].Close
			previousClose := data[i-1].Close
			previousSMA := 0.0
			for j := i - window; j < i; j++ {
				previousSMA += data[j].Close
			}
			previousSMA /= float64(window)

			if previousClose <= previousSMA && currentClose > sma {
				signal = "BUY"
			} else if previousClose >= previousSMA && currentClose < sma {
				signal = "SELL"
			}
		}

		// Print results
		fmt.Printf("Date: %s, Close: %.2f, SMA(20): %.2f, Signal: %s\n",
			data[i].Date.Format("2006-01-02"), data[i].Close, sma, signal)
	}
}
