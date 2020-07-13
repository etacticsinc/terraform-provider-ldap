package main

type Attributes struct {
	Map map[string][]string
}

func (a *Attributes) Get(key string) []string {
	return a.Map[key]
}

func (a *Attributes) GetFirst(key string) string {
	elem := a.Map[key]
	if elem != nil && len(elem) > 0 {
		return elem[0]
	}
	return ""
}

func (a *Attributes) ForEach(fn func(string, []string)) {
	for key, elem := range a.Map {
		if elem != nil && len(elem) > 0 && elem[0] != "" {
			fn(key, elem)
		}
	}
}

func (a *Attributes) HasValue(key string) bool {
	elem, ok := a.Map[key]
	return ok && elem != nil && len(elem) > 0 && elem[0] != ""
}

func (a *Attributes) Keys() []string {
	keys := make([]string, len(a.Map))
	i := 0
	for key := range a.Map {
		keys[i] = key
		i++
	}
	return keys
}
