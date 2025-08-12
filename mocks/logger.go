package mocks

import (
	"github.com/stretchr/testify/mock"
)

type Logger struct {
	mock.Mock
}

func (m *Logger) Info(args ...interface{}) {
	m.Called(args...)
}

func (m *Logger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *Logger) Error(args ...interface{}) {
	m.Called(args...)
}

func (m *Logger) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *Logger) Warn(args ...interface{}) {
	m.Called(args...)
}

func (m *Logger) Warnf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *Logger) Debug(args ...interface{}) {
	m.Called(args...)
}

func (m *Logger) Debugf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *Logger) Fatal(args ...interface{}) {
	m.Called(args...)
}

func (m *Logger) Fatalf(format string, args ...interface{}) {
	m.Called(format, args)
}
