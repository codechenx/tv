package main

import (
	"errors"
)

type sepDetecor struct {
	char []rune
	freq map[rune][]int
}

//There are some features that a separator has.
//Firstly, a separator can split every line into same number of parts(exclude the affection of `"`) and is't zero.(requirement)
//Secondly, delimited files trend to use ',' or '\t' as their separator.
//So if the number of ',' or '\t' can meet requirement . I prefer to use it as a separator.
//If ',' or '\t' can't meet requirement and there is only one char can meet requirement, sepDetect() prefer to use it as a separator.
//Otherwise, it is hard to know which char can be a separator. Machine learning may be a available choice.

func (sd *sepDetecor) sepDetect(s []string) rune {
	preferChar := []rune{',', '	'}
	sd.calFreq(s)

	//sum of int array
	sum := func(ia []int) int {
		var iSum int
		for _, i := range ia {
			iSum += i
		}
		return iSum
	}
	//exclude char appear 0 time in some lines
	var candidate []rune
	for _, char := range sd.char {
		if sum(sd.freq[char]) != 0 {
			candidate = append(candidate, char)
		}
	}

	// if char is in preferChar list and char frequency in every line is also equal, return it
	for _, pc := range preferChar {
		if v := sd.freq[pc]; allIntItemEqual(v) {
			return pc
		}
	}

	//check all char which one's frequency in every line is  equal
	var tempCandidate []rune
	for _, char := range candidate {
		if allIntItemEqual(sd.freq[char]) {
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
	if sd.freq == nil {
		sd.freq = map[rune][]int{}
	}
	charArray := []rune(s[0])
	sd.char = uniqueChar(charArray[:len(charArray)-1])
	for _, char := range sd.char {
		var charFreq []int
		for _, line := range s {
			record, err := lineCSVParse(line, char)
			if err != nil {
				fatalError(errors.New("tv can't identify separator, you need to set it manual"))
			}

			charFreq = append(charFreq, len(record))
		}
		sd.freq[char] = charFreq
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
