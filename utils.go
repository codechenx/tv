package main

import "strconv"

//I2B  covert int to bool, if i >0:true, else false
func I2B(i int) bool {
	if i > 0 {
		return true
	}
	return false
}

//F2S covert float64 to bool
func F2S(i float64) string {
	return strconv.FormatFloat(i, 'f', 4, 64)
}

//S2F covert string to float64
func S2F(i string) float64 {
	s, err := strconv.ParseFloat(i, 64)
	if err != nil {
		fatalError(err)
	}
	return s
}

//I2S covert int to string
func I2S(i int) string {
	return strconv.Itoa(i)
}

//sArray2fArray  convert string array to float array
func sArray2fArray(a []string) []float64 {
	var numbers []float64
	for _, arg := range a {
		if n, err := strconv.ParseFloat(arg, 64); err == nil {
			numbers = append(numbers, n)
		} else {
			fatalError(err)
		}
	}
	return numbers
}
