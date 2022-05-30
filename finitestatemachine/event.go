package finitestatemachine

type Event struct {
	ID      string
	Message interface{}
}

type Events []Event

func (es Events) GetMessage(ID string) (bool, interface{}) {
	for _, e := range es {
		if e.ID == ID {
			return true, e.Message
		}
	}
	return false, nil
}

func (es Events) index(e Event) int {
	for i, o := range es {
		if o.ID == e.ID {
			return i
		}
	}
	return -1
}

func (es Events) eventIDs() []string {
	var eventIDs []string
	for _, o := range es {
		eventIDs = append(eventIDs, o.ID)
	}
	return eventIDs
}
