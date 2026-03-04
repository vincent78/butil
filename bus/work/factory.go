package work

import "sync"

type Factory struct {
	shop sync.Map
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) OnDuty(w *Worker) {}

func (f *Factory) OffDuty(w *Worker) {}
