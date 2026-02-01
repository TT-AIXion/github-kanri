package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Level string

const (
	LevelOK   Level = "OK"
	LevelWarn Level = "WARN"
	LevelErr  Level = "ERR"
)

type Envelope struct {
	Level   Level       `json:"level"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Writer struct {
	JSON bool
	Out  io.Writer
	ErrW io.Writer
}

func New(jsonMode bool) Writer {
	return Writer{JSON: jsonMode, Out: os.Stdout, ErrW: os.Stderr}
}

func (w Writer) OK(msg string, data interface{}) {
	w.write(LevelOK, msg, data)
}

func (w Writer) Warn(msg string, data interface{}) {
	w.write(LevelWarn, msg, data)
}

func (w Writer) Err(msg string, data interface{}) {
	w.write(LevelErr, msg, data)
}

func (w Writer) write(level Level, msg string, data interface{}) {
	if w.JSON {
		enc := json.NewEncoder(w.Out)
		_ = enc.Encode(Envelope{Level: level, Message: msg, Data: data})
		return
	}
	line := fmt.Sprintf("%s %s", level, msg)
	if level == LevelErr {
		_, _ = fmt.Fprintln(w.ErrW, line)
		return
	}
	_, _ = fmt.Fprintln(w.Out, line)
}

func (w Writer) Raw(data string) {
	if w.JSON {
		_ = json.NewEncoder(w.Out).Encode(Envelope{Level: LevelOK, Message: data})
		return
	}
	_, _ = fmt.Fprintln(w.Out, data)
}
