package main

func min(a, b int, c ...int) int {
	minVal := a
	if b < a {
		minVal = b
	}
	for _, val := range c {
		if val < minVal {
			minVal = val
		}
	}
	return minVal
}

func max(a, b int, c ...int) int {
	maxVal := a
	if b > a {
		maxVal = b
	}
	for _, val := range c {
		if val > maxVal {
			maxVal = val
		}
	}
	return maxVal
}