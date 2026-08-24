package rise

// BranchName 返回本次抬升所走分支的中文描述，供 CLI 与页面展示。
func (r RiseResult) BranchName() string {
	switch {
	case r.Flux <= 0:
		return "无浮力（烟气不高于环境温度，不抬升）"
	case r.Calm:
		return "无风支 ΔH=5.3·Fb^(1/4)"
	case r.Stable:
		return "有风支（稳定级，位温梯度公式）"
	default:
		return "有风支（不稳定/中性，Fb 分段公式）"
	}
}

// FinalRiseKind 返回抬升的短标签：calm / lowflux / highflux / stable / none。
func (r RiseResult) FinalRiseKind() string {
	switch {
	case r.Flux <= 0:
		return "none"
	case r.Calm:
		return "calm"
	case r.Stable:
		return "stable"
	case r.Flux < FbThreshold:
		return "lowflux"
	default:
		return "highflux"
	}
}
