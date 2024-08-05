package dispatcher

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/khan-lau/go-rabbitmq/internal/logger"
)

// Dispatcher -
type Dispatcher struct {
	log           logger.Logger
	subscribers   map[int]dispatchSubscriber
	subscribersMu *sync.Mutex
}

type dispatchSubscriber struct {
	notifyCancelOrCloseChan chan error
	closeCh                 <-chan struct{}
}

// NewDispatcher -
func NewDispatcher(logger logger.Logger) *Dispatcher {
	return &Dispatcher{
		log:           logger,
		subscribers:   make(map[int]dispatchSubscriber),
		subscribersMu: &sync.Mutex{},
	}
}

// Dispatch -
func (d *Dispatcher) Dispatch(err error) error {
	d.subscribersMu.Lock()
	defer d.subscribersMu.Unlock()
	for _, subscriber := range d.subscribers {
		select {
		case <-time.After(time.Second * 5):
			d.log.Warnf("Unexpected rabbitmq error: timeout in dispatch")
		case subscriber.notifyCancelOrCloseChan <- err:
		}
	}
	return nil
}

// AddSubscriber -
func (d *Dispatcher) AddSubscriber() (<-chan error, chan<- struct{}) {
	const maxRand = math.MaxInt
	const minRand = 0
	id := rand.Intn(maxRand-minRand) + minRand

	closeCh := make(chan struct{})
	notifyCancelOrCloseChan := make(chan error)

	d.subscribersMu.Lock()
	d.subscribers[id] = dispatchSubscriber{
		notifyCancelOrCloseChan: notifyCancelOrCloseChan,
		closeCh:                 closeCh,
	}
	d.subscribersMu.Unlock()

	go func(id int) {
		<-closeCh
		d.subscribersMu.Lock()
		defer d.subscribersMu.Unlock()
		sub, ok := d.subscribers[id]
		if !ok {
			return
		}
		close(sub.notifyCancelOrCloseChan)
		delete(d.subscribers, id)
	}(id)
	return notifyCancelOrCloseChan, closeCh
}
