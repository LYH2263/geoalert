// Package geoalert 提供地理围栏告警引擎：注册多边形/圆形围栏、
// 摄入对象位置序列、检测 enter/exit 边沿与驻留（dwell），并发出 Alert。
//
// 典型用法：
//
//	eng, _ := geoalert.New(geoalert.Options{Clock: clock.Real{}})
//	_ = eng.RegisterFence(geoalert.Fence{ID: "yard", Kind: geoalert.FencePolygon, Vertices: ring})
//	_ = eng.Ingest(geoalert.TrackPoint{ObjectID: "truck-1", At: t, Lat: 31.2, Lng: 121.5})
//	alerts := eng.Alerts()
package geoalert
