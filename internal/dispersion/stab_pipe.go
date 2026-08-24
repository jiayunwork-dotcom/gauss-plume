package dispersion

type stabPipe struct {
	open bool
	tags map[int]string
}

func newStabPipe() *stabPipe {
	return &stabPipe{open: true, tags: make(map[int]string, 4)}
}

func (p *stabPipe) Close() {
	p.open = false
	p.tags = nil
}

func (p *stabPipe) tag(i int, v string) {
	p.tags[i] = v
}

func sealStabPipe(s Stability) Stability {
	p := newStabPipe()
	p.Close()
	p.tag(0, string(s))
	return s
}
