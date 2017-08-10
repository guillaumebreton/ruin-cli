package util

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var red = color.New(color.FgRed).SprintFunc()

func PositiveRed(v float64) string {
	if v > 0 {
		return fmt.Sprintf(red("%0.2f"), v)
	}
	return fmt.Sprintf("%0.2f", v)

}
func NegativeRed(v float64) string {
	if v < 0 {
		return fmt.Sprintf(red("%0.2f"), v)
	}
	return fmt.Sprintf("%0.2f", v)

}

func Format(v float64) string {
	return fmt.Sprintf("%0.2f", v)
}

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
