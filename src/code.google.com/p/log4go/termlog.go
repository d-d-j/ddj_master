// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"fmt"
	"io"
	"os"
)

const (
	RED    string = "\033[31m"
	GREEN  string = "\033[32m"
	YELLOW string = "\033[33m"
	BLUE   string = "\033[34m"
	GRAY   string = "\033[90m"
	RESET  string = "\033[0m"
)

var (
	levelColor = [...]string{GRAY, GRAY, RESET, RESET, BLUE, YELLOW, RED, RED}
)

var stdout io.Writer = os.Stdout

// This is the standard writer that prints to standard output.
type ConsoleLogWriter chan *LogRecord

// This creates a new ConsoleLogWriter
func NewConsoleLogWriter() ConsoleLogWriter {
	records := make(ConsoleLogWriter, LogBufferLength)
	go records.run(stdout)
	return records
}

func (w ConsoleLogWriter) run(out io.Writer) {
	var timestr string
	var timestrAt int64

	for rec := range w {
		if at := rec.Created.UnixNano() / 1e9; at != timestrAt {
			timestr, timestrAt = rec.Created.Format("15:04:05"), at
		}
		fmt.Fprint(out, "[", GRAY, timestr, RESET, "] [", levelColor[rec.Level], levelStrings[rec.Level], RESET, "] (", GREEN, rec.Source, RESET, ") ", levelColor[rec.Level], rec.Message, RESET, "\n")
	}
}

// This is the ConsoleLogWriter's output method.  This will block if the output
// buffer is full.
func (w ConsoleLogWriter) LogWrite(rec *LogRecord) {
	w <- rec
}

// Close stops the logger from sending messages to standard output.  Attempts to
// send log messages to this logger after a Close have undefined behavior.
func (w ConsoleLogWriter) Close() {
	close(w)
}
