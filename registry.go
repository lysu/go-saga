package saga

import (
	"reflect"
)

type Registry struct {
	idToValue  map[string]reflect.Value
	valueToID  map[reflect.Value]string
	nameToType map[string]reflect.Type
	typeToName map[reflect.Type]string
}

func NewRegistry() *Registry {
	r := Registry{
		idToValue:  make(map[string]reflect.Value),
		valueToID:  make(map[reflect.Value]string),
		nameToType: make(map[string]reflect.Type),
		typeToName: make(map[reflect.Type]string),
	}
	return &r
}

func (r *Registry) Add(funcID string, method interface{}) *Registry {
	funcValue := reflect.ValueOf(method)
	if funcValue.Kind() != reflect.Func {
		panic("Regist object must be a func")
	}
	if funcValue.Type().NumIn() < 1 ||
		funcValue.Type().In(0) != activityContextType {
		panic("First argument must use ActivityContext.")
	}
	r.idToValue[funcID] = funcValue
	r.valueToID[funcValue] = funcID
	r.addParams(funcValue)
	return r
}

func (r *Registry) addParams(funcValue reflect.Value) {
	funcType := funcValue.Type()
	for i := 0; i < funcType.NumIn(); i++ {
		paramType := funcType.In(i)
		r.nameToType[paramType.Name()] = paramType
		r.typeToName[paramType] = paramType.Name()
	}
	for i := 0; i < funcType.NumOut(); i++ {
		returnType := funcType.Out(i)
		r.nameToType[returnType.Name()] = returnType
		r.typeToName[returnType] = returnType.Name()
	}
}

func (r *Registry) FindFuncID(method reflect.Value) string {
	return r.valueToID[method]
}

func (r *Registry) FindFunction(funcID string) *reflect.Value {
	if f, ok := r.idToValue[funcID]; ok {
		return &f
	}
	return nil
}

func (r *Registry) FindTypeName(typ reflect.Type) string {
	return r.typeToName[typ]
}

func (r *Registry) FindType(typeName string) (reflect.Type, bool) {
	f, ok := r.nameToType[typeName]
	return f, ok
}
