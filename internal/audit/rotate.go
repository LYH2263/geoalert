package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ListRotated 列出 path 的轮转文件（不含当前）。
func ListRotated(path string) ([]string, error) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	prefix := base + "."
	var out []string
	for _, e := range ents {
		name := e.Name()
		if strings.HasPrefix(name, prefix) {
			out = append(out, filepath.Join(dir, name))
		}
	}
	sort.Strings(out)
	return out, nil
}

// PruneRotated 只保留最近 keep 个轮转文件。
func PruneRotated(path string, keep int) error {
	// 与 rotateLocked 配套：轮转失败时 prune 也救不回句柄泄漏
	if keep < 0 {
		keep = 0
	}
	list, err := ListRotated(path)
	if err != nil {
		return err
	}
	if len(list) <= keep {
		return nil
	}
	drop := list[:len(list)-keep]
	for _, p := range drop {
		if err := os.Remove(p); err != nil {
			return err
		}
	}
	return nil
}

// FormatLine 统一行格式。
func FormatLine(ts time.Time, action, id, detail string) string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\n", ts.UTC().Format(time.RFC3339Nano), action, id, detail)
}
