package main

import (
	"strings"
	"unicode"
)

type sepDetecor struct {
}

// Fast separator detection algorithm with improved heuristics
// Key improvements:
// 1. Priority-based candidate selection
// 2. Early exit for common separators
// 3. Optimized character counting
// 4. Better validation logic

func (sd *sepDetecor) sepDetect(s []string) rune {
	if len(s) < 1 {
		return 0
	}

	// Fast path: Check common separators first (99% of cases)
	commonSeps := []rune{',', '\t', '|', ';'}
	for _, sep := range commonSeps {
		if sd.isValidSeparator(s, sep) {
			return sep
		}
	}

	// Fallback: Analyze all potential separators
	return sd.detectBestSeparator(s)
}

// Fast validation: Check if a separator is valid for all lines
func (sd *sepDetecor) isValidSeparator(lines []string, sep rune) bool {
	if len(lines) == 0 {
		return false
	}

	// Count separator occurrences in first line
	firstCount := countRuneFast(lines[0], sep)
	if firstCount == 0 {
		return false // Separator not found
	}

	// Verify all lines have same count
	for i := 1; i < len(lines); i++ {
		if countRuneFast(lines[i], sep) != firstCount {
			return false
		}
	}

	return true
}

// Optimized rune counter - much faster than strings.Count for single runes
func countRuneFast(s string, r rune) int {
	count := 0
	for _, c := range s {
		if c == r {
			count++
		}
	}
	return count
}

// Analyze all potential separators when common ones don't work
func (sd *sepDetecor) detectBestSeparator(lines []string) rune {
	if len(lines) == 0 {
		return 0
	}

	// Build candidate list from first line
	candidates := sd.getCandidates(lines[0])
	if len(candidates) == 0 {
		return 0
	}

	// Score each candidate
	type candidateScore struct {
		sep   rune
		score int
		count int
	}

	var scored []candidateScore

	for _, sep := range candidates {
		counts := make([]int, len(lines))
		allEqual := true

		// Count occurrences in each line
		for i, line := range lines {
			counts[i] = countRuneFast(line, sep)
		}

		// Check if all counts are equal and non-zero
		firstCount := counts[0]
		if firstCount == 0 {
			continue
		}

		for i := 1; i < len(counts); i++ {
			if counts[i] != firstCount {
				allEqual = false
				break
			}
		}

		if !allEqual {
			continue
		}

		// Calculate score based on separator quality
		score := sd.scoreSeparator(sep, firstCount)
		scored = append(scored, candidateScore{
			sep:   sep,
			score: score,
			count: firstCount,
		})
	}

	// Return separator with highest score
	if len(scored) == 0 {
		return 0
	}

	best := scored[0]
	for _, candidate := range scored[1:] {
		if candidate.score > best.score {
			best = candidate
		}
	}

	return best.sep
}

// Get candidate separators from first line
func (sd *sepDetecor) getCandidates(line string) []rune {
	// Use map for deduplication
	seen := make(map[rune]bool)
	var candidates []rune

	// Priority characters to check first
	priority := []rune{',', '\t', '|', ';', ':', ' '}
	for _, r := range priority {
		if strings.ContainsRune(line, r) && !seen[r] {
			seen[r] = true
			candidates = append(candidates, r)
		}
	}

	// Check other non-alphanumeric characters
	for _, r := range line {
		if seen[r] || unicode.IsLetter(r) || unicode.IsDigit(r) {
			continue
		}
		// Skip quotes and other problematic chars
		if r == '"' || r == '\'' || r == '\\' {
			continue
		}
		seen[r] = true
		candidates = append(candidates, r)
	}

	return candidates
}

// Score separator quality (higher is better)
func (sd *sepDetecor) scoreSeparator(sep rune, count int) int {
	score := 0

	// Prefer common separators
	switch sep {
	case ',':
		score += 1000 // Highest priority
	case '\t':
		score += 900
	case '|':
		score += 800
	case ';':
		score += 700
	case ':':
		score += 600
	case ' ':
		score += 100 // Lowest priority (can be ambiguous)
	default:
		score += 500 // Moderate priority for other chars
	}

	// Prefer separators with reasonable column counts (2-100)
	if count >= 2 && count <= 100 {
		score += count * 10
	} else if count > 100 {
		score -= 100 // Penalize too many columns
	}

	return score
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

// check if all item in []int is equal, false for empty array
func allIntItemEqual(r []int) bool {
	if len(r) == 0 {
		return false
	}
	for i := 1; i < len(r); i++ {
		if r[i] != r[0] {
			return false
		}
	}
	return true
}
