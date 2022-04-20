package util

import (
	"fmt"

	"github.com/gookit/color"
)

//PrintTabbed Prints a message prepended with a tab to STDOUT
func PrintTabbed(message string) {
	fmt.Printf("	%s\n", message)
}

//PrintSuccess Prints a success message with a green indicator to STDOUT
func PrintSuccess(message string) {
	printIcon(color.FgGreen)
	fmt.Printf(" %s\n", message)
}

//PrintInfo Prints an informational message with a blue indicator to STDOUT
func PrintInfo(message string) {
	printIcon(color.FgBlue)
	fmt.Printf(" %s\n", message)
}

//PrintWarning Prints a warning message with a yellow indicator to STDOUT
func PrintWarning(message string) {
	printIcon(color.FgYellow)
	fmt.Printf(" %s\n", message)
}

//PrintErr Prints an error message with a yellow indicator to STDOUT
func PrintErr(prefix string, err error) {
	printIcon(color.FgRed)
	fmt.Printf(" %s  %s\n", prefix, err)
}

//printIcon Prints a colored square icon to STDOUT
func printIcon(c color.Color) {
	color.Style{c}.Print("■")
}
