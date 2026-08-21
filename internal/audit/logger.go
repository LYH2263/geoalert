package audit

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Logger 追加写审计日志。
type Logger struct {
	mu     sync.Mutex
	f      *os.File
	path   string
	bytes  int64
	maxRot int64
}

func Open(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	st, _ := f.Stat()
	var n int64
	if st != nil {
		n = st.Size()
	}
	return &Logger{f: f, path: path, bytes: n, maxRot: 1 << 20}, nil
}

func (l *Logger) Write(action, id, detail string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f == nil {
		return fmt.Errorf("audit: closed")
	}
	line := fmt.Sprintf("%s\t%s\t%s\t%s\n", time.Now().UTC().Format(time.RFC3339Nano), action, id, detail)
	n, err := l.f.WriteString(line)
	l.bytes += int64(n)
	if err != nil {
		return err
	}
	if l.maxRot > 0 && l.bytes >= l.maxRot {
		return l.rotateLocked()
	}
	return nil
}

func (l *Logger) Flush() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f == nil {
		return nil
	}
	return l.f.Sync()
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f == nil {
		return nil
	}
	err := l.f.Close()
	l.f = nil
	return err
}

func (l *Logger) Path() string { return l.path }

func (l *Logger) SetMaxBytes(n int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.maxRot = n
}

// rotateLocked 轮转：关闭旧文件再打开新文件。
func (l *Logger) rotateLocked() error {
	if l.f == nil {
		return nil
	}
	_ = l.f.Sync()
	if err := l.f.Close(); err != nil {
		return err
	}
	l.f = nil
	rotated := fmt.Sprintf("%s.%d", l.path, time.Now().Unix())
	_ = os.Rename(l.path, rotated)
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	l.f = f
	l.bytes = 0
	return nil
}

// Rotate 强制轮转。
func (l *Logger) Rotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.rotateLocked()
}
