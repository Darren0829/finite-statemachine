package finitestatemachine

import (
	"fmt"
	"golang.org/x/exp/slices"
)

type Transformation struct {
	// todo 重复接收同一事件怎么处理？忽略
	Events   Events
	From     State
	To       State
	Commands Commands

	waitingEvents  []string
	receivedEvents Events
}

func (t *Transformation) try(cur State, e Event) bool {
	if cur != t.From {
		return false
	}

	idx := slices.Index(t.waitingEvents, e.ID)
	if idx == -1 {
		return false
	}
	t.waitingEvents = append(t.waitingEvents[:idx], t.waitingEvents[idx+1:]...)
	t.receivedEvents = append(t.receivedEvents, e)

	return len(t.waitingEvents) == 0
}

func (t *Transformation) reset() {
	t.waitingEvents = t.Events.eventIDs()
	t.receivedEvents = Events{}
}

type Transformations []Transformation

func (ts Transformations) enforce(cur State, e Event) (bool, Transformation) {
	fmt.Println(fmt.Sprintf("current state: [%s], accept event: [%v]", cur, e))
	for _, t := range ts {
		ok := t.try(cur, e)
		if ok {
			return true, t
		}
	}
	return false, Transformation{}
}

func (ts Transformations) reset() {
	for _, t := range ts {
		t.reset()
	}
}
