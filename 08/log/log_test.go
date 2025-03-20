package log

import (
	"bytes"
	"context"
	_log "log"
	"strings"
	"testing"
)

////////// Test Helpers //////////

func captureLogOutput() (*bytes.Buffer, func()) {
	var buf bytes.Buffer
	logger := _log.New(&buf, "", 0)
	old_output := _log.Writer()
	_log.SetOutput(logger.Writer())
	return &buf, func() {
		_log.SetOutput(old_output)
	}
}

func assertLogOutput(t *testing.T, buf *bytes.Buffer, expected string) {
	if !strings.Contains(buf.String(), expected) {
		t.Errorf("Expected log output to contain %q, got %q", expected, buf.String())
	}
}

////////// Test Cases //////////

func TestPrintlnNoID(t *testing.T) {
	ctx := context.Background()

	buf, cleanup := captureLogOutput()
	defer cleanup()

	Println(ctx, "Hello, World!")

	expected := "Hello, World!"
	assertLogOutput(t, buf, expected)
}

type idKey struct{}

func TestPrintlnInt(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, idKey{}, 1)

	// Create a buffer to capture log output
	buf, cleanup := captureLogOutput()
	defer cleanup()

	Println(ctx, "Hello, World!")

	expected := "[1] Hello, World!"
	assertLogOutput(t, buf, expected)
}

type MyID struct {
	value int
}

func (m MyID) ID() int {
	return m.value
}

func TestPrintlnID(t *testing.T) {
	my_id := MyID{value: 1}

	ctx := context.Background()
	ctx = context.WithValue(ctx, idKey{}, my_id)

	buf, cleanup := captureLogOutput()
	defer cleanup()

	Println(ctx, "Hello, World!")

	expected := "[1] Hello, World!"
	assertLogOutput(t, buf, expected)
}
