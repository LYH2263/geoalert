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

// rotateLocked 轮转：先 Sync/Close 旧句柄再 Rename，最后打开新文件。
// Windows 下若不先关闭占用句柄，Rename 会因文件被占用而失败，
// 导致轮转文件出不来。
func (l *Logger) rotateLocked() error {
	if l.f == nil {
		return nil
	}

	// 先刷盘再关闭旧句柄，释放对 path 的占用，
	// 这样 Windows 上的 Rename 才不会报文件占用。
	_ = l.f.Sync()
	old := l.f
	l.f = nil
	if err := old.Close(); err != nil {
		return err
	}

	rotated := fmt.Sprintf("%s.%d", l.path, time.Now().Unix())
	if err := os.Rename(l.path, rotated); err != nil {
		// Rename 失败时回退：重新打开原文件，保持 Logger 可用。
		if f, ferr := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644); ferr == nil {
			l.f = f
		}
		return err
	}

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
