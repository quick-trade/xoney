package dataread

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
	"xoney/common/data"
)

func LoadChartFromCSV(filePath string) (data.Chart, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return data.Chart{}, fmt.Errorf("Error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()
	if err != nil {
		return data.Chart{}, fmt.Errorf("Error reading CSV: %w", err)
	}

	var chart data.Chart

	for _, row := range rows[1:] {
		timestampStr := row[0]
		openStr := row[1]
		highStr := row[2]
		lowStr := row[3]
		closeStr := row[4]
		volumeStr := row[5]

		timestamp, err := time.Parse("2006-01-02 15:04:05", timestampStr)
		if err != nil {
			return data.Chart{}, fmt.Errorf("Error parsing time: %w", err)
		}

		open, err := strconv.ParseFloat(openStr, 64)
		if err != nil {
			return data.Chart{}, fmt.Errorf("Error parsing: %w", err)
		}

		high, err := strconv.ParseFloat(highStr, 64)
		if err != nil {
			return data.Chart{}, fmt.Errorf("Error parsing: %w", err)
		}

		low, err := strconv.ParseFloat(lowStr, 64)
		if err != nil {
			return data.Chart{}, fmt.Errorf("Error parsing: %w", err)
		}

		closePrice, err := strconv.ParseFloat(closeStr, 64)
		if err != nil {
			return data.Chart{}, fmt.Errorf("Error parsing: %w", err)
		}

		volume, err := strconv.ParseFloat(volumeStr, 64)
		if err != nil {
			return data.Chart{}, fmt.Errorf("Error parsing: %w", err)
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

func WriteSlice(data []float64, columnName string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{columnName}); err != nil {
		return err
	}

	for _, value := range data {
		row := []string{fmt.Sprintf("%f", value)}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
