package parser

import (
	"encoding/json"
	"strconv"
	"time"
)

func parseJSON(data []byte) (*Entry, bool) {
	m := make(map[string]interface{})
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, false
	}
	e := newEntry()
	for k, v := range m {
		setFieldFromIface(e, k, v)
	}
	return e, true
}

func setFieldFromIface(e *Entry, name string, val interface{}) {
	switch val := val.(type) {
	case string:
		if t, err := tryParseTime(val); err == nil {
			e.setField(name, TimeField{Value: t})
		} else if d, err := time.ParseDuration(val); err == nil {
			e.setField(name, DurationField{Value: d})
		} else {
			e.setField(name, StringField{Value: val})
		}
	case float64:
		e.setField(name, NumberField{Value: val})
	case bool:
		e.setField(name, BooleanField{Value: val})
	case []byte:
		e.setField(name, RawField{Value: val})
	case []interface{}:
		for i, arrVal := range val {
			setFieldFromIface(e, name+"["+strconv.Itoa(i)+"]", arrVal)
		}
	case map[string]interface{}:
		for key, iface := range val {
			setFieldFromIface(e, name+"."+key, iface)
		}
	case nil:
		setFieldFromIface(e, name, NilField{})
	default:
		setFieldFromIface(e, name, UnknownField{Value: val})
	}
}
