package log

import "testing"

func TestLog(t *testing.T) {
	SetLevel(LevelError)
	Debug("Debug message")
	Info("Info message")
	Warn("Warn message")
	Error("Error message")

	Printfln("Printfln message: %s", "always printed")
	Println("Println message: always printed")
}
