package finitestatemachine

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type State string

type StateMachine struct {
	id              string
	name            string
	state           State
	states          []State
	transformations Transformations
	context         SMContext

	//sync.Mutex
	running  bool
	shutdown chan bool
	event    chan Event
}

func NewStateMachine(name string, s State, states []State) StateMachine {
	if slices.Index(states, s) == -1 {
		states = append(states, s)
	}

	id := uuid.NewString()
	sm := StateMachine{
		id:     id,
		name:   name,
		state:  s,
		states: states,
		context: SMContext{
			id:    id,
			name:  name,
			state: s,
		},
		shutdown: make(chan bool),
		event:    make(chan Event),
	}
	return sm
}

func (s *StateMachine) Set(receive Events, cur, dst State, cmd Commands) *StateMachine {
	t := Transformation{
		Events:   receive,
		From:     cur,
		To:       dst,
		Commands: cmd,
	}
	t.reset()
	s.transformations = append(s.transformations, t)
	return s
}

func (s *StateMachine) ID() string {
	return s.id
}

func (s *StateMachine) Matrix() {

}

func (s *StateMachine) run(ctx context.Context) {
	if s.running {
		return
	}
	s.running = true
	fmt.Println(fmt.Sprintf("state machine [%s] is running...", s.name))

	for {
		select {
		case e := <-s.event:
			s.accept(e)
		case <-ctx.Done():
			s.destroy()
			return
		case <-s.shutdown:
			s.reset()
			return
		}
	}
}

func (s *StateMachine) accept(e Event) {
	fmt.Println("accept a event", e)
	should, t := s.transformations.enforce(s.state, e)
	if should {
		s.transformations.reset()
		_ = s.turnInto(t.To)
		t.Commands.exec(&s.context, t.receivedEvents)
	}
}

func (s *StateMachine) turnInto(state State) error {
	if slices.Index(s.states, state) == -1 {
		return fmt.Errorf("init state [%s] is none of valid states", state)
	}
	s.state = state
	s.context.state = state
	return nil
}

func (s *StateMachine) reset() {
	s.running = false
}

func (s *StateMachine) destroy() {
	fmt.Println(fmt.Sprintf("state machine [%s] is destroying...", s.name))
	close(s.event)
	close(s.shutdown)
	s.running = false
	fmt.Println(fmt.Sprintf("state machine[%s] is destroyed", s.name))
}
