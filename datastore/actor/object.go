package actor

func (a *Actor) GetObjects() ([]string, error) {
	n, ok := (*a)["objects"]
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

func (a *Actor) GetObjectsCount() int {
	objects, err := a.GetObjects()
	if err != nil {
		return 0
	}

	return len(objects)
}

func (a *Actor) AppendObject(objectID string) {
	objects, err := a.GetObjects()
	if err != nil {
		objects = []string{}
	}

	objects = append(objects, objectID)
	(*a)["objects"] = objects
}
