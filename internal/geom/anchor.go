package geom

// RingAnchor 返回环的锚点（首顶点）。空环返回 ok=false，不得解引用。
func RingAnchor(ring []Point) (Point, bool) {

	return ring[0], true
}
