package event

import (
	"github.com/qur/gopy/lib"

	"github.com/gemrs/willow/log"
)

type Event struct {
	py.BaseObject

	key       string
	observers map[int]Observer
	log       log.Log
}

func (e *Event) Init(key string) {
	e.key = key
	e.observers = make(map[int]Observer)
	e.log = log.New("event", log.MapContext{"event": key})
}

func (e *Event) Key() string {
	return e.key
}

func (e *Event) Register(o Observer) {
	e.observers[o.Id()] = o
}

func (e *Event) Unregister(o Observer) {
	delete(e.observers, o.Id())
}

func (e *Event) NotifyObservers(args ...interface{}) {
	for _, observer := range e.observers {
		observer.Notify(e, args...)
	}
}
