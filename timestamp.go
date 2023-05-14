package matcher

import (
	"fmt"
	"time"

	"github.com/digitalocean/captainslog"
)

// Timestamp represents a syslog timestamp matcher
type Timestamp struct {
	MatchType MatchType
	Timestamp captainslog.Time
}

// NewTimestamp returns a new timestamp matcher with the specified time and matchType
func NewTimestamp(m MatchType, t captainslog.Time) *Timestamp {
	return &Timestamp{
		MatchType: m,
		Timestamp: t,
	}
}

// String converts a Timestamp matcher to its corresponding string representation.
func (t Timestamp) String() string {
	return fmt.Sprintf("timestamp(%s, %s)", t.MatchType.String(), t.Timestamp.Time.Format(time.Stamp))
}

// Matches returns true if the Timestamp aligns with the supplied SyslogMsg timestamp and MatchType
func (t *Timestamp) Matches(m captainslog.SyslogMsg) bool {
	// Comparison operators that include equals work the same as ones without when comparing time
	switch t.MatchType {
	case Equals:
		return t.Timestamp.Time.Equal(m.Time)
	case LessThan, LessThanEqual:
		return t.Timestamp.Time.Before(m.Time)
	case GreaterThan, GreaterThanEqual:
		return t.Timestamp.Time.After(m.Time)
	default:
		return false
	}
}

// Decode decodes a matcher map into a Timestamp type.
func (t *Timestamp) Decode(m map[string]interface{}) error {
	foundMatchType := false
	foundTimestamp := false
	for k, v := range m {
		switch k {
		case "match_type":
			foundMatchType = true

			if mt, ok := v.(string); ok {
				if err := t.MatchType.FromString(mt); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed to decode timestamp matcher, match_type is not a string")
			}
		case "timestamp":
			foundTimestamp = true

			if f, ok := v.(string); ok {
				if _, err := time.Parse(time.Stamp, f); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed to decode timestamp matcher, timestamp is not a string")
			}
		}
	}

	if !(foundTimestamp && foundMatchType) {
		return fmt.Errorf("failed to decode timestamp matcher, missing fields")
	}

	return nil
}

// Encode encodes a Timestamp into the matcher map.
func (t *Timestamp) Encode(out map[string]interface{}) {
	out["match_type"] = t.MatchType.String()
	out["timestamp"] = t.Timestamp.Time.Format(time.Stamp)
}
