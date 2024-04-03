package crawler

import "github.com/go-rod/rod"

type Pool interface {
	Get() *rod.Page
	Put(p *rod.Page)
}

type RodPool struct {
	pool   *rod.PagePool
	create func() *rod.Page
}

func NewRodPool(pool *rod.PagePool,
	create func() *rod.Page) Pool {
	return &RodPool{pool, create}
}

func (p *RodPool) Get() *rod.Page {
	return p.pool.Get(p.create)
}

func (p *RodPool) Put(page *rod.Page) {
	p.pool.Put(page)
}
