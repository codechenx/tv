package main

import (
	"sort"

	"github.com/montanaflynn/stats"
)

type statsSummary interface {
	summary(a []string)
	getSummaryData() [][]string
	getSummaryStr(a []string) string
}

type DiscreteStats struct {
	data        []string
	summaryData [][]string
	count       int
	unique      int
	missing     int
	counter     map[string]int
}

type ContinuousStats struct {
	data        []float64
	summaryData [][]string
	count       int
	min         float64
	max         float64
	mean        float64
	median      float64
	sd          float64
	variance    float64
	sum         float64
	q1          float64
	q2          float64
	q3          float64
	iqr         float64
	mode        float64
	modeCount   int
	missing     int
}

func (s *ContinuousStats) summary(a []string) {
	originalCount := len(a)
	data := stats.LoadRawData(a)
	s.data = data
	s.count = len(data)
	s.missing = originalCount - s.count

	if s.count == 0 {
		s.summaryData = [][]string{{"Total values", I2S(originalCount)}, {"Missing/Invalid", I2S(s.missing)}, {"No numeric data", ""}}
		return
	}

	s.min, _ = stats.Min(data)
	s.max, _ = stats.Max(data)
	s.mean, _ = stats.Mean(data)
	s.median, _ = stats.Median(data)
	s.sd, _ = stats.StandardDeviation(data)
	s.variance, _ = stats.Variance(data)
	s.sum, _ = stats.Sum(data)

	q, _ := stats.Quartile(data)
	s.q1, s.q2, s.q3 = q.Q1, q.Q2, q.Q3
	s.iqr = s.q3 - s.q1

	// Calculate mode
	s.mode, s.modeCount = calculateMode(data)

	summaryArray := [][]string{
		{"Total values", I2S(originalCount)},
		{"Valid numbers", I2S(s.count)},
		{"Missing/Invalid", I2S(s.missing)},
		{"", ""},
		{"Min", F2S(s.min)},
		{"Max", F2S(s.max)},
		{"Range", F2S(s.max - s.min)},
		{"Sum", F2S(s.sum)},
		{"", ""},
		{"Mean", F2S(s.mean)},
		{"Median", F2S(s.median)},
		{"Mode", F2S(s.mode) + " (" + I2S(s.modeCount) + "x)"},
		{"", ""},
		{"Std Dev", F2S(s.sd)},
		{"Variance", F2S(s.variance)},
		{"", ""},
		{"Q1 (25%)", F2S(s.q1)},
		{"Q2 (50%)", F2S(s.q2)},
		{"Q3 (75%)", F2S(s.q3)},
		{"IQR", F2S(s.iqr)},
	}
	s.summaryData = summaryArray
}

func calculateMode(data []float64) (float64, int) {
	if len(data) == 0 {
		return 0, 0
	}

	freq := make(map[float64]int)
	for _, v := range data {
		freq[v]++
	}

	var mode float64
	maxCount := 0
	for k, v := range freq {
		if v > maxCount {
			maxCount = v
			mode = k
		}
	}

	return mode, maxCount
}

func (s *ContinuousStats) getSummaryData() [][]string {
	return s.summaryData
}

func (s *ContinuousStats) getSummaryStr(a []string) string {
	result := ""
	s.summary(a)
	summaryArray := s.getSummaryData()

	for _, i := range summaryArray {
		var n, v string
		n, v = i[0], i[1]
		result = result + "#" + n + " : " + v + "\n"
	}

	return result
}

func (s *DiscreteStats) summary(a []string) {
	s.data = a
	s.count = len(a)

	// Count empty/missing values
	s.missing = 0
	s.counter = make(map[string]int)

	for _, row := range a {
		if row == "" {
			s.missing++
		}
		s.counter[row]++
	}

	s.unique = len(s.counter)

	type kv struct {
		Key   string
		Value int
	}

	// Sort map by value (frequency)
	var ss []kv
	for k, v := range s.counter {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	// Build summary with overview first
	s.summaryData = [][]string{
		{"Total values", I2S(s.count)},
		{"Unique values", I2S(s.unique)},
		{"Empty/Missing", I2S(s.missing)},
		{"", ""},
		{"Value", "Frequency"},
		{"─────", "─────────"},
	}

	// Add top values (limit to prevent excessive display)
	maxDisplay := 50
	for i, kv := range ss {
		if i >= maxDisplay {
			remaining := len(ss) - maxDisplay
			s.summaryData = append(s.summaryData, []string{"...", "(" + I2S(remaining) + " more)"})
			break
		}

		displayKey := kv.Key
		if displayKey == "" {
			displayKey = "(empty)"
		}
		// Truncate very long values
		if len(displayKey) > 40 {
			displayKey = displayKey[:37] + "..."
		}

		percent := float64(kv.Value) / float64(s.count) * 100
		s.summaryData = append(s.summaryData, []string{
			displayKey,
			I2S(kv.Value) + " (" + F2S(percent) + "%)",
		})
	}
}
func (s *DiscreteStats) getSummaryData() [][]string {
	return s.summaryData
}

func (s *DiscreteStats) getSummaryStr(a []string) string {
	s.summary(a)
	summaryArray := s.getSummaryData()
	result := ""
	for _, i := range summaryArray {
		var n, v string
		n, v = i[0], i[1]
		result = result + "#" + n + " : " + v + "\n"
	}
	result = result + "----------\n" + "Top 20 variable\n\n"
	type kv struct {
		Key   string
		Value int
	}

	//sortByStr map by value
	var ss []kv
	for k, v := range s.counter {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	for _, kv := range ss {
		result = result + "#" + kv.Key + " : " + I2S(kv.Value) + "\n"
	}

	return result
}
