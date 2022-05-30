package finitestatemachine

import (
	"fmt"
	"sync"
)

type stateMachineController struct {
	groups map[string]*StateMachineGroup
	m      sync.Mutex
}

var controller *stateMachineController
var once sync.Once

func GetStateMachineController() *stateMachineController {
	once.Do(func() {
		controller = &stateMachineController{
			groups: map[string]*StateMachineGroup{},
			m:      sync.Mutex{},
		}
	})
	return controller
}

func (c *stateMachineController) Handle(groupName string, e Event) error {
	g := c.groups[groupName]
	if g == nil {
		return fmt.Errorf("state machine group not exist")
	}

	g.handle(e)
	return nil
}

func (c *stateMachineController) launch(g *StateMachineGroup) error {
	c.m.Lock()
	defer c.m.Unlock()

	name := g.name
	if c.groups[name] != nil {
		return fmt.Errorf("state machine group [%s] already exits", name)
	}

	err := g.run()
	if err != nil {
		return err
	}

	c.groups[name] = g
	fmt.Println(fmt.Sprintf("state machine group [%s] is running", name))
	return nil
}

func (c *stateMachineController) destroy(g *StateMachineGroup) error {
	c.m.Lock()
	defer c.m.Unlock()

	name := g.name
	delete(c.groups, name)
	return g.destroy()
}
