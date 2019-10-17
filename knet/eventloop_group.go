package knet

type eventLoopGroup struct {
	nextLoopIndex int
	eventLoops    []*loop
	size          int
}

type IEventLoopGroup interface {
	register(*loop)
	next() *loop
	len() int
	iterate(func(int, *loop) bool)
}

func (g *eventLoopGroup) register(l *loop) {
	g.eventLoops = append(g.eventLoops, l)
	g.size++
}

func (g *eventLoopGroup) next() *loop {
	lp := g.eventLoops[g.nextLoopIndex]
	g.nextLoopIndex++
	if g.nextLoopIndex >= g.size {
		g.nextLoopIndex = 0
	}

	return lp
}

func (g *eventLoopGroup) iterate(f func(int, *loop) bool) {
	for i, lp := range g.eventLoops {
		if !f(i, lp) {
			break
		}
	}
}

func (g *eventLoopGroup) len() int {
	return g.size
}
