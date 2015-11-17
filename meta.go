package activity

import "reflect"

type Registry map[string]reflect.Method

func NewRegistry() Registry {
	return make(Registry)
}

func (r *Registry) Add(funcID string, method reflect.Method) *Registry {
	(*r)[funcID] = method
	return r
}