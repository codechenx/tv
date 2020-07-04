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
	//count       int
	counter map[string]int
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
	q1          float64
	q2          float64
	q3          float64
}

func (s *ContinuousStats) summary(a []string) {
	s.count = len(a)
	data := stats.LoadRawData(a)
	s.data = data
	s.min, _ = stats.Min(data)
	s.max, _ = stats.Max(data)
	s.mean, _ = stats.Mean(data)
	s.median, _ = stats.Median(data)
	s.sd, _ = stats.StandardDeviation(data)
	q, _ := stats.Quartile(data)
	s.q1, s.q2, s.q3 = q.Q1, q.Q2, q.Q3
	summaryArray := [][]string{{"#Count", I2S(s.count)}, {"#Min", F2S(s.min)},
		{"#Max", F2S(s.max)}, {"#Mean", F2S(s.median)}, {"#Median", F2S(s.median)}, {"#SD", F2S(s.sd)},
		{"#Q1", F2S(s.q1)}, {"#Q2", F2S(s.q2)}, {"#Q3", F2S(s.q3)}}
	s.summaryData = summaryArray
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
	//s.count = len(a)
	//summaryArray := [][]string{{"#Count", I2S(s.count)}}
	//s.summaryData = summaryArray

	//catalogue counter
	s.counter = make(map[string]int)
	for _, row := range a {
		s.counter[row]++
	}
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
		s.summaryData = append(s.summaryData, []string{kv.Key, I2S(kv.Value)})
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
