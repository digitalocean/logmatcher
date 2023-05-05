package matcher

import (
	"fmt"

	"github.com/digitalocean/captainslog"
)

// Severity represents a syslog facility matcher
type Severity struct {
	MatchType MatchType
	Severity  captainslog.Severity
}

// NewSeverity returns a new Severity with the specified value.
func NewSeverity(t MatchType, f captainslog.Severity) *Severity {
	return &Severity{
		MatchType: t,
		Severity:  f,
	}
}

// String converts a Severity to its corresponding string representation.
func (s Severity) String() string {
	return fmt.Sprintf("severity(%s, %s)", s.MatchType.String(), s.Severity)
}

// Matches returns true if the Severity matches the supplied SyslogMsg.
func (s *Severity) Matches(m captainslog.SyslogMsg) bool {
	// Syslog severity values are lower for higher severities.
	switch s.MatchType {
	case Equals:
		return m.Pri.Severity == s.Severity
	case LessThan:
		return m.Pri.Severity > s.Severity
	case LessThanEqual:
		return m.Pri.Severity >= s.Severity
	case GreaterThan:
		return m.Pri.Severity < s.Severity
	case GreaterThanEqual:
		return m.Pri.Severity <= s.Severity
	default:
		return false
	}
}

// Decode decodes a matcher map into a Severity type.
func (s *Severity) Decode(m map[string]interface{}) error {
	foundMatchType := false
	foundSeverity := false
	for k, v := range m {
		switch k {
		case "match_type":
			foundMatchType = true

			if mt, ok := v.(string); ok {
				if err := s.MatchType.FromString(mt); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed to decode severity matcher, match_type is not a string")
			}
		case "severity":
			foundSeverity = true

			if f, ok := v.(string); ok {
				if err := s.Severity.FromString(f); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed to decode severity matcher, severity is not a string")
			}
		}
	}

	if !(foundSeverity && foundMatchType) {
		return fmt.Errorf("failed to decode severity matcher, missing fields")
	}

	return nil
}

// Encode encodes a Severity into the matcher map.
func (s *Severity) Encode(out map[string]interface{}) {
	out["match_type"] = s.MatchType.String()
	out["severity"] = s.Severity.String()
}
