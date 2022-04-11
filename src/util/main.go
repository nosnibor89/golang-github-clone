package util

import (
	"fmt"
)

func StringIsEmpty(text string) bool {
	return len(text) == 0
}

var reset = "\033[0m"
var red = "\033[31m"
var cyan = "\033[36m"

func PrintRed(data interface{}) {
	fmt.Printf("%s %v %s\n", red, data, reset)
}

func PrintCyan(data interface{}) {
	fmt.Printf("%s %v %s\n", cyan, data, reset)
}
