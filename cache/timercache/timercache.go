package timercache

import "time"

type TimerCache[T any] struct {
	data         *T
	Refresh      func() T
	expired      time.Duration
	once         bool
	loop         bool
	updateThread func()
}

func (h *TimerCache[T]) Start() {
	if !h.loop {
		h.loop = true
		go h.updateThread()
	}
}

func (h *TimerCache[T]) Update() {
	t := h.Refresh()
	h.data = &t
}

func (h *TimerCache[T]) Stop() {
	h.loop = false
}

func (h *TimerCache[T]) GetData() *T {
	return h.data
}

func (h *TimerCache[T]) TryGetData(timeout time.Duration) *T {
	now := time.Now()
	for {
		if h.data == nil && time.Since(now) < timeout {
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}

	}
	return h.data
}

func (h *TimerCache[T]) GetInterval() time.Duration {
	return h.expired
}

func (h *TimerCache[T]) GetOnce() bool {
	return h.once
}

func (h *TimerCache[T]) setOnce(once bool) {
	h.once = once
}

func Build[T any](expired time.Duration, refresh func() T) *TimerCache[T] {
	instance := &TimerCache[T]{expired: expired, Refresh: refresh}
	instance.updateThread = func() {
		for instance.loop {
			instance.Update()
			instance.once = true
			time.Sleep(instance.expired)
		}
	}
	return instance
}
