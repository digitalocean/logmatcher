package matcher

import "fmt"

// MatchType is the enum class for representing different match types.
type MatchType int

// Match types. New MatchType values MUST be appended to the end of the list for
// database consistency.
const (
	// String types
	ExactMatch MatchType = iota
	PrefixMatch
	Contains
	Regex

	// Numeric types
	LessThan
	LessThanEqual
	GreaterThan
	GreaterThanEqual

	// Universal types
	Equals
)

// String converts a MatchType to its corresponding string representation.
func (m MatchType) String() string {
	switch m {
	case ExactMatch:
		return "exact_match"
	case PrefixMatch:
		return "prefix_match"
	case Contains:
		return "contains"
	case Regex:
		return "regex"
	case Equals:
		return "equals"
	case LessThan:
		return "lt"
	case LessThanEqual:
		return "lte"
	case GreaterThan:
		return "gt"
	case GreaterThanEqual:
		return "gte"
	default:
		return "invalid type"
	}
}

// FromString converts the MatchType to the value corresponding to the supplied
// string representation.
func (m *MatchType) FromString(s string) error {
	switch s {
	case "exact_match":
		*m = ExactMatch
	case "prefix_match":
		*m = PrefixMatch
	case "contains":
		*m = Contains
	case "regex":
		*m = Regex
	case "equals":
		*m = Equals
	case "lt":
		*m = LessThan
	case "lte":
		*m = LessThanEqual
	case "gt":
		*m = GreaterThan
	case "gte":
		*m = GreaterThanEqual
	default:
		return fmt.Errorf("failed to convert string to MatchType")
	}

	return nil
}
