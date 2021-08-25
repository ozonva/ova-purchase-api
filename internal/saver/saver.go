package saver

import (
	"errors"
	"github.com/ozonva/ova-purchase-api/internal/flusher"
	"github.com/ozonva/ova-purchase-api/internal/purchase"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

var CapacityNotValidError = errors.New("capacity must be greater than 0")
var TimeoutNotValidError = errors.New("timeout must be greater than 0")
var CapacityReachedError = errors.New("can't save purchase, saver's capacity is reached")

type Saver interface {
	Save(purchase purchase.Purchase) error
	Close()
}

type saver struct {
	mutex    sync.Mutex
	flusher  flusher.Flusher
	timeout  time.Duration
	buffer   []purchase.Purchase
	done     chan struct{}
	capacity uint
}

// NewSaver возвращает Saver с поддержкой переодического сохранения
func NewSaver(
	capacity uint,
	flusher flusher.Flusher,
	timeout time.Duration,
) (Saver, error) {
	if timeout <= 0 {
		return nil, TimeoutNotValidError
	}
	if capacity <= 0 {
		return nil, CapacityNotValidError
	}
	instance := &saver{
		mutex:    sync.Mutex{},
		flusher:  flusher,
		timeout:  timeout,
		buffer:   make([]purchase.Purchase, 0, capacity),
		done:     make(chan struct{}),
		capacity: capacity,
	}
	go instance.start()
	return instance, nil
}

func (s *saver) Save(purchase purchase.Purchase) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.buffer) >= int(s.capacity) {
		return CapacityReachedError
	}
	s.buffer = append(s.buffer, purchase)
	return nil
}

func (s *saver) Close() {
	s.done <- struct{}{}
	close(s.done)
}

func (s *saver) start() {
	timer := time.NewTimer(s.timeout)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			if len(s.buffer) > 0 {
				notFlushed := s.flusher.Flush(s.buffer)
				newBuffer := make([]purchase.Purchase, 0, s.capacity)
				s.buffer = append(newBuffer, notFlushed...)
			}
		case <-s.done:
			notFlushed := s.flusher.Flush(s.buffer)
			if len(notFlushed) > 0 {
				log.Error().Interface("Purchases", notFlushed).Msg("failed to flush purchases %s")
			}
			return
		}
	}
}
