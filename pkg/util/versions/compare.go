package versions

import (
	"sort"
	"strconv"
	"strings"
)

// compare compares two version strings
// returns -1 if v1 < v2, 1 if v1 > v2, 0 otherwise.
func compare(v1, v2 string) int {
	var (
		currTab  = strings.Split(v1, ".")
		otherTab = strings.Split(v2, ".")
	)

	max := len(currTab)
	if len(otherTab) > max {
		max = len(otherTab)
	}
	for i := 0; i < max; i++ {
		var currInt, otherInt int

		if len(currTab) > i {
			currInt, _ = strconv.Atoi(currTab[i])
		}
		if len(otherTab) > i {
			otherInt, _ = strconv.Atoi(otherTab[i])
		}
		if currInt > otherInt {
			return 1
		}
		if otherInt > currInt {
			return -1
		}
	}
	return 0
}

// LessThan checks if a version is less than another
func LessThan(v, other string) bool {
	return compare(v, other) == -1
}

// SortEngines orders engine names such as docker-29.10.0 or k3s-v1.34.1-k3s1
// by type and then by version, so that 29.10 comes after 29.9 (a plain string
// sort puts it before). The version is everything after the first "-" and is
// compared as the sequence of numbers it contains, so v1.34.1-k3s1 < v1.34.1-k3s2.
func SortEngines(engines []string) {
	sort.SliceStable(engines, func(i, j int) bool {
		typeI, versionI := splitEngine(engines[i])
		typeJ, versionJ := splitEngine(engines[j])
		if typeI != typeJ {
			return typeI < typeJ
		}
		return compareNumbers(versionI, versionJ) < 0
	})
}

func splitEngine(engine string) (string, string) {
	if i := strings.Index(engine, "-"); i > 0 {
		return engine[:i], engine[i+1:]
	}
	return engine, ""
}

// compareNumbers compares the sequences of decimal numbers in two strings:
// "v1.34.1-k3s1" is [1 34 1 3 1]. A missing number counts as zero.
func compareNumbers(v1, v2 string) int {
	n1, n2 := numbers(v1), numbers(v2)
	for i := 0; i < len(n1) || i < len(n2); i++ {
		var a, b int
		if i < len(n1) {
			a = n1[i]
		}
		if i < len(n2) {
			b = n2[i]
		}
		if a != b {
			if a < b {
				return -1
			}
			return 1
		}
	}
	return 0
}

func numbers(s string) []int {
	result := []int{}
	current := -1
	for _, r := range s {
		if r >= '0' && r <= '9' {
			if current < 0 {
				current = 0
			}
			current = current*10 + int(r-'0')
			continue
		}
		if current >= 0 {
			result = append(result, current)
			current = -1
		}
	}
	if current >= 0 {
		result = append(result, current)
	}
	return result
}
