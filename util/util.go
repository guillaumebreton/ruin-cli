package util

import (
	"fmt"
	"os"
)

func HandleError(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, msg)
	}
}
func Exit(msg string) {
	fmt.Fprintf(os.Stderr, msg)
	os.Exit(1)
}
func ExitOnError(err error, msg string) {
	if err != nil {
		Exit(msg)
	}
}
