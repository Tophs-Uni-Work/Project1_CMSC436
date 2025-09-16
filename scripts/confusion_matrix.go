package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type DataPoint struct {
	X, Y  float64
	Class int
}

type ConfusionMatrix struct {
	TruePositive  int
	FalsePositive int
	TrueNegative  int
	FalseNegative int
}

func (cm ConfusionMatrix) Precision() float64 {
	if cm.TruePositive+cm.FalsePositive == 0 {
		return 0
	}
	return float64(cm.TruePositive) / float64(cm.TruePositive+cm.FalsePositive)
}

func (cm ConfusionMatrix) Recall() float64 {
	if cm.TruePositive+cm.FalseNegative == 0 {
		return 0
	}
	return float64(cm.TruePositive) / float64(cm.TruePositive+cm.FalseNegative)
}

func (cm ConfusionMatrix) F1Score() float64 {
	precision := cm.Precision()
	recall := cm.Recall()
	if precision+recall == 0 {
		return 0
	}
	return 2 * (precision * recall) / (precision + recall)
}

func (cm ConfusionMatrix) Accuracy() float64 {
	total := cm.TruePositive + cm.FalsePositive + cm.TrueNegative + cm.FalseNegative
	if total == 0 {
		return 0
	}
	return float64(cm.TruePositive+cm.TrueNegative) / float64(total)
}

// For original data: use diagonal divider (y = -x + 70000)
func classifyOriginal(x, y float64) int {
	if y > (-x + 70000) {
		return 1 // Big car
	}
	return 0 // Small car
}

// For normalized data: use linear classifier y = -x + 1
func classifyNormalized(x, y float64) int {
	if y > (-x + 1) {
		return 1 // Big car
	}
	return 0 // Small car
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

func evaluateClassifier(data []DataPoint, isNormalized bool) ConfusionMatrix {
	var cm ConfusionMatrix

	for _, point := range data {
		var predicted int
		if isNormalized {
			predicted = classifyNormalized(point.X, point.Y)
		} else {
			predicted = classifyOriginal(point.X, point.Y)
		}
		actual := point.Class

		// TP: Predicted big car (1) and actually big car (1)
		// FP: Predicted big car (1) but actually small car (0) 
		// TN: Predicted small car (0) and actually small car (0)
		// FN: Predicted small car (0) but actually big car (1)

		if predicted == 1 && actual == 1 {
			cm.TruePositive++
		} else if predicted == 1 && actual == 0 {
			cm.FalsePositive++
		} else if predicted == 0 && actual == 0 {
			cm.TrueNegative++
		} else if predicted == 0 && actual == 1 {
			cm.FalseNegative++
		}
	}

	return cm
}

func writeConfusionMatrixData(filename string, cm ConfusionMatrix) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write confusion matrix as 2x2 matrix for gnuplot
	// Format: row col value
	fmt.Fprintf(file, "0 0 %d\n", cm.TruePositive)   // Actual Big, Predicted Big
	fmt.Fprintf(file, "0 1 %d\n", cm.FalseNegative)  // Actual Big, Predicted Small
	fmt.Fprintf(file, "1 0 %d\n", cm.FalsePositive)  // Actual Small, Predicted Big
	fmt.Fprintf(file, "1 1 %d\n", cm.TrueNegative)   // Actual Small, Predicted Small

	return nil
}

func printResults(writer io.Writer, groupName string, data []DataPoint, cm ConfusionMatrix, dataType string) {
	fmt.Fprintf(writer, "\n=== %s (%s) Results ===\n", groupName, dataType)
	fmt.Fprintf(writer, "Total data points: %d\n", len(data))
	fmt.Fprintf(writer, "\nConfusion Matrix:\n")
	fmt.Fprintf(writer, "                 Predicted\n")
	fmt.Fprintf(writer, "               Big  Small\n")
	fmt.Fprintf(writer, "Actual   Big   %3d   %3d\n", cm.TruePositive, cm.FalseNegative)
	fmt.Fprintf(writer, "        Small  %3d   %3d\n", cm.FalsePositive, cm.TrueNegative)
	
	fmt.Fprintf(writer, "\nMetrics:\n")
	fmt.Fprintf(writer, "Accuracy:  %.4f\n", cm.Accuracy())
	fmt.Fprintf(writer, "Precision: %.4f\n", cm.Precision())
	fmt.Fprintf(writer, "Recall:    %.4f\n", cm.Recall())
	fmt.Fprintf(writer, "F1-Score:  %.4f\n", cm.F1Score())
}

func main() {
	groups := []string{"groupA.txt", "groupB.txt", "groupC.txt"}
	groupNames := []string{"Group A", "Group B", "Group C"}

	// Create directory for confusion matrix data (normalized only)
	os.MkdirAll("plots/normalized/confusion", 0755)

	// Create output file for results
	outputFile, err := os.Create("plots/normalized/confusion/analysis_results.txt")
	if err != nil {
		log.Fatal("Failed to create output file:", err)
	}
	defer outputFile.Close()

	// Create a multi-writer to write to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, outputFile)

	fmt.Fprintln(multiWriter, "Confusion Matrix Analysis - Normalized Data Only")
	fmt.Fprintln(multiWriter, "===============================================")

	for i, group := range groups {
		// Process normalized data only
		normalizedPath := filepath.Join("Data/normalized", group)
		if data, err := readData(normalizedPath); err == nil {
			cm := evaluateClassifier(data, true)
			printResults(multiWriter, groupNames[i], data, cm, "Normalized")
			
			// Write confusion matrix data for plotting
			matrixFile := fmt.Sprintf("plots/normalized/confusion/group%s.dat", string(rune('A'+i)))
			writeConfusionMatrixData(matrixFile, cm)
		} else {
			fmt.Fprintf(multiWriter, "Failed to read %s: %v\n", normalizedPath, err)
		}
	}

	fmt.Fprintln(multiWriter, "\n=== Classification Rules ===")
	fmt.Fprintln(multiWriter, "Original data: y > -2.22*x + 107643 → Big Car (1), y ≤ -2.22*x + 107643 → Small Car (0)")
	fmt.Fprintln(multiWriter, "Normalized data: y > -x + 1 → Big Car (1), y ≤ -x + 1 → Small Car (0)")
	fmt.Fprintln(multiWriter, "Classes: 0 = Small Car, 1 = Big Car")
	fmt.Fprintln(multiWriter, "True Positive (TP): Correctly identified big cars (class 1)")
	fmt.Fprintln(multiWriter, "Data files for plotting generated in plots/normalized/confusion/")
	fmt.Fprintln(multiWriter, "Analysis results saved to: plots/normalized/confusion/analysis_results.txt")
}