package log

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
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

func Printfln(msg string, args ...any) {
	fmt.Printf(msg+"\n", args...)
}

func Debug(msg string, args ...any) {
	logger.Debug(fmt.Sprintf(msg, args...))
}

func Warn(msg string, args ...any) {
	logger.Warn(fmt.Sprintf(msg, args...))
}

func Info(msg string, args ...any) {
	logger.Info(fmt.Sprintf(msg, args...))
}

func Error(msg string, args ...any) {
	buf := make([]byte, 1024)
	runtime.Stack(buf, false)
	logger.Error(fmt.Sprintf(msg, args...), "Stacktrace", string(buf))
}
