# gauss-plume

gauss-plume 是高斯烟羽点源扩散 Web 核算服务。用户输入源强 Q（g/s）、源描述（直接给有效源高，或给烟囱高度、出口速度、出口内径与烟气/环境温度，由 Briggs 公式计算热抬升）、风速 u、Pasquill 稳定度（A–F）与下风向距离，后端按 Pasquill–Gifford 幂律给出 σy/σz，并用含地面镜像反射项的高斯烟羽公式计算受体浓度。`-http :8080` 提供 `/api/conc` 与 `/api/axis`。面向单点源地面/受体浓度与下风向轴线核算。

## 用法

启动 Web 控制台（默认监听 `:8080`，前端页面由 Go 同进程提供）：

```text
go run . -http :8080
```

打开 `http://localhost:8080`：输入算例或点击「加载示例（地面源）」，即可列出受体浓度并画出下风向曲线；非法参数会在页面和 API 错误体中显示后端返回的说明。

离线 CLI：

```text
go run . axis example/ground-level.json
go run . conc example/point.json
go run . axis example/stack.json
go run . check example/check.json
```

`check` 对给定基例执行交叉规则自检：Q 加倍 ⇒ 处处浓度加倍；u 加倍 ⇒ 轴线浓度明显下降（地面源约减半）；不稳定→稳定 ⇒ 同一距离 σz 变小且轴线形态改变；远场衰减指数符合预期。

## API

| 方法 | 路径 | 输入 | 返回 |
|------|------|------|------|
| POST | `/api/conc` | 源 + 受体 | 浓度、σy、σz、有效源高 |
| POST | `/api/axis` | 源 + 下风向网格 | 轴线浓度点列 |
| GET | `/api/version` | — | 服务名与版本 |

请求示例：

```json
{
  "q": 5,
  "source": { "height": 0 },
  "wind_speed": 3,
  "stability": "D",
  "receptor": { "x": 500, "y": 0, "z": 0 }
}
```

烟囱模式把 `source` 换成：

```json
{ "stack": {
    "height": 40, "exit_velocity": 12, "radius": 1.2,
    "gas_temperature": 420, "ambient_temperature": 288,
    "inversion_top": 0, "gradient": 0 } }
```

任何非法输入（风速 ≤ 0、源强 < 0、距离 ≤ 0、稳定度非 A–F、参数缺失或非有限数值）统一返回 HTTP 400 与错误体 `{"error": "..."}`；CLI 则写 stderr 并以非零退出码结束，绝不静默给数值。

## 关键约定

- 浓度公式钉死 2π 形式：

  ```text
  C = Q/(2π·u·σy·σz) · exp(−y²/2σy²) · [exp(−(z−H)²/2σz²) + exp(−(z+H)²/2σz²)]
  ```

  地面受体 z=0 时两个镜像项相等，合并为 π 形式；H=0 的地面源轴线浓度即 `Q/(π·u·σy·σz)`。
- σ 为 Pasquill–Gifford 农村幂律：`σy = a·x^b`、`σz = c·x^d`，x 与 σ 均以米计，指数分段在 0.80–1.00。
- Briggs 热抬升按风速分段：u < 1 m/s 走无风支 `ΔH = 5.3·Fb^(1/4)`；有风支不稳定/中性（A–D）按 `Fb < 55` 与 `Fb ≥ 55` 两段，稳定级（E/F）走位温梯度公式。有效源高 `H = hs + ΔH`，可选逆温层底高截断。
- 烟气温度不高于环境温度时如实报告不抬升，不把烟囱几何高度冒充已抬升的有效源高。
- 下风向距离上限 1000 km；轴线网格点数上限 10000。

## 构建与测试

```text
go build ./...
go test ./...
go run . -http :8080
```

## 目录

- `internal/dispersion/` — 稳定度解析、Pasquill–Gifford 系数表、σ 评估与一致性校验
- `internal/rise/` — Briggs 热抬升（Fb 通量、无风/有风两支、逆温顶截断）
- `internal/plume/` — 高斯烟羽浓度公式、源解析、单点/轴线核算与交叉规则自检
- `internal/server/` — 薄 HTTP 层与静态前端托管
- `example/` — 离线算例
- `web/` — 交互页面（仅渲染后端返回的数值）
