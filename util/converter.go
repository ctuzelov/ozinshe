package util

import "strconv"

func ToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}
