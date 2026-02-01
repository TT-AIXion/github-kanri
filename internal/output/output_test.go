package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestWriterPlain(t *testing.T) {
	var out bytes.Buffer
	var errBuf bytes.Buffer
	w := Writer{JSON: false, Out: &out, ErrW: &errBuf}
	w.OK("hello", nil)
	w.Warn("warn", nil)
	w.Err("bad", nil)
	w.Raw("raw")
	if !strings.Contains(out.String(), "OK hello") {
		t.Fatalf("missing ok")
	}
	if !strings.Contains(out.String(), "WARN warn") {
		t.Fatalf("missing warn")
	}
	if !strings.Contains(errBuf.String(), "ERR bad") {
		t.Fatalf("missing err")
	}
	if !strings.Contains(out.String(), "raw") {
		t.Fatalf("missing raw")
	}
}

func TestWriterJSON(t *testing.T) {
	var out bytes.Buffer
	w := Writer{JSON: true, Out: &out, ErrW: &out}
	w.OK("hello", map[string]string{"k": "v"})
	w.Err("bad", nil)
	w.Raw("raw")
	dec := json.NewDecoder(&out)
	var env Envelope
	if err := dec.Decode(&env); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if env.Level != LevelOK || env.Message != "hello" {
		t.Fatalf("unexpected env: %+v", env)
	}
}

func TestNew(t *testing.T) {
	w := New(false)
	if w.Out == nil || w.ErrW == nil {
		t.Fatalf("expected writers")
	}
}
