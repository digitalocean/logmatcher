package matcher

import (
	"fmt"
	"github.com/digitalocean/captainslog"
	"regexp"
	"strings"
)

// Hostname represents a syslog hostname matcher
type Hostname struct {
	MatchType   MatchType
	NameMatcher string
}

// NewHostname returns a new hostname matcher
func NewHostname(m MatchType, n string) *Hostname {
	return &Hostname{
		MatchType:   m,
		NameMatcher: n,
	}
}

// String converts a Hostname matcher to its string representation
func (h Hostname) String() string {
	return fmt.Sprintf("hostname(%s, %s)", h.MatchType.String(), h.NameMatcher)
}

// Matches returns true if the Hostname aligns with the supplied SyslogMsg hostname and MatchType
func (h *Hostname) Matches(m captainslog.SyslogMsg) bool {
	switch h.MatchType {
	case ExactMatch:
		return h.NameMatcher == m.Host
	case PrefixMatch:
		return strings.HasPrefix(m.Host, h.NameMatcher)
	case Contains:
		return strings.Contains(m.Host, h.NameMatcher)
	case Regex:
		matched, _ := regexp.MatchString(h.NameMatcher, m.Host)
		return matched
	default:
		return false
	}
}

// Decode decodes a matcher map into a Hostname type.
func (h *Hostname) Decode(m map[string]interface{}) error {
	foundMatchType := false
	hostIsString := false
	for k, v := range m {
		switch k {
		case "match_type":
			foundMatchType = true

			if mt, ok := v.(string); ok {
				if err := h.MatchType.FromString(mt); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed to decode hostname matcher, match_type is not a string")
			}
		case "hostname":
			hostIsString = true

			if _, ok := v.(string); ok {
			} else {
				return fmt.Errorf("failed to decode hostname matcher, hostname is not a string")
			}
		}
	}

	if !(foundMatchType && hostIsString) {
		return fmt.Errorf("failed to decode hostname matcher, missing fields")
	}

	return nil
}

// Encode encodes a Hostname into a matcher map.
func (h *Hostname) Encode(out map[string]interface{}) {
	out["match_type"] = h.MatchType.String()
	out["hostname"] = fmt.Sprintf(h.NameMatcher)
}
