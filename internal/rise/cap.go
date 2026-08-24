package rise

// capHeight 应用逆温顶截断：给定烟囱高度 hs 与抬升 rise，
// 若逆温层底高 zInv > 0 且 hs+rise 超过 zInv，则有效源高被截断到 zInv。
func capHeight(hs, rise, zInv float64) (eff float64, capped bool) {
	eff = hs + rise
	if zInv > 0 && eff > zInv {
		return zInv, true
	}
	return eff, false
}
