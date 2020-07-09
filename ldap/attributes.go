package main

type Attributes struct {
	Map map[string][]string
}

func (a *Attributes) Get(key string) ([]string, bool) {
	elem, ok := a.Map[key]
	return elem, ok
}

func (a *Attributes) GetFirst(key string) (string, bool) {
	elem := a.Map[key]
	if elem != nil {
		if len(elem) > 0 {
			return elem[0], true
		}
		return "", true
	}
	return "", false
}
