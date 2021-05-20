package models

type Event interface {
	Type() string
}

type Subscribe struct {
	ch   chan<- Event
	Type string
}

var (
	ch         chan Event
	subscribes map[string]*Subscribe
)

func init() {
	ch = make(chan Event, 50)
	subscribes = make(map[string]*Subscribe)

	go dispatcher()
}

func dispatcher() {
	for {
		select {
		case evt := <-ch:
			dispatch(evt)
		}
	}
}

func dispatch(evt Event) {
	for _, sub := range subscribes {
		if sub.Type != "" {
			if evt.Type() == sub.Type {
				sub.ch <- evt
			}
		} else {
			sub.ch <- evt
		}
	}
}

func Pub(evt Event) {
	ch <- evt
}

func Sub(id string) (<-chan Event, func()) {
	var (
		ch       = make(chan Event, 50)
		cancelFn func()
	)

	subscribes[id] = &Subscribe{
		ch: ch,
	}

	return ch, cancelFn
}
