package validate

import "math"

// CheckPolygon 校验多边形顶点坐标。
func CheckPolygon(lats, lngs []float64) error {
	if len(lats) < 3 || len(lats) != len(lngs) {
		return ErrTooFewVertices
	}
	for i := range lats {
		if err := LatLng(lats[i], lngs[i]); err != nil {
			return err
		}
		if i > 0 && almostEq(lats[i], lats[i-1]) && almostEq(lngs[i], lngs[i-1]) {
			return ErrDupVertex
		}
	}
	return nil
}

func almostEq(a, b float64) bool {
	return math.Abs(a-b) < 1e-12
}
