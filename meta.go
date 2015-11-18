package activity

import "reflect"

type Registry map[string]reflect.Value

func NewRegistry() *Registry {
	r := make(Registry)
	return &r
}

func (r *Registry) Add(funcID string, method interface{}) *Registry {
	(*r)[funcID] = reflect.ValueOf(method)
	return r
}
