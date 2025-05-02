package cache

import (
	"github.com/vincent78/butil/utils/objUtil"
	"sync"
)

type ObjPool struct {
	sync.Mutex
	Inuse     []interface{}
	Available []interface{}
	New       func() interface{}
	Reset     func(obj interface{})
}

func NewObjPool(new func() interface{}) *ObjPool {
	return &ObjPool{New: new}
}

func (p *ObjPool) Acquire() interface{} {
	p.Lock()
	var object interface{}

	if len(p.Inuse) != 0 {
		object = p.Available[0]
		// TODO: remove one of Available
		p.Available = append(p.Available[:0], p.Available[1:]...)
		// TODO: add one in Inuse
		p.Inuse = append(p.Inuse, object)
	} else {
		object = p.New()
		p.Inuse = append(p.Inuse, object)
	}
	p.Unlock()
	return object
}

func (p *ObjPool) Release(object interface{}) {
	p.Lock()
	if !objUtil.IsNil(object) && !objUtil.IsNil(p.Reset) {
		p.Reset(object)
	}
	p.Available = append(p.Available, object)
	for i, v := range p.Inuse {
		if v == object {
			// TODO: remove object from available list
			p.Inuse = append(p.Inuse[:i], p.Inuse[i+1:]...)
			break
		}
	}
	p.Unlock()
}
