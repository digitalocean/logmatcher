package matcher

import (
	"fmt"

	"github.com/digitalocean/captainslog"
)

// Matcher is the generic interface type for different exclusion behaviors.
type Matcher interface {
	Matches(captainslog.SyslogMsg) bool
	String() string
	Encode(map[string]interface{})
	Decode(map[string]interface{}) error
}

// Matchers is a slice of Matcher.
type Matchers []Matcher

// DecodeArray is a generic method to decode an array of matchers.
func DecodeArray(s []interface{}) (Matchers, error) {
	var o Matchers

	for _, v := range s {
		if m, ok := v.(map[string]interface{}); ok {
			matcher, err := Decode(m)
			if err != nil {
				return nil, err
			}
			o = append(o, matcher)
		} else {
			return nil, fmt.Errorf("failed to decode matchers, found item that wasn't a map")
		}
	}

	return o, nil
}

// Decode is a generic method to decode a matcher.
func Decode(m map[string]interface{}) (Matcher, error) {
	if len(m) > 1 {
		return nil, fmt.Errorf("failed to decode matcher, found too many keys")
	}

	for k, v := range m {
		matcher, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to decode matcher into map")
		}

		switch k {
		case "unary_op":
			unaryop := &UnaryOp{}
			err := unaryop.Decode(matcher)
			return unaryop, err
		case "n_ary_op":
			naryop := &NAryOp{}
			err := naryop.Decode(matcher)
			return naryop, err
		case "value_matcher":
			value := &Value{}
			err := value.Decode(matcher)
			return value, err
		case "kv_matcher":
			kv := &KV{}
			err := kv.Decode(matcher)
			return kv, err
		case "facility_matcher":
			facility := &Facility{}
			err := facility.Decode(matcher)
			return facility, err
		case "severity_matcher":
			severity := &Severity{}
			err := severity.Decode(matcher)
			return severity, err
		case "timestamp_matcher":
			timestamp := &Timestamp{}
			err := timestamp.Decode(matcher)
			return timestamp, err
		}
	}

	return nil, fmt.Errorf("failed to decode matcher, found no valid types")
}

// Encode is a generic method to encode a matcher into a map.
func Encode(in Matcher, out map[string]interface{}) {
	switch in.(type) {
	case *Severity:
		out["severity_matcher"] = make(map[string]interface{})
		m := in.(*Severity)
		m.Encode(out["severity_matcher"].(map[string]interface{}))
	case *Facility:
		out["facility_matcher"] = make(map[string]interface{})
		m := in.(*Facility)
		m.Encode(out["facility_matcher"].(map[string]interface{}))
	case *Timestamp:
		out["timestamp_matcher"] = make(map[string]interface{})
		m := in.(*Timestamp)
		m.Encode(out["timestamp_matcher"].(map[string]interface{}))
	case *KV:
		out["kv_matcher"] = make(map[string]interface{})
		m := in.(*KV)
		m.Encode(out["kv_matcher"].(map[string]interface{}))
	case *Value:
		out["value_matcher"] = make(map[string]interface{})
		m := in.(*Value)
		m.Encode(out["value_matcher"].(map[string]interface{}))
	case *UnaryOp:
		out["unary_op"] = make(map[string]interface{})
		m := in.(*UnaryOp)
		m.Encode(out["unary_op"].(map[string]interface{}))
	case *NAryOp:
		out["n_ary_op"] = make(map[string]interface{})
		m := in.(*NAryOp)
		m.Encode(out["n_ary_op"].(map[string]interface{}))
	}
}
