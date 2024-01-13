package events

import (
	hook "github.com/robotn/gohook"
)

type EventCatcher struct {
	sendEventCh chan *hook.Event
	doneCh chan struct{}
	readEventCh chan hook.Event
}

func New() *EventCatcher{
	return &EventCatcher{
		sendEventCh: make(chan *hook.Event, 1),
		doneCh: make(chan struct{}),
	}
}

func (e *EventCatcher) isAllowedEvent(event *hook.Event) bool{
	switch event.Kind{
	case hook.MouseDown, hook.MouseMove, hook.MouseWheel, hook.MouseDrag, hook.KeyDown:
		return true
	default:
		return false
	}
}

func (e *EventCatcher) cleanup(){
	close(e.sendEventCh)
	close(e.readEventCh)
	hook.StopEvent()

}

func (e *EventCatcher) Stop(){
	close(e.doneCh)
}

func (e *EventCatcher) Start() chan *hook.Event{
	
	e.readEventCh = hook.Start()

	go func(){
		for{
			select{
			case <- e.doneCh:
				e.cleanup()
				return

			case ev := <-e.readEventCh:
				if e.isAllowedEvent(&ev){
					e.sendEventCh <- &ev
				}
			}

		}
	}()

	return e.sendEventCh
}
