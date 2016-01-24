package saga

import (
	"reflect"
)

// ParamData presents sub-transaction input parameter data.
// This structure used to store and restore tx input data into log.
type ParamData struct {
	ParamType string `json:"paramType,omitempty"`
	Data      string `json:"data,omitempty"`
}

// MarshalParam convert args into ParamData.
// This method will lookup typeName in given SEC.
func MarshalParam(sec *ExecutionCoordinator, args []interface{}) []ParamData {
	p := make([]ParamData, 0, len(args))
	for _, arg := range args {
		typ := sec.MustFindParamName(reflect.ValueOf(arg).Type())
		p = append(p, ParamData{
			ParamType: typ,
			Data:      mustMarshal(arg),
		})
	}
	return p
}

// UnmarshalParam convert ParamData back to parameter values to function call usage.
// This method will lookup reflect.Type in given SEC.
func UnmarshalParam(sec *ExecutionCoordinator, paramData []ParamData) []reflect.Value {
	var values []reflect.Value
	for _, param := range paramData {
		ptyp := sec.MustFindParamType(param.ParamType)
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
