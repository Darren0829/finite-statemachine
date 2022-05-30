package finitestatemachine

import (
	"context"
	"fmt"
	"golang.org/x/exp/slices"
)

type StateMachineGroup struct {
	name          string
	stateMachines StateMachines

	ctx    context.Context
	cancel context.CancelFunc
}

func NewStateMachineGroup(name string) StateMachineGroup {
	ctx, cancel := context.WithCancel(context.Background())
	return StateMachineGroup{
		name:   name,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (g *StateMachineGroup) AddMachine(m *StateMachine) {
	g.stateMachines = append(g.stateMachines, m)
}

func (g *StateMachineGroup) AddMachines(ms StateMachines) {
	g.stateMachines = append(g.stateMachines, ms...)
}

func (g *StateMachineGroup) Start() error {
	err := GetStateMachineController().launch(g)
	if err != nil {
		return err
	}
	return nil
}

func (g *StateMachineGroup) Stop() error {
	err := GetStateMachineController().destroy(g)
	if err != nil {
		return err
	}
	return nil
}

func (g *StateMachineGroup) run() error {
	if len(g.stateMachines) == 0 {
		return fmt.Errorf("state machine group must contain at least one machine")
	}
	g.stateMachines.run(g.ctx)
	return nil
}

func (g *StateMachineGroup) handle(e Event) {
	g.stateMachines.handle(e)
}

func (g *StateMachineGroup) destroy() error {
	g.cancel()
	slices.Delete(g.stateMachines, 0, len(g.stateMachines)-1)
	return nil
}

type StateMachines []*StateMachine

func (ms StateMachines) run(ctx context.Context) {
	for _, m := range ms {
		go m.run(ctx)
	}
}

func (ms StateMachines) handle(e Event) {
	for _, m := range ms {
		m.event <- e
	}
}
