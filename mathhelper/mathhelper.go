package mathhelper

// Because there is no Max method for integer in the math package...
func Max(a int, b int) int {
	if (a>b) {
		return a
	} else {
		return b
	}
}
