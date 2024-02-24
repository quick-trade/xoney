package dataread

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/quick-trade/xoney/common/data"
)

func LoadChartFromCSV(filePath string, tf data.TimeFrame, contains_index int) (data.Chart, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return data.Chart{}, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()
	if err != nil {
		return data.Chart{}, fmt.Errorf("error reading CSV: %w", err)
	}

	chart := data.Chart{
		Open: make([]float64, 0, 1000),
		High: make([]float64, 0, 1000),
		Low: make([]float64, 0, 1000),
		Close: make([]float64, 0, 1000),
		Volume: make([]float64, 0, 1000),
		Timestamp: data.NewTimeStamp(tf, 1000),
	}

	for _, row := range rows[1:] {
		timestampStr := row[contains_index]
		openStr := row[1+contains_index]
		highStr := row[2+contains_index]
		lowStr := row[3+contains_index]
		closeStr := row[4+contains_index]
		volumeStr := row[5+contains_index]

		timestamp, err := time.Parse("2006-01-02 15:04:05", timestampStr)
		if err != nil {
			return data.Chart{}, fmt.Errorf("error parsing time: %w", err)
		}

		open, err := strconv.ParseFloat(openStr, 64)
		if err != nil {
			return data.Chart{}, fmt.Errorf("error parsing: %w", err)
		}

		high, err := strconv.ParseFloat(highStr, 64)
		if err != nil {
			return data.Chart{}, fmt.Errorf("error parsing: %w", err)
		}

		low, err := strconv.ParseFloat(lowStr, 64)
		if err != nil {
			return data.Chart{}, fmt.Errorf("error parsing: %w", err)
		}

		closePrice, err := strconv.ParseFloat(closeStr, 64)
		if err != nil {
			return data.Chart{}, fmt.Errorf("error parsing: %w", err)
		}

		volume, err := strconv.ParseFloat(volumeStr, 64)
		if err != nil {
			return data.Chart{}, fmt.Errorf("error parsing: %w", err)
		}

		chart.Timestamp.Append(timestamp)
		chart.Open = append(chart.Open, open)
		chart.High = append(chart.High, high)
		chart.Low = append(chart.Low, low)
		chart.Close = append(chart.Close, closePrice)
		chart.Volume = append(chart.Volume, volume)
	}

	return chart, nil
}

func WriteMap(data_map map[data.Currency][]float64, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	columns := make([]string, 0, 10)
	currencies := make([]data.Currency, 0, 10)

	for k := range data_map {
		columns = append(columns, k.String())
		currencies = append(currencies, k)
	}

	if err := writer.Write(columns); err != nil {
		return err
	}

	maxLen := maxRowLength(data_map)

	for i := 0; i < maxLen; i++ {
		row := make([]string, 0, 10)
		for _, currency := range currencies {
			val := 0.0
			if i < len(data_map[currency]){
				val = data_map[currency][i]
			}

			row = append(row, strconv.FormatFloat(val, 'f', -1, 64))
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func maxRowLength(data map[data.Currency][]float64) int {
	maxLength := 0
	for _, values := range data {
		if len(values) > maxLength {
			maxLength = len(values)
		}
	}
	return maxLength
}
