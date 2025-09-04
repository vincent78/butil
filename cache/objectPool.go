package cache

import (
	"sync"

	"github.com/vincent78/butil/utils/objUtil"
)

type ObjPool struct {
	sync.Mutex
	Inuse     []interface{}
	Available []interface{}
	new       func() interface{}
	Reset     func(obj interface{})
	Capacity  int
}

func NewObjPool(new func() interface{}) *ObjPool {
	op := &ObjPool{new: new}
	op.Capacity = 100
	op.Inuse = make([]interface{}, 0)
	op.Available = make([]interface{}, 0)
	return op
}

func (p *ObjPool) Acquire() interface{} {
	p.Lock()
	defer p.Unlock()
	var obj interface{}

	if len(p.Inuse) != 0 && len(p.Available) > 0 {
		obj = p.Available[0]
		p.Available = append(p.Available[:0], p.Available[1:]...)
		p.Inuse = append(p.Inuse, obj)
	} else {
		obj = p.new()
		p.Inuse = append(p.Inuse, obj)
	}
	return obj
}

func (p *ObjPool) Release(object interface{}) {
	p.Lock()
	defer p.Unlock()
	if (len(p.Inuse) + len(p.Available)) <= p.Capacity {
		if !objUtil.IsNil(object) && !objUtil.IsNil(p.Reset) {
			p.Reset(object)
		}
		p.Available = append(p.Available, object)

	} else {
		object = nil
	}
	for i, v := range p.Inuse {
		if v == object {
			p.Inuse = append(p.Inuse[:i], p.Inuse[i+1:]...)
			break
		}
	}
}
