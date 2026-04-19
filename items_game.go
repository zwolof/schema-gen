package main

import (
	"fmt"
)

func GetFormattedItemName(title string, count int, length int) string {
	// end result: title            (count) entries
	buf := title

	for len(buf) < length {
		buf += " "
	}

	numbuf := fmt.Sprintf("\033[32m%d\033[0m", count)
	for len(numbuf) < 14 {
		numbuf += " "
	}
	buf += fmt.Sprintf("%s entries", numbuf)

	return buf
}
