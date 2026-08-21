package geom

import "math"

// GridIndex 粗粒度经纬度网格，用于围栏候选过滤。
type GridIndex struct {
	cellDeg float64
	cells   map[cellKey][]string
}

type cellKey struct {
	i, j int
}

func NewGridIndex(cellDeg float64) *GridIndex {
	if cellDeg <= 0 {
		cellDeg = 0.01
	}
	return &GridIndex{cellDeg: cellDeg, cells: make(map[cellKey][]string)}
}

func (g *GridIndex) key(lat, lng float64) cellKey {
	return cellKey{
		i: int(math.Floor(lat / g.cellDeg)),
		j: int(math.Floor(lng / g.cellDeg)),
	}
}

// InsertBBox 将 ID 插入包围盒覆盖的格子。
func (g *GridIndex) InsertBBox(id string, b BBox) {
	i0 := int(math.Floor(b.MinLat / g.cellDeg))
	i1 := int(math.Floor(b.MaxLat / g.cellDeg))
	j0 := int(math.Floor(b.MinLng / g.cellDeg))
	j1 := int(math.Floor(b.MaxLng / g.cellDeg))
	for i := i0; i <= i1; i++ {
		for j := j0; j <= j1; j++ {
			k := cellKey{i: i, j: j}
			g.cells[k] = appendUnique(g.cells[k], id)
		}
	}
}

func appendUnique(xs []string, id string) []string {
	for _, x := range xs {
		if x == id {
			return xs
		}
	}
	return append(xs, id)
}

// Query 查询点所在格子的候选 ID。
func (g *GridIndex) Query(lat, lng float64) []string {
	return append([]string(nil), g.cells[g.key(lat, lng)]...)
}

// Remove 删除 ID（全表扫描）。
func (g *GridIndex) Remove(id string) {
	for k, xs := range g.cells {
		out := xs[:0]
		for _, x := range xs {
			if x != id {
				out = append(out, x)
			}
		}
		if len(out) == 0 {
			delete(g.cells, k)
		} else {
			g.cells[k] = out
		}
	}
}
