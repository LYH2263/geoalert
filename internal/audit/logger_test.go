package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRotateProducesRotatedFile 验证写满触发 Rotate 后会生成轮转文件，
// 且后续写入落到新文件、句柄可用。这是 Windows 下 Rename 文件占用
// 问题的回归测试：旧句柄未关闭时 os.Rename 会因文件被占用而失败。
func TestRotateProducesRotatedFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	l, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = l.Close() })

	// 阈值设为 1：第一次写必然超阈值，触发 write 路径的 rotateLocked。
	l.SetMaxBytes(1)

	// 第一次写：写完后 bytes>=1，触发轮转。当前文件含 id1 的行被重命名
	// 为轮转文件，新当前文件为空。若 Windows 上 Rename 因句柄占用失败，
	// 此处 Write 会返回错误。
	if err := l.Write("action", "id1", "detail"); err != nil {
		t.Fatalf("Write #1 (triggers rotate): %v", err)
	}

	// 关闭自动轮转，保证第二次写不触发轮转、落在新当前文件。
	l.SetMaxBytes(0)
	if err := l.Write("action", "id2", "detail"); err != nil {
		t.Fatalf("Write #2: %v", err)
	}

	// 轮转文件应存在并含首次写入内容。
	rotated, err := ListRotated(path)
	if err != nil {
		t.Fatalf("ListRotated: %v", err)
	}
	if len(rotated) == 0 {
		t.Fatal("expected at least one rotated file, got none")
	}
	got, err := os.ReadFile(rotated[0])
	if err != nil {
		t.Fatalf("ReadFile rotated: %v", err)
	}
	if !strings.Contains(string(got), "id1") {
		t.Errorf("rotated file missing first write content; got %q", got)
	}

	// 当前文件应含第二次写入，且不含第一次（已轮转出去）。
	cur, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile current: %v", err)
	}
	if !strings.Contains(string(cur), "id2") {
		t.Errorf("current file missing post-rotate write; got %q", cur)
	}
	if strings.Contains(string(cur), "id1") {
		t.Errorf("current file unexpectedly contains pre-rotate write; got %q", cur)
	}
}
