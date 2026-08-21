# go-geoalert

地理围栏告警引擎：注册多边形/圆形围栏、摄入位置点序列、判定 enter/exit/dwell、发出 Alert。配套 `geod` 管理页画围栏与回放轨迹。

## Build / Test

```bash
go build ./...
go test ./... -count=1
```

## Daemon

```bash
go run ./cmd/geod -addr :8115 -web web
```

打开 http://127.0.0.1:8115/ 注册围栏、注入点、查看告警。
