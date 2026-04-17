// Package output handles JSON / text / jq / template rendering uniformly.
package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

type Format int

const (
	FormatText Format = iota
	FormatJSON
)

type Writer struct {
	Out      io.Writer
	Err      io.Writer
	Format   Format
	Timing   bool
	JQ       string
	Template string

	start time.Time
}

type Options struct {
	JSON     bool
	JQ       string
	Template string
	Timing   bool
}

func New(opts Options) *Writer {
	w := &Writer{
		Out:      os.Stdout,
		Err:      os.Stderr,
		Timing:   opts.Timing,
		JQ:       opts.JQ,
		Template: opts.Template,
		start:    time.Now(),
	}
	if opts.JSON || opts.JQ != "" || opts.Template != "" {
		w.Format = FormatJSON
	}
	return w
}

func (w *Writer) IsJSON() bool { return w.Format == FormatJSON }

func (w *Writer) JSON(v interface{}) error {
	if w.JQ != "" {
		return w.pipeJQ(v)
	}
	if w.Template != "" {
		return w.pipeTemplate(v)
	}
	enc := json.NewEncoder(w.Out)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func (w *Writer) Text(format string, args ...interface{}) {
	fmt.Fprintf(w.Out, format, args...)
}

func (w *Writer) Warn(format string, args ...interface{}) {
	fmt.Fprintf(w.Err, format, args...)
}

func (w *Writer) Finish() {
	if w.Timing {
		elapsed := time.Since(w.start)
		fmt.Fprintf(w.Err, "%.1fms\n", float64(elapsed.Microseconds())/1000)
	}
}

func (w *Writer) pipeJQ(v interface{}) error {
	payload, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal for --jq: %w", err)
	}
	bin, err := exec.LookPath("jq")
	if err != nil {
		return fmt.Errorf("--jq requires jq on PATH: %w", err)
	}
	cmd := exec.Command(bin, w.JQ)
	cmd.Stdin = bytes.NewReader(payload)
	cmd.Stdout = w.Out
	cmd.Stderr = w.Err
	return cmd.Run()
}

func (w *Writer) pipeTemplate(v interface{}) error {
	t, err := template.New("out").Parse(w.Template)
	if err != nil {
		return fmt.Errorf("parse --template: %w", err)
	}
	payload, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal for --template: %w", err)
	}
	var generic interface{}
	if err := json.Unmarshal(payload, &generic); err != nil {
		return fmt.Errorf("unmarshal for --template: %w", err)
	}
	return t.Execute(w.Out, generic)
}

func WordWrap(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}
	var lines []string
	words := strings.Fields(text)
	current := ""
	for _, word := range words {
		switch {
		case current == "":
			current = word
		case len(current)+1+len(word) <= width:
			current += " " + word
		default:
			lines = append(lines, current)
			current = word
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}
