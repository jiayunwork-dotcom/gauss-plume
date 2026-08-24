# gauss-plume：高斯烟羽点源扩散核算工具

gauss-plume 读入源强 Q、源描述（有效源高或烟囱+烟气参数）、风速、稳定度 A–F 与下风向距离，按 Pasquill–Gifford 幂律计算 σy/σz，用含地面镜像反射项的高斯烟羽公式给出受体浓度，并可选叠加 Briggs 热抬升。

## 构建 / 运行 / 测试

```text
go build ./...
go run . -http :8080
go run . axis example/ground-level.json
go test ./...
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
