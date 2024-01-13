package eventsmgr

import (
	"fmt"

	hook "github.com/robotn/gohook"
)

type EventManager struct{}

func (e *EventManager) HandleEvent(ev *hook.Event){
	fmt.Println(ev)
}