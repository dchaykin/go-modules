package log

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	SetLevel(LevelDebug)
	Debug("Debug message")
	Info("Info message: %s", "I am saying you something not important")
	Warn("Warn message: %s", "I am saying you something important")
	Error("Error message: %v", fmt.Errorf("I am saying you something very important"))

	Printfln("Printfln message: %s", "always printed")
	Println("Println message: always printed")
}
