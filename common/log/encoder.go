package log

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"sync"
	"time"
)

type PlainEncoder struct {
	writer       io.Writer
	mu           sync.Mutex
	EnableBuffer bool
}

func (en *PlainEncoder) setWriter(writer io.Writer) func() error {
	if en.EnableBuffer {
		bufWriter := bufio.NewWriter(writer)
		en.writer = bufWriter
		return func() error {
			return bufWriter.Flush()
		}
	} else {
		en.writer = writer
		return func() error { return nil }
	}
}

func (en *PlainEncoder) write(level Level, path string, trace string, msg string, kvs []interface{}) {
	now := time.Now().Format("2006-01-02 15:04:05")
	buf := &bytes.Buffer{}
	buf.WriteByte('[')
	buf.WriteString(level.CapitalString())
	buf.WriteByte(']')
	buf.WriteByte(' ')
	buf.WriteString(now)
	if len(path) > 0 {
		buf.WriteByte(' ')
		buf.WriteString(path)
	}
	buf.WriteByte(' ')
	if len(trace) == 0 {
		buf.WriteString(msg)
	} else {
		buf.WriteString(trace)
	}
	buf.WriteByte('\n')
	en.mu.Lock()
	defer en.mu.Unlock()
	en.writer.Write(buf.Bytes())
}

type JsonEncoder struct {
	writer       io.Writer
	mu           sync.Mutex
	EnableBuffer bool
}

func (en *JsonEncoder) setWriter(writer io.Writer) func() error {
	if en.EnableBuffer {
		bufWriter := bufio.NewWriter(writer)
		en.writer = bufWriter
		return func() error {
			return bufWriter.Flush()
		}
	} else {
		en.writer = writer
		return func() error { return nil }
	}
}

func (en *JsonEncoder) write(level Level, path string, trace string, msg string, kvs []interface{}) {
	params := make(map[string]interface{}, len(kvs)/2+5)
	for i := 0; i < len(kvs); i += 2 {
		if k, ok := kvs[i].(string); ok {
			params[k] = kvs[i+1]
		}
	}
	params["level"] = level.CapitalString()
	params["datetime"] = time.Now().Format("2006-01-02 15:04:05")
	params["path"] = path
	params["msg"] = msg
	params["trace"] = trace
	value, err := json.Marshal(params)
	if err != nil {
		return
	}
	value = append(value, '\n')
	en.mu.Lock()
	defer en.mu.Unlock()
	en.writer.Write(value)
}
