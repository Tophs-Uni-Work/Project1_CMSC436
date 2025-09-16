package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type DataPoint struct {
	X, Y  float64
	Class int
}

func readData(filename string) ([]DataPoint, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []DataPoint
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			continue
		}

		x, err1 := strconv.ParseFloat(parts[0], 64)
		y, err2 := strconv.ParseFloat(parts[1], 64)
		class, err3 := strconv.Atoi(parts[2])

		if err1 != nil || err2 != nil || err3 != nil {
			continue
		}

		data = append(data, DataPoint{X: x, Y: y, Class: class})
	}

	return data, scanner.Err()
}

func normalizeData(data []DataPoint) []DataPoint {
	if len(data) == 0 {
		return data
	}

	// Find min and max for X and Y
	minX, maxX := data[0].X, data[0].X
	minY, maxY := data[0].Y, data[0].Y

	for _, point := range data {
		minX = math.Min(minX, point.X)
		maxX = math.Max(maxX, point.X)
		minY = math.Min(minY, point.Y)
		maxY = math.Max(maxY, point.Y)
	}

	// Normalize using formula: (value - min) / (max - min)
	var normalized []DataPoint
	for _, point := range data {
		normalizedX := (point.X - minX) / (maxX - minX)
		normalizedY := (point.Y - minY) / (maxY - minY)
		normalized = append(normalized, DataPoint{
			X:     normalizedX,
			Y:     normalizedY,
			Class: point.Class,
		})
	}

	return normalized
}

func writeData(filename string, data []DataPoint) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, point := range data {
		_, err := fmt.Fprintf(writer, "%.6f,%.6f,%d\n", point.X, point.Y, point.Class)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	dataDir := "Data"
	outputDir := "Data/normalized"

	// Create output directory
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatal("Failed to create output directory:", err)
	}

	// Process each group file
	groups := []string{"groupA.txt", "groupB.txt", "groupC.txt"}

	for _, group := range groups {
		inputPath := filepath.Join(dataDir, group)
		outputPath := filepath.Join(outputDir, group)

		fmt.Printf("Processing %s...\n", group)

		// Read data
		data, err := readData(inputPath)
		if err != nil {
			log.Printf("Failed to read %s: %v\n", inputPath, err)
			continue
		}

		fmt.Printf("  Read %d data points\n", len(data))

		// Normalize data
		normalized := normalizeData(data)

		// Write normalized data
		err = writeData(outputPath, normalized)
		if err != nil {
			log.Printf("Failed to write %s: %v\n", outputPath, err)
			continue
		}

		fmt.Printf("  Normalized data written to %s\n", outputPath)
	}

	fmt.Println("Data normalization complete!")
}