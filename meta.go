package activity

import "reflect"

type Registry struct {
	idToValue map[string]reflect.Value
	valueToID map[reflect.Value]string
}

func NewRegistry() *Registry {
	r := Registry{
		idToValue: make(map[string]reflect.Value),
		valueToID: make(map[reflect.Value]string),
	}
	return &r
}

func (r *Registry) Add(funcID string, method interface{}) *Registry {
	r.idToValue[funcID] = reflect.ValueOf(method)
	r.valueToID[reflect.ValueOf(method)] = funcID
	return r
}

func (r *Registry) FindFuncID(method reflect.Value) string {
	return r.valueToID[method]
}
