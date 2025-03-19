package log

import (
	"fmt"
	"log/slog"
	"os"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

func SetLevel(level int) {
	l := slog.LevelInfo
	switch level {
	case LevelDebug:
		l = slog.LevelDebug
	case LevelInfo:
		l = slog.LevelInfo
	case LevelWarn:
		l = slog.LevelWarn
	case LevelError:
		l = slog.LevelError
	}
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l}))
}

func Println(text ...any) {
	fmt.Println(text...)
}

func Printfln(msg string, text ...any) {
	fmt.Printf(msg+"\n", text...)
}

func Debug(msg string, text ...any) {
	logger.Debug(msg, text...)
}

func Warn(msg string, text ...any) {
	logger.Warn(msg, text...)
}

func Info(msg string, text ...any) {
	logger.Info(msg, text...)
}

func Error(msg string, text ...any) {
	logger.Error(msg, text...)
}
