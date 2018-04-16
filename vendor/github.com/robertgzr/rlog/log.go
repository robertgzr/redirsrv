package rlog

import (
	apex "github.com/apex/log"
)

var std = Logger{apex.Logger{
	Handler: Default,
	Level:   apex.DebugLevel,
}}

type Logger struct {
	apex.Logger
}

type Fields apex.Fields

func Debug(msg string) {
	std.Debug(msg)
}
func Info(msg string) {
	std.Info(msg)
}
func Warn(msg string) {
	std.Warn(msg)
}
func Error(msg string) {
	std.Error(msg)
}
func Fatal(msg string) {
	std.Fatal(msg)
}

func Trace(msg string) *apex.Entry {
	return std.Trace(msg)
}
func WithError(err error) *apex.Entry {
	return std.WithError(err)
}
func WithField(key string, value interface{}) *apex.Entry {
	return std.WithField(key, value)
}
func WithFields(fields Fields) *apex.Entry {
	return std.WithFields(apex.Fields(fields))
}

// SetHandler overwrites the custom internal handler with whatever is passed as `h`
func SetHandler(h apex.Handler) {
	std.Handler = h
}
