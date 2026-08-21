package geoalert

import (
	"encoding/json"
	"os"
)

// DumpSnapshotJSON 将快照写入文件。
func (e *Engine) DumpSnapshotJSON(path string) error {
	s := e.Snapshot()
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// ImportFencesJSON 从 JSON 数组导入围栏。
func (e *Engine) ImportFencesJSON(path string) (int, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	var list []Fence
	if err := json.Unmarshal(b, &list); err != nil {
		return 0, err
	}
	n := 0
	for _, f := range list {
		if err := e.RegisterFence(f); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}
