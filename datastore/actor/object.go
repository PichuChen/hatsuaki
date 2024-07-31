package actor

func (a *Actor) GetOutboxObjects() ([]string, error) {
	n, ok := (*a)["outbox"]
	if !ok {
		return []string{}, nil
	}
	var objects []string
	if s, ok := n.([]string); ok {
		return s, nil
	}

	is, _ := n.([]interface{})
	objects = make([]string, len(is))
	for i, v := range is {
		objects[i] = v.(string)
	}
	return objects, nil
}

func (a *Actor) GetOutboxObjectsCount() int {
	objects, err := a.GetOutboxObjects()
	if err != nil {
		return 0
	}

	return len(objects)
}

func (a *Actor) AppendOutboxObject(objectID string) {
	objects, err := a.GetOutboxObjects()
	if err != nil {
		objects = []string{}
	}

	objects = append(objects, objectID)
	(*a)["outbox"] = objects
}

func (a *Actor) AppendInboxObject(objectID string) {
	objects, err := a.GetInboxObjects()
	if err != nil {
		objects = []string{}
	}

	objects = append(objects, objectID)
	(*a)["inbox"] = objects
}

func (a *Actor) GetInboxObjects() ([]string, error) {
	n, ok := (*a)["inbox"]
	if !ok {
		return []string{}, nil
	}

	is := n.([]interface{})
	objects := make([]string, len(is))
	for i, v := range is {
		objects[i] = v.(string)
	}
	return objects, nil
}

func (a *Actor) GetInboxObjectsCount() int {
	objects, err := a.GetInboxObjects()
	if err != nil {
		return 0
	}

	return len(objects)
}
