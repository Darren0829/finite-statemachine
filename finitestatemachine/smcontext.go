package finitestatemachine

type SMContext struct {
	id             string
	name           string
	state          State
	variable       interface{}
	receivedEvents []Event
}

func (c SMContext) GetID() string {
	return c.id
}

func (c SMContext) GetName() string {
	return c.name
}

func (c SMContext) GetCurrentState() State {
	return c.state
}

func (c SMContext) GetVariable() interface{} {
	return c.variable
}
