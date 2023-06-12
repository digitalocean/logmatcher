package matcher

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/digitalocean/captainslog"
)

func TestNewValue(t *testing.T) {
	v := NewValue(Program, PrefixMatch, "topo")
	if v == nil {
		t.Errorf("got nil ValueMatcher")
	}
}

func TestNewUnaryOp(t *testing.T) {
	o := NewUnaryOp(Not, NewValue(Program, Contains, "bre"))
	if o == nil {
		t.Errorf("got nil UnaryOp")
	}
}

func TestNewNAryOp(t *testing.T) {
	o := NewNAryOp(And,
		NewHostname(Contains, "bre"),
		NewUnaryOp(Not,
			NewHostname(Regex, "bad-host.*nyc3.internal.digitalocean.com")))
	if o == nil {
		t.Errorf("got nil UnaryOp")
	}
}

func TestNot(t *testing.T) {
	o := NewUnaryOp(Not, NewValue(Program, ExactMatch, "foo"))

	m := captainslog.NewSyslogMsg()
	m.SetHost("")

	m.SetProgram("bar")
	if want, got := true, o.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.SetProgram("foo")
	if want, got := false, o.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
}

func TestAnd(t *testing.T) {
	o := NewNAryOp(And,
		NewHostname(ExactMatch, "foo"),
		NewValue(Program, ExactMatch, "bar"))

	m := captainslog.NewSyslogMsg()

	m.SetHost("foo")
	m.SetProgram("bar")
	if want, got := true, o.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
	m.SetProgram("baz")
	if want, got := false, o.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
	m.SetHost("for")
	m.SetProgram("bar")
	if want, got := false, o.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
	m.SetProgram("baz")
	if want, got := false, o.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
}

func TestOr(t *testing.T) {
	o := NewNAryOp(Or,
		NewHostname(ExactMatch, "foo"),
		NewValue(Program, ExactMatch, "bar"))

	m := captainslog.NewSyslogMsg()

	m.SetHost("foo")
	m.SetProgram("bar")
	if want, got := true, o.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
	m.SetProgram("baz")
	if want, got := true, o.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
	m.SetHost("for")
	m.SetProgram("bar")
	if want, got := true, o.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
	m.SetProgram("baz")
	if want, got := false, o.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
}

func TestKVMatcherString(t *testing.T) {
	kvmExact := NewKV("foo", ExactMatch, "bar")
	kvmPrefix := NewKV("foo", PrefixMatch, "ba")
	kvmContains := NewKV("foo", Contains, "a")
	kvmRegex := NewKV("foo", Regex, "b.r")

	m := captainslog.NewSyslogMsg()
	m.IsJSON = false

	if want, got := false, kvmExact.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.IsJSON = true
	m.JSONValues = make(map[string]interface{})

	if want, got := false, kvmExact.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = "bar"
	if want, got := true, kvmExact.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = "bar"
	if want, got := true, kvmExact.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = "bar"
	if want, got := true, kvmPrefix.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = "bar"
	if want, got := true, kvmContains.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = "bar"
	if want, got := true, kvmRegex.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
}

func TestKVMatcherInt(t *testing.T) {
	kvmEq := NewKV("foo", Equals, 10)
	kvmLT := NewKV("foo", LessThan, 10)
	kvmLTE := NewKV("foo", LessThanEqual, 10)
	kvmGT := NewKV("foo", GreaterThan, 10)
	kvmGTE := NewKV("foo", GreaterThanEqual, 10)

	m := captainslog.NewSyslogMsg()
	m.IsJSON = true
	m.JSONValues = make(map[string]interface{})

	m.JSONValues["foo"] = float64(10)
	if want, got := true, kvmEq.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = float64(9)
	if want, got := true, kvmLT.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = float64(10)
	if want, got := false, kvmLT.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = float64(11)
	if want, got := false, kvmLT.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = float64(10)
	if want, got := true, kvmLTE.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = float64(11)
	if want, got := true, kvmGT.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = float64(10)
	if want, got := true, kvmGTE.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	// String field, int matcher
	m.JSONValues["foo"] = "bar"
	if want, got := false, kvmEq.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
}

func TestKVMatcherFloat(t *testing.T) {
	kvmEq := NewKV("foo", Equals, 10.473)
	kvmLT := NewKV("foo", LessThan, 10.473)
	kvmLTE := NewKV("foo", LessThanEqual, 10.473)
	kvmGT := NewKV("foo", GreaterThan, 10.473)
	kvmGTE := NewKV("foo", GreaterThanEqual, 10.473)

	m := captainslog.NewSyslogMsg()
	m.IsJSON = true
	m.JSONValues = make(map[string]interface{})

	m.JSONValues["foo"] = 10.473
	if want, got := true, kvmEq.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = 10.472
	if want, got := true, kvmLT.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = 10.473
	if want, got := true, kvmLTE.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = 10.474
	if want, got := true, kvmGT.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = 10.472
	if want, got := false, kvmGTE.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = 10.473
	if want, got := true, kvmGTE.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = 10.474
	if want, got := true, kvmGTE.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
}

func TestKVMatcherBool(t *testing.T) {
	kvm := NewKV("foo", Equals, false)

	m := captainslog.NewSyslogMsg()
	m.IsJSON = true

	m.JSONValues["foo"] = false
	if want, got := true, kvm.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues["foo"] = true
	if want, got := false, kvm.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
}

func clearMap(m map[string]interface{}) {
	for k := range m {
		delete(m, k)
	}
}

func TestKVMatcherDeep(t *testing.T) {
	kvm := NewKV("foo.bar.baz", Equals, false)

	m := captainslog.NewSyslogMsg()
	m.IsJSON = true
	m.JSONValues = make(map[string]interface{})

	if want, got := false, kvm.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	jsonStr := `{
		"foo": {
			"bar": {
				"baz": "airforce"
			}
		}
	}`
	clearMap(m.JSONValues)
	_ = json.Unmarshal([]byte(jsonStr), &m.JSONValues)

	if want, got := false, kvm.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	jsonStr = `{
		"foo": {
			"bar": {
				"baz": true
			}
		}
	}`
	clearMap(m.JSONValues)
	_ = json.Unmarshal([]byte(jsonStr), &m.JSONValues)
	if want, got := false, kvm.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	jsonStr = `{
		"foo": {
			"bar": {
				"baz": 7.42
			}
		}
	}`
	clearMap(m.JSONValues)
	_ = json.Unmarshal([]byte(jsonStr), &m.JSONValues)
	if want, got := false, kvm.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	jsonStr = `{
		"foo": {
			"baz": {
				"bar": false
			}
		}
	}`
	clearMap(m.JSONValues)
	_ = json.Unmarshal([]byte(jsonStr), &m.JSONValues)
	if want, got := false, kvm.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	jsonStr = `{
		"foo": {
			"bar": {
				"baz": false
			}
		}
	}`
	clearMap(m.JSONValues)
	_ = json.Unmarshal([]byte(jsonStr), &m.JSONValues)
	if want, got := true, kvm.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	// This isn't valid JSON
	jsonStr = `{
		"foo": {
			10: {
				"baz": false
			}
		}
	}`
	clearMap(m.JSONValues)
	_ = json.Unmarshal([]byte(jsonStr), &m.JSONValues)
	if want, got := false, kvm.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.JSONValues = make(map[string]interface{})
	m.JSONValues["foo"] = make(map[int]float64)
	if want, got := false, kvm.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
}

func TestFacilityMatcher(t *testing.T) {
	e := NewFacility(captainslog.Kern)
	m := captainslog.NewSyslogMsg()

	_ = m.SetFacility(captainslog.Kern)
	if want, got := true, e.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	_ = m.SetFacility(captainslog.UUCP)
	if want, got := false, e.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
}

func TestSeverityMatcher(t *testing.T) {
	e := NewSeverity(Equals, captainslog.Debug)
	m := captainslog.NewSyslogMsg()

	_ = m.SetSeverity(captainslog.Debug)
	if want, got := true, e.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	_ = m.SetSeverity(captainslog.Info)
	if want, got := false, e.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	e = NewSeverity(LessThan, captainslog.Warning)

	_ = m.SetSeverity(captainslog.Debug)
	if want, got := true, e.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	_ = m.SetSeverity(captainslog.Info)
	if want, got := true, e.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	_ = m.SetSeverity(captainslog.Notice)
	if want, got := true, e.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	_ = m.SetSeverity(captainslog.Warning)
	if want, got := false, e.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	_ = m.SetSeverity(captainslog.Err)
	if want, got := false, e.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	_ = m.SetSeverity(captainslog.Crit)
	if want, got := false, e.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
}

func TestTimestampMatcher(t *testing.T) {
	e := NewTimestamp(Equals, captainslog.Time{
		Time:       time.Now(),
		TimeFormat: time.Stamp,
	})

	m := captainslog.NewSyslogMsg()
	m.SetTime(time.Now())

	if want, got := false, e.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.SetTime(time.Date(1969, time.April, 20, 0, 0, 0, 0, time.UTC))
	e.MatchType = LessThan

	if want, got := false, e.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v, message time is %v, compared with %v", want, got, m.Time.Format(time.Stamp), e.Timestamp.Time.Format(time.Stamp))
	}

	e.MatchType = GreaterThan

	if want, got := true, e.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v, time is %v", want, got, m.Time)
	}

}

func TestHostnameMatcher(t *testing.T) {
	hostExact := NewHostname(ExactMatch, "cool.website.com:5757")
	hostPrefix := NewHostname(PrefixMatch, "coo")
	hostContains := NewHostname(Contains, "website")
	hostRegex := NewHostname(Regex, "c.o")

	m := captainslog.NewSyslogMsg()
	m.Host = "cool.website.com:5757"

	if want, got := true, hostExact.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
	m.Host = "foo"
	if want, got := false, hostExact.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.Host = "cool.website.com:5757"
	if want, got := true, hostPrefix.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
	m.Host = "bar"
	if want, got := false, hostPrefix.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.Host = "cool.website.com:5757"
	if want, got := true, hostContains.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
	m.Host = "foo"
	if want, got := false, hostContains.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}

	m.Host = "cool.website.com:5757"
	if want, got := true, hostRegex.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
	m.Host = "bar"
	if want, got := false, hostRegex.Matches(m); want != got {
		t.Errorf("want != got, want = %v, got = %v", want, got)
	}
}
