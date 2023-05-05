package matcher

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/digitalocean/captainslog"
)

// ValueType is the enum class for representing different value types.
type ValueType int

// Value types.
const (
	Host ValueType = iota
	Program
	Content
)

// String converts a ValueType to its corresponding string representation.
func (t ValueType) String() string {
	switch t {
	case Host:
		return "host"
	case Program:
		return "program"
	case Content:
		return "content"
	default:
		return "invalid type"
	}
}

// FromString converts the ValueType to the value corresponding to the supplied
// string representation.
func (t *ValueType) FromString(s string) error {
	switch s {
	case "host":
		*t = Host
	case "program":
		*t = Program
	case "content":
		*t = Content
	default:
		return fmt.Errorf("failed to convert string to ValueType")
	}

	return nil
}

// Value represents the core of exclusions
type Value struct {
	Type      ValueType
	MatchType MatchType
	Value     string
}

// NewValue returns a new Value with the specified value and match
// types and string value.
func NewValue(t ValueType, m MatchType, v string) *Value {
	return &Value{
		Type:      t,
		MatchType: m,
		Value:     v,
	}
}

// String converts a Value to its corresponding string representation.
func (v Value) String() string {
	return fmt.Sprintf("%s(%s, \"%s\")", v.Type, v.MatchType, v.Value)
}

// Matches returns true if the Value matches the supplied SyslogMsg.
func (v *Value) Matches(m captainslog.SyslogMsg) bool {
	var val string
	switch v.Type {
	case Host:
		val = m.Host
	case Program:
		val = m.Tag.Program
	case Content:
		val = m.Content
	}

	switch v.MatchType {
	case ExactMatch, Equals:
		return v.Value == val
	case PrefixMatch:
		return strings.HasPrefix(val, v.Value)
	case Contains:
		return strings.Contains(val, v.Value)
	case Regex:
		matched, _ := regexp.MatchString(v.Value, val)
		return matched
	}

	return false
}

// Decode decodes the matcher map into a Value type.
func (v *Value) Decode(m map[string]interface{}) error {
	foundType := false
	foundMatchType := false
	foundValue := false
	for k, val := range m {
		switch k {
		case "type":
			foundType = true

			if t, ok := val.(string); ok {
				if err := v.Type.FromString(t); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed to decode value matcher, type is not a string")
			}
		case "match_type":
			foundMatchType = true

			if mt, ok := val.(string); ok {
				if err := v.MatchType.FromString(mt); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed to decode value matcher, match_type is not a string")
			}
		case "value":
			foundValue = true

			if val, ok := val.(string); ok {
				v.Value = val
			} else {
				return fmt.Errorf("failed to decode value matcher, value is not a string")
			}
		}
	}

	if !(foundType && foundMatchType && foundValue) {
		return fmt.Errorf("failed to decode value matcher, missing fields")
	}

	return nil
}

// Encode encodes the Value into a matcher map.
func (v *Value) Encode(out map[string]interface{}) {
	out["type"] = v.Type.String()
	out["match_type"] = v.MatchType.String()
	out["value"] = v.Value
}
