package utils

import (
	"fmt"
	"log"
)

type StdType int

const (
	// Std simple stdout print
	Std StdType = iota
	// StdLog print stdout with date-time
	StdLog
	// StdPanic panic std-out with date-time
	StdPanic
)

func PrintStd(t StdType, key, msg string, args ...interface{}) {
	placeholder := ""

	if key != "" {
		placeholder = fmt.Sprintf("[%s]", key)
	}

	placeholder += msg
	output := fmt.Sprintf(placeholder, args...)

	switch t {
	case Std:
		fmt.Println(output)
	case StdLog:
		log.Println(output)
	case StdPanic:
		log.Fatal(output)
	default:
		log.Printf("[std unknown type]%s\n", output)
	}
}
