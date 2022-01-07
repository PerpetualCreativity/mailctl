package utils

import (
	"os"

	"github.com/fatih/color"
)

// Prints message and err in magenta and then exits
func fancyExit(message string, err interface{}) {
	color.Red("%s: %s", message, err)
	os.Exit(1)
}

func Success(message string) {
	color.Green("Command successfully ran: %s", message)
}

// fancyExits if err != comp
func ErrComp(err interface{}, comp interface{}, message string) {
	if err != comp {
		fancyExit(message, err)
	}
}

// fancyExits if err == comp
func ErrNComp(err interface{}, comp interface{}, message string) {
	if err == comp {
		fancyExit(message, err)
	}
}

// fancyExits if err != nil
func ErrCheck(err interface{}, message string) {
	ErrComp(err, nil, message)
}

// fancyExits if err != nil and err != exp
func ErrExp(err interface{}, exp interface{}, message string) {
	if err != nil && err != exp {
		fancyExit(message, err)
	}
}
