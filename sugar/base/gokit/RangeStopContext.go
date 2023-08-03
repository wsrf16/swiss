package gokit

type RangeStopContext struct {
	stopChan chan interface{}
}

func NewStopContext() *RangeStopContext {
	var s *RangeStopContext
	s.stopChan = make(chan interface{})
	return s
}

func (context RangeStopContext) RunBackground(task func()) {
	go func(ch chan interface{}) {
		for range ch {
			task()
		}
	}(context.stopChan)
}

func (context RangeStopContext) Start() {
	context.stopChan <- 1
}

func (context RangeStopContext) Stop() {
	close(context.stopChan)
}
