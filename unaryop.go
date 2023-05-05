package matcher

import (
	"fmt"

	"github.com/digitalocean/captainslog"
)

// UnaryOpType is the enum class for representing unary operation types.
type UnaryOpType int

// Unary operation types.
const (
	Not UnaryOpType = iota
)

// String converts a UnaryOpType to its corresponding string representation.
func (o UnaryOpType) String() string {
	switch o {
	case Not:
		return "not"
	default:
		return "invalid type"
	}
}

// FromString converts the UnaryOpType to the value corresponding to the supplied
// string representation.
func (o *UnaryOpType) FromString(s string) error {
	switch s {
	case "not":
		*o = Not
	default:
		return fmt.Errorf("failed to convert string to UnaryOpType")
	}

	return nil
}

// UnaryOp encapsulates a type of unary operation and the abstract Matcher on
// which it applies.
type UnaryOp struct {
	Type    UnaryOpType
	Matcher Matcher
}

// NewUnaryOp returns a new UnaryOp with the specified operation type and
// exclusion.
func NewUnaryOp(t UnaryOpType, m Matcher) *UnaryOp {
	return &UnaryOp{
		Type:    t,
		Matcher: m,
	}
}

// String converts a UnaryOp to its corresponding string representation.
func (o UnaryOp) String() string {
	return fmt.Sprintf("%s %s", o.Type, o.Matcher)
}

// Matches returns true if the UnaryOp matches the supplied SyslogMsg.
func (o *UnaryOp) Matches(m captainslog.SyslogMsg) bool {
	switch o.Type {
	case Not:
		return !o.Matcher.Matches(m)
	default:
		return false
	}
}

// Decode decodes the matcher map into a UnaryOp type.
func (o *UnaryOp) Decode(m map[string]interface{}) error {
	foundType := false
	foundMatcher := false
	for k, v := range m {
		switch k {
		case "type":
			foundType = true

			if t, ok := v.(string); ok {
				if err := o.Type.FromString(t); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed to decode unary op, type is not a string")
			}
		case "matcher":
			foundMatcher = true

			if matcher, ok := v.(map[string]interface{}); ok {
				val, err := Decode(matcher)
				if err != nil {
					return err
				}
				o.Matcher = val
			} else {
				return fmt.Errorf("failed to decode unary op, matcher is not a map")
			}
		}
	}

	if !(foundType && foundMatcher) {
		return fmt.Errorf("failed to decode unary op, missing fields")
	}

	return nil
}

// Encode encodes a UnaryOp into a matcher map.
func (o *UnaryOp) Encode(out map[string]interface{}) {
	out["type"] = o.Type.String()
	out["matcher"] = make(map[string]interface{})
	Encode(o.Matcher, out["matcher"].(map[string]interface{}))
}
