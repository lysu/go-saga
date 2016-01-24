package saga

import (
	"reflect"
)

type ParamData struct {
	ParamType string `json:"paramType,omitempty"`
	Data      string `json:"data,omitempty"`
}

func (s *Saga) MarshalParam(args []interface{}) []ParamData {
	p := make([]ParamData, 0, len(args))
	for _, arg := range args {
		typ := s.coordinator.FindParamName(reflect.ValueOf(arg).Type())
		p = append(p, ParamData{
			ParamType: typ,
			Data:      mustMarshal(arg),
		})
	}
	return p
}

func (s *Saga) UnmarshalParam(paramData []ParamData) []reflect.Value {
	var values []reflect.Value
	for _, param := range paramData {
		ptyp := s.coordinator.FindParamType(param.ParamType)
		obj := reflect.New(ptyp).Interface()
		mustUnmarshal([]byte(param.Data), obj)
		objV := reflect.ValueOf(obj)
		if objV.Type().Kind() == reflect.Ptr && objV.Type() != ptyp {
			objV = objV.Elem()
		}
		values = append(values, objV)
	}
	return values
}
