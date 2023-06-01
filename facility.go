package matcher

import (
	"fmt"

	"github.com/digitalocean/captainslog"
)

// Facility represents a syslog facility matcher
type Facility struct {
	Facility captainslog.Facility
}

// NewFacility returns a new Facility with the specified value.
func NewFacility(f captainslog.Facility) *Facility {
	return &Facility{
		Facility: f,
	}
}

// String converts a Facility to its corresponding string representation.
func (f Facility) String() string {
	return fmt.Sprintf("facility(%s)", f.Facility)
}

// Matches returns true if the Facility matches the supplied SyslogMsg.
func (f *Facility) Matches(m captainslog.SyslogMsg) bool {
	return f.Facility == m.Pri.Facility
}

// Decode is a helper function to decode a facility matcher.
func (f *Facility) Decode(m map[string]interface{}) error {
	foundFacility := false
	for k, v := range m {
		switch k {
		case "facility":
			foundFacility = true

			if t, ok := v.(string); ok {
				if err := f.Facility.FromString(t); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed to decode facility matcher, facility is not a string")
			}
		}
	}

	if !foundFacility {
		return fmt.Errorf("failed to decode facility matcher, missing fields")
	}

	return nil
}

// Encode is a helper function to encode a facility into a matcher map.
func (f *Facility) Encode(out map[string]interface{}) {
	out["facility"] = f.Facility.String()
}
