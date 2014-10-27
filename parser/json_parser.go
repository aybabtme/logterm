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
			e.setField(name, TimeField{t})
		} else if d, err := time.ParseDuration(val); err == nil {
			e.setField(name, DurationField{d})
		} else {
			e.setField(name, StringField(val))
		}
	case float64:
		e.setField(name, NumberField(val))
	case bool:
		e.setField(name, BooleanField(val))
	case []byte:
		e.setField(name, RawField(val))
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
		e.setField(name, UnknownField(val))
	}
}
