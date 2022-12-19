/*
Copyright © 2022 PaoloB <paolo.bellardone@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package utils

import (
	"fmt"

	"github.com/fatih/color"
)

var NoColor bool
var Verbose bool

func PrintVerbose(msg string) {
	if Verbose {
		printMessage(msg, "trace")
	}
}

func PrintInfo(msg string) {
	printMessage(msg, "info")
}

func PrintWarning(msg string) {
	printMessage(msg, "warning")
}

func PrintError(msg string) {
	printMessage(msg, "error")
}

func Print(msg string) {
	printMessage(msg, "")
}

func PrintBold(msg string) {
	b := color.New(color.Bold)
	b.Println(msg)
}

func PrintKV(key string, value string) {
	fmt.Print(key)
	PrintInfo(value)
}

func printMessage(msg string, level string) {
	if NoColor {
		color.NoColor = true
	}

	switch level {
	case "trace":
		fmt.Println(color.HiBlueString(msg))
	case "info":
		fmt.Println(color.HiGreenString(msg))
	case "warning":
		fmt.Println(color.HiYellowString(msg))
	case "error":
		fmt.Println(color.HiRedString(msg))
	default:
		fmt.Println(msg)
	}
}
