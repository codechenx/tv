package main

import (
	"errors"
	"strings"
)

type sepDetecor struct {
	char []rune
	freq []map[rune]int
}

func (sd *sepDetecor) sepDetect(s []string) rune {
	preferChar := []rune{',', '\t'}
	sd.calFreq(s)
	//exclude char appear 0 time in some lines
	var candidate []rune
	charFreqByLine := map[rune][]int{}
	for _, char := range sd.char {
		flag := true
		var freqByLine []int
		for _, freq := range sd.freq {
			if freq[char] == 0 {
				flag = false
			}
			freqByLine = append(freqByLine, freq[char])
		}
		if flag {
			candidate = append(candidate, char)
			charFreqByLine[char] = freqByLine
		}
	}

	// if char is in preferChar list and char frequency in every line is also equal, return it
	for _, pc := range preferChar {
		if v := charFreqByLine[pc]; allIntItemEqual(v) {
			return pc
		}
	}

	//check all char which one's frequency in every line is  equal
	var tempCandidate []rune
	for _, char := range candidate {
		if allIntItemEqual(charFreqByLine[char]) {
			tempCandidate = append(tempCandidate, char)
		}
	}
	candidate = tempCandidate
	//if only one char which one's frequency in every line is  equal, return it
	if len(candidate) == 1 {
		return candidate[0]
	}
	return 0
}

// calculate char(only chars in first line except for last char) frequency in every line
func (sd *sepDetecor) calFreq(s []string) {
	if len(s) < 1 {
		fatalError(errors.New("tv can't identify separator, you need to set it manual"))
	}
	charArray := []rune(s[0])
	sd.char = uniqueChar(charArray[:len(charArray)-1])
	for _, line := range s {
		charFreqSL := map[rune]int{}
		for _, char := range sd.char {
			charFreqSL[char] = strings.Count(line, string(char))
		}
		sd.freq = append(sd.freq, charFreqSL)
	}
}

// remove duplication item in []rune
func uniqueChar(intSlice []rune) []rune {
	keys := make(map[rune]bool)
	var list []rune
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//check if all item in []int is equal, false for empty array
func allIntItemEqual(r []int) bool {
	if len(r) == 0 {
		return false
	}
	flag := true
	for _, i := range r {
		if i != r[0] {
			flag = false
		}
	}
	return flag
}
