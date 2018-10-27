package log

import (
	kitlog "github.com/go-kit/kit/log"
	kitlevel "github.com/go-kit/kit/log/level"

	"io"
)

const (
	msgKey    = "message"
	moduleKey = "module"
)

type mainLogger struct {
	srcLogger kitlog.Logger
}

// Interface assertions
var _ Logger = (*mainLogger)(nil)

// NewMainLogger returns a logger that encodes msg and keyvals to the Writer
// using go-kit's log as an underlying logger and our custom formatter. Note
// that underlying logger could be swapped with something else.
func NewMainLogger(w io.Writer) Logger {
	srcLogger := kitlog.NewLogfmtLogger(w)

	// return main logger
	return &mainLogger{srcLogger}
}

// Info logs a message at level Info.
func (l *mainLogger) Info(msg string, keyvals ...interface{}) {
	lWithLevel := kitlevel.Info(l.srcLogger)
	if err := kitlog.With(lWithLevel, msgKey, msg).Log(keyvals...); err != nil {
		errLogger := kitlevel.Error(l.srcLogger)
		kitlog.With(errLogger, msgKey, msg).Log("err", err)
	}
}

// Debug logs a message at level Debug.
func (l *mainLogger) Debug(msg string, keyvals ...interface{}) {
	lWithLevel := kitlevel.Debug(l.srcLogger)
	if err := kitlog.With(lWithLevel, msgKey, msg).Log(keyvals...); err != nil {
		errLogger := kitlevel.Error(l.srcLogger)
		kitlog.With(errLogger, msgKey, msg).Log("err", err)
	}
}

// Error logs a message at level Error.
func (l *mainLogger) Error(msg string, keyvals ...interface{}) {
	lWithLevel := kitlevel.Error(l.srcLogger)
	lWithMsg := kitlog.With(lWithLevel, msgKey, msg)
	if err := lWithMsg.Log(keyvals...); err != nil {
		lWithMsg.Log("err", err)
	}
}

// With returns a new contextual logger with keyvals prepended to those passed
// to calls to Info, Debug or Error.
func (l *mainLogger) With(keyvals ...interface{}) Logger {
	return &mainLogger{kitlog.With(l.srcLogger, keyvals...)}
}
