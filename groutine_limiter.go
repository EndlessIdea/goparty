package goparty

import "errors"

type GLimiter interface {
	Get() bool
	Put() bool
}

type myGLimiter struct {
	total 	int32
	pChan	chan struct{}
}

func NewGLimiter(count int32) (gl GLimiter, err error) {
	p := &myGLimiter{total: count}
	p.pChan = make(chan struct{}, count)
	if count <= 0 {
		return nil, errors.New("count param must be bigger than zero")
	}
	n := int(count)
	for i := 0; i < n; i++  {
		p.pChan <- struct{}{}
	}

	return p, err
}

func (p *myGLimiter) Get() bool {
	select {
	case <-p.pChan:
		return true
	default:
		return false
	}
}

func (p *myGLimiter) Put() bool {
	c := struct {}{}
	select {
	case p.pChan <- c:
		return true
	default:
		return false
	}
}
