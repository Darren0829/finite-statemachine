package test

import (
	fsm "finite-statemachine/finitestatemachine"
	"fmt"
	"testing"
	"time"
)

func Test_Single_Machine(t *testing.T) {
	g := fsm.NewStateMachineGroup("Group-001")
	states := []fsm.State{
		fsm.State("red"),
		fsm.State("green"),
		fsm.State("yellow"),
	}
	sm001 := fsm.NewStateMachine("Machine-001", fsm.State("initialed"), states)

	events := fsm.Events{}
	events = append(events, fsm.Event{
		ID: "e-init",
	})
	commands := fsm.Commands{}
	commands = append(commands, fsm.Command{
		Name: "初始化命令",
		Fn: func(ctx *fsm.SMContext, es fsm.Events) error {
			fmt.Println(fmt.Sprintf("exec command after state machine [%s] turn into [%s]", ctx.GetName(), ctx.GetCurrentState()))
			return nil
		},
	})
	sm001.Set(events, fsm.State("initialed"), fsm.State("red"), commands)

	events = fsm.Events{}
	events = append(events, fsm.Event{
		ID: "e-turn-into-green",
	})
	sm001.Set(events, fsm.State("red"), fsm.State("green"), commands)

	events = fsm.Events{}
	events = append(events, fsm.Event{
		ID: "e-turn-into-yellow",
	})
	sm001.Set(events, fsm.State("green"), fsm.State("yellow"), commands)

	events = fsm.Events{}
	events = append(events, fsm.Event{
		ID: "e-turn-into-red",
	})
	sm001.Set(events, fsm.State("yellow"), fsm.State("red"), commands)
	g.AddMachine(&sm001)

	err := g.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = fsm.GetStateMachineController().Handle("Group-001", fsm.Event{
		ID: "e-init",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(time.Second * 5)
	err = fsm.GetStateMachineController().Handle("Group-001", fsm.Event{
		ID: "e-turn-into-green",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Second * 5)

	err = g.Stop()
	if err != nil {
		return
	}
	time.Sleep(time.Second * 5)
}
