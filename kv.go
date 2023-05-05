package matcher

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/digitalocean/captainslog"
)

// KV represents a key-value matcher
type KV struct {
	Key       string
	MatchType MatchType
	Value     interface{}
}

// NewKV returns a new KV with the specified key, match type, and
// string value.
func NewKV(k string, m MatchType, v interface{}) *KV {
	var vNew interface{}

	if reflect.TypeOf(v).Kind() == reflect.Int {
		vNew = float64(reflect.ValueOf(v).Int())
	} else {
		vNew = v
	}

	return &KV{
		Key:       k,
		MatchType: m,
		Value:     vNew,
	}
}

// String converts a KV to its corresponding string representation.
func (kv KV) String() string {
	if reflect.ValueOf(kv.Value).Kind() == reflect.String {
		return fmt.Sprintf("kv(\"%s\", %s, \"%s\")", kv.Key, kv.MatchType, kv.Value)
	}
	return fmt.Sprintf("kv(\"%s\", %s, %v)", kv.Key, kv.MatchType, kv.Value)
}

// Matches returns true if the KV matches the supplied SyslogMsg.
func (kv *KV) Matches(m captainslog.SyslogMsg) bool {
	if !m.IsJSON {
		return false
	}

	keyChain := strings.Split(kv.Key, ".")
	var next interface{}
	next = m.JSONValues

	for _, key := range keyChain {
		typ := reflect.TypeOf(next)
		if typ.Kind() != reflect.Map || typ.Key().Kind() != reflect.String || typ.Elem().Kind() != reflect.Interface {
			return false
		}
		mapVal := reflect.ValueOf(next).MapIndex(reflect.ValueOf(key))
		if !mapVal.IsValid() {
			return false
		}
		next = mapVal.Interface()
	}

	v := reflect.ValueOf(next)
	t := reflect.TypeOf(next)
	kvr := reflect.ValueOf(kv.Value)

	switch kvr.Kind() {
	case reflect.String:
		comp := kvr.String()

		if v.Kind() != reflect.String {
			return false
		}
		val := v.String()

		switch kv.MatchType {
		case ExactMatch, Equals:
			return comp == val
		case PrefixMatch:
			return strings.HasPrefix(val, comp)
		case Contains:
			return strings.Contains(val, comp)
		case Regex:
			matched, _ := regexp.MatchString(comp, val)
			return matched
		}
	case reflect.Float64:
		comp := kvr.Float()

		var val float64
		if t.String() == "json.Number" {
			jsonNum, ok := v.Interface().(json.Number)
			if !ok {
				return false
			}
			tVal, err := jsonNum.Float64()
			if err != nil {
				return false
			}
			val = tVal
		} else if v.Kind() == reflect.Float64 {
			val = v.Float()
		} else {
			return false
		}

		switch kv.MatchType {
		case Equals:
			return comp == val
		case LessThan:
			return val < comp
		case LessThanEqual:
			return val <= comp
		case GreaterThan:
			return val > comp
		case GreaterThanEqual:
			return val >= comp
		}
	case reflect.Bool:
		comp := kvr.Bool()

		if v.Kind() != reflect.Bool {
			return false
		}
		val := v.Bool()

		switch kv.MatchType {
		case Equals:
			return comp == val
		}
	}

	return false
}

// Decode is a helper function to decode a matcher to a key-value matcher.
func (kv *KV) Decode(m map[string]interface{}) error {
	foundKey := false
	foundMatchType := false
	foundValue := false
	for k, v := range m {
		switch k {
		case "key":
			foundKey = true

			if val, ok := v.(string); ok {
				kv.Key = val
			} else {
				return fmt.Errorf("failed to decode kv matcher, key is not a string")
			}
		case "match_type":
			foundMatchType = true

			if mt, ok := v.(string); ok {
				if err := kv.MatchType.FromString(mt); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed to decode kv matcher, match_type is not a string")
			}
		case "str_value":
			if v != nil {
				foundValue = true
				kv.Value = v
			}
		case "num_value":
			if v != nil {
				foundValue = true
				val := reflect.ValueOf(v)
				switch val.Kind() {
				case reflect.Int, reflect.Int64:
					kv.Value = float64(val.Int())
				case reflect.Float32, reflect.Float64:
					kv.Value = val.Float()
				default:
					foundValue = false
				}
			}
		case "bool_value":
			if v != nil {
				foundValue = true
				kv.Value = v
			}
		}
	}

	if !(foundKey && foundMatchType && foundValue) {
		return fmt.Errorf("failed to decode kv matcher, missing fields")
	}

	return nil
}

// Encode encodes a key-value object into a matcher map.
func (kv *KV) Encode(out map[string]interface{}) {
	out["key"] = kv.Key
	out["match_type"] = kv.MatchType.String()

	out["str_value"] = nil
	out["num_value"] = nil
	out["bool_value"] = nil

	switch reflect.ValueOf(kv.Value).Kind() {
	case reflect.String:
		out["str_value"] = kv.Value
	case reflect.Float64:
		out["num_value"] = kv.Value
	case reflect.Bool:
		out["bool_value"] = kv.Value
	}
}
