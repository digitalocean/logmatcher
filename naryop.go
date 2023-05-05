package matcher

import (
	"bytes"
	"fmt"

	"github.com/digitalocean/captainslog"
)

// NAryOpType is the enum class for representing n-ary operation types.
type NAryOpType int

// N-ary operation types.
const (
	And NAryOpType = iota
	Or
)

// String converts a NAryOpType to its corresponding string representation.
func (o NAryOpType) String() string {
	switch o {
	case And:
		return "and"
	case Or:
		return "or"
	default:
		return "invalid type"
	}
}

// FromString converts the NAryOpType to the value corresponding to the supplied
// string representation.
func (o *NAryOpType) FromString(s string) error {
	switch s {
	case "and":
		*o = And
	case "or":
		*o = Or
	default:
		return fmt.Errorf("failed to convert string to NAryOpType")
	}

	return nil
}

// NAryOp encapsulates a type of n-ary operation and the slice of abstract
// Matchers on which it applies.
type NAryOp struct {
	Type     NAryOpType
	Matchers Matchers
}

// NewNAryOp returns a new NAryOp with the specified operation type and
// set of exclusions.
func NewNAryOp(t NAryOpType, v ...Matcher) *NAryOp {
	return &NAryOp{
		Type:     t,
		Matchers: v,
	}
}

// String converts an NAryOp to its corresponding string representation.
func (o NAryOp) String() string {
	var b bytes.Buffer
	b.WriteByte('(')
	for i, m := range o.Matchers {
		if i != 0 {
			b.WriteByte(' ')
		}
		b.WriteString(m.String())
		if i != len(o.Matchers)-1 {
			b.WriteByte(' ')
			b.WriteString(o.Type.String())
		}
	}
	b.WriteByte(')')
	return b.String()
}

// Matches returns true if the NAryOp matches the supplied SyslogMsg.
func (o *NAryOp) Matches(m captainslog.SyslogMsg) bool {
	switch o.Type {
	case And:
		for _, v := range o.Matchers {
			if !v.Matches(m) {
				return false
			}
		}
		return true
	case Or:
		for _, v := range o.Matchers {
			if v.Matches(m) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// Decode decodes the map to an NAryOp type.
func (o *NAryOp) Decode(m map[string]interface{}) error {
	foundType := false
	foundMatchers := false
	for k, v := range m {
		switch k {
		case "type":
			foundType = true

			if t, ok := v.(string); ok {
				if err := o.Type.FromString(t); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed to decode n-ary op, type is not a string")
			}
		case "matchers":
			foundMatchers = true

			if matchers, ok := v.([]interface{}); ok {
				vals, err := DecodeArray(matchers)
				if err != nil {
					return err
				}
				o.Matchers = vals
			} else {
				return fmt.Errorf("failed to decode n-ary op, matchers is not a slice")
			}
		}
	}

	if !(foundType && foundMatchers) {
		return fmt.Errorf("failed to decode n-ary op, missing fields")
	}

	return nil
}

// Encode encodes the NAryOp to a map.
func (o *NAryOp) Encode(out map[string]interface{}) {
	out["type"] = o.Type.String()
	var v []map[string]interface{}
	for _, item := range o.Matchers {
		x := make(map[string]interface{})
		Encode(item, x)
		v = append(v, x)
	}
	out["matchers"] = v
}
