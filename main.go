// Command gauss-plume 是高斯烟羽点源扩散核算工具。
//
// 用户给出源强 Q、源（直接有效源高，或烟囱+烟气参数由 Briggs 抬升）、
// 风速、稳定度 A–F 与下风向距离，按 Pasquill–Gifford 幂律计算 σy/σz，
// 用含地面镜像反射项的高斯烟羽公式给出受体浓度。非法输入一律报错，
// 绝不静默给数值。
//
// 子命令：
//
//	gauss-plume -http :8080      启动 Web 控制台（默认 :8080）
//	gauss-plume serve [-http]    同上
//	gauss-plume axis <算例.json>  打印下风向轴线浓度点列
//	gauss-plume conc <算例.json>  打印单受体浓度
//	gauss-plume check <算例.json> 交叉规则自检（Q 加倍 / 风速加倍 / 稳定度 / 远场衰减）
//	gauss-plume peak <算例.json>  定位地面轴线峰值位置与浓度
//	gauss-plume sigma [A–F]       打印 σ 参数表或单等级明细
//	gauss-plume version          打印版本
//	gauss-plume help             显示帮助
package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"gauss-plume/internal/dispersion"
	"gauss-plume/internal/plume"
	"gauss-plume/internal/server"
)

//go:embed web example
var assets embed.FS

// version 是程序版本号。
const version = "1.0.0"

// usageText 是帮助文本。
const usageText = `gauss-plume —— 高斯烟羽点源扩散核算

用法:
  gauss-plume -http :8080       启动 Web 控制台（默认 :8080）
  gauss-plume serve [-http :9090]
  gauss-plume axis <算例.json>   打印下风向轴线浓度点列
  gauss-plume conc <算例.json>   打印单受体浓度
  gauss-plume check <算例.json>  交叉规则自检
  gauss-plume peak <算例.json>   定位地面轴线峰值位置与浓度
  gauss-plume sigma [A–F]        打印 σ 参数表或单等级明细
  gauss-plume version / help

算例示例:
  gauss-plume axis example/ground-level.json
  gauss-plume conc example/point.json
  gauss-plume check example/check.json

算例 JSON 字段（POST /api/axis 与 CLI 通用）:
  q            源强（g/s，必须 ≥ 0）
  source       有效源高 {"height": 60} 或烟囱参数 {"stack": {...}}
  wind_speed   风速（m/s，必须 > 0）
  stability    稳定度 A–F
  axis         {"start":50,"end":2000,"count":40,"y":0,"z":0}
  单点另有 receptor {"x":500,"y":0,"z":0}

非法输入（风速 ≤ 0、源强 < 0、距离 ≤ 0、稳定度非 A–F 等）一律报错：
HTTP 返回 {"error": "..."}，CLI 写 stderr 并以非零退出码结束。
`

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fail("缺少子命令，运行 gauss-plume help 查看用法")
	}
	switch args[0] {
	case "serve":
		runServe(args[1:])
	case "axis":
		need(args, "axis 需要一个算例文件参数")
		runAxis(args[1])
	case "conc":
		need(args, "conc 需要一个算例文件参数")
		runConc(args[1])
	case "check":
		need(args, "check 需要一个算例文件参数")
		runCheck(args[1])
	case "peak":
		need(args, "peak 需要一个算例文件参数")
		runPeak(args[1])
	case "sigma":
		runSigma(args[1:])
	case "version", "-v", "--version":
		fmt.Printf("gauss-plume %s\n", version)
	case "help", "-h", "--help":
		fmt.Print(usageText)
	default:
		if strings.HasPrefix(args[0], "-") {
			runServe(args)
			return
		}
		fail("未知子命令 %q，运行 gauss-plume help 查看用法", args[0])
	}
}

// need 检查子命令的参数个数。
func need(args []string, msg string) {
	if len(args) < 2 {
		fail("%s", msg)
	}
}

// runServe 启动 Web 控制台。
func runServe(args []string) {
	fsets := flag.NewFlagSet("serve", flag.ContinueOnError)
	fsets.SetOutput(os.Stderr)
	addr := fsets.String("http", ":8080", "HTTP 监听地址")
	if err := fsets.Parse(args); err != nil {
		os.Exit(2)
	}
	webFS, err := fs.Sub(assets, "web")
	if err != nil {
		fail("web 资源缺失：%v", err)
	}
	exampleFS, err := fs.Sub(assets, "example")
	if err != nil {
		fail("example 资源缺失：%v", err)
	}
	srv := server.New(webFS, exampleFS, version)
	fmt.Printf("gauss-plume %s 已启动，打开 http://localhost%s\n", version, *addr)
	if err := http.ListenAndServe(*addr, srv.Handler()); err != nil {
		fail("HTTP 服务退出：%v", err)
	}
}

// runAxis 加载轴线算例并打印结果。
func runAxis(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fail("%v", err)
	}
	req, err := plume.AxisRequestFromJSON(data)
	if err != nil {
		fail("算例 JSON 解析失败：%v", err)
	}
	resp, err := plume.ComputeAxis(req)
	if err != nil {
		fail("%v", err)
	}
	fmt.Print(resp.String())
}

// runConc 加载单点算例并打印结果。
func runConc(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fail("%v", err)
	}
	req, err := plume.PointRequestFromJSON(data)
	if err != nil {
		fail("算例 JSON 解析失败：%v", err)
	}
	resp, err := plume.ComputePoint(req)
	if err != nil {
		fail("%v", err)
	}
	fmt.Print(resp.String())
}

// runCheck 加载基例并运行交叉规则自检。
func runCheck(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fail("%v", err)
	}
	var c plume.CheckCase
	if err := jsonUnmarshalStrict(data, &c); err != nil {
		fail("算例 JSON 解析失败：%v", err)
	}
	rep, err := plume.RunChecks(c)
	if err != nil {
		fail("%v", err)
	}
	fmt.Print(rep.String())
}

// runPeak 加载轴线算例并定位峰值。
func runPeak(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fail("%v", err)
	}
	req, err := plume.AxisRequestFromJSON(data)
	if err != nil {
		fail("算例 JSON 解析失败：%v", err)
	}
	st, err := plume.ValidateAxisRequest(req)
	if err != nil {
		fail("%v", err)
	}
	state, err := plume.ResolveSource(req.Source, st, req.WindSpeed)
	if err != nil {
		fail("%v", err)
	}
	peak, err := plume.PeakOnAxis(req.Q, req.WindSpeed, state.Height, st, req.Axis.Start, req.Axis.End)
	if err != nil {
		fail("%v", err)
	}
	fmt.Printf("峰值定位（%s，H=%.4g m）：x* = %.4g m，Cmax = %.6g g/m³\n",
		st, state.Height, peak.Distance, peak.Concentration)
}

// runSigma 打印扩散参数表。
func runSigma(args []string) {
	if len(args) == 0 {
		xs, err := dispersion.Grid(100, 10000, 6)
		if err != nil {
			fail("%v", err)
		}
		tab, err := dispersion.Tabulate(xs)
		if err != nil {
			fail("%v", err)
		}
		fmt.Print(dispersion.FormatTable(tab))
		return
	}
	st, err := dispersion.ParseStability(args[0])
	if err != nil {
		fail("%v", err)
	}
	line, err := dispersion.DescribeClass(st)
	if err != nil {
		fail("%v", err)
	}
	fmt.Println(line)
	xs, err := dispersion.Grid(100, 10000, 6)
	if err != nil {
		fail("%v", err)
	}
	for _, x := range xs {
		sg, err := dispersion.Dispersion(st, x)
		if err != nil {
			fail("%v", err)
		}
		fmt.Printf("  x=%8.0f m: %s\n", x, sg.String())
	}
}

// jsonUnmarshalStrict 用严格模式解析算例 JSON（拒绝未知字段）。
func jsonUnmarshalStrict(data []byte, dst any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

// fail 把错误写入 stderr 并以退出码 1 结束。
func fail(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "gauss-plume: "+format+"\n", a...)
	os.Exit(1)
}
