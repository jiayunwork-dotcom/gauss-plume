# gauss-plume：Go 高斯烟羽点源扩散 Web 服务（Pasquill–Gifford + 镜像项 + 前端控制台）

给定源强、有效源高或烟囱参数、风速与稳定度，核算受体浓度与下风向轴线；提供 `/api/conc`、`/api/axis` 与嵌入网页。

## 构建 / 运行 / 测试

```text
go build ./...
./gauss-plume -http :8080
curl -s http://127.0.0.1:8080/api/version
go run . axis example/ground-level.json
go test ./...
```

## 评测镜像

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -d -P --name gauss-plume-b14 <image-name>:latest
curl -s http://127.0.0.1:$(docker port gauss-plume-b14 8080 | cut -d: -f2)/api/version
docker rm -f gauss-plume-b14
```
