package pool

type glimit struct {
	num int
	c   chan struct{}
}

// 协程限制器
func NewGlimit(num int) *glimit {
	return &glimit{
		num: num,
		c:   make(chan struct{}, num),
	}
}

func (g *glimit) Run(f func()) {
	g.c <- struct{}{}
	go func() {
		f()
		<-g.c
	}()
}
