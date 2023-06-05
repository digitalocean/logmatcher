# Matcher Documentation

The `matcher` package contains the core expression language
for matching `captainslog.SyslogMsg` data in an intuitive and expressive way.

The [Matcher](matcher.go#L9) interface most importantly provides the method
`Matches(captainslog.SyslogMsg) bool`. There are several implementations of
the interface in this package. Those include the following:
* [Value](#value-matcher)
* [Facility](#facility-matcher)
* [Severity](#severity-matcher)
* [KV](#key-value-matcher)
* [UnaryOp](#unary-operator)
* [NAryOp](#n-ary-operator)

From these primitive matchers, essentially any complex expression can be
defined. Matchers are used in three forms, and as such, the configuration for
each form will be described in their respective sections. Technically there are
four commonly used forms, but the fourth (JSON) is simply a special case of
another (YAML), so it will not be discussed. The three forms are:
* Golang
* CLI
* YAML

**Note:** A few of the matchers rely on a special [MatchType](#match-types)
struct described below. Please read through that section first.

## Match Types

**MatchType** is essentially an enumerator that represents the many possible
comparison operators on which some of the matchers rely. Those types are
included here for reference in both their Golang and string-encoded
representations:

| Golang                | Encoded      |
|-----------------------|--------------|
| _**String Types**_    |              |
| ExactMatch            | exact_match  |
| PrefixMatch           | prefix_match |
| Contains              | contains     |
| Regex                 | regex        |
| _**Numeric Types**_   |              |
| LessThan              | lt           |
| LessThanEqual         | lte          |
| GreaterThan           | gt           |
| GreaterThanEqual      | gte          |
| _**Universal Types**_ |              |
| Equals                | equals       |

It should become clearer in the following sections how these types are used.

## Value Matcher

The **Value** matcher is intended to match some basic syslog fields such as
the host, program or content of a message.

### Golang

To instantiate a new value matcher in Golang, you can call:

```golang
func NewValue(t ValueType, m MatchType, v string) *Value
```

Here, the `ValueType` can be one of (in Golang and string-encoded forms):

| Golang  | Encoded |
|---------|---------|
| Program | program |
| Host    | host    |
| Content | content |

**Ex. Usage**
```golang
v := NewValue(Program, PrefixMatch, "logCatcher_")
```

### CLI

To make this package more accessible to CLI tools, an expression language is provided which exposes the 
value matchers as functions. e.g:

```
program(prefix_match, "logCatcher_")
```

### YAML

The YAML encoding of a value matcher is as follows:

```yaml
---
value_matcher:
  type: program
  match_type: prefix_match
  value: 'logCatcher_'
```

## Facility Matcher

The **Facility** matcher is essentially a wrapper around the
`captainslog.Facility` class.

### Golang

To instantiate in Go, simply call:

```golang
func NewFacility(f captainslog.Facility) *Facility
```

**Ex. Usage**

```golang
f := NewFacility(captainslog.Local6)
```

### CLI

A convenience function is also supplied in the CLI form:

```
'facility("local6")'
```

### YAML

And in YAML:

```yaml
---
facility_matcher:
  facility: local6
```

## Severity Matcher

The **Severity** matcher is very similar to the Facility matcher, with the
added ability to compare using the equals, less than, or greater than
operators.

### Golang

To instantiate in Go, simply call:

```golang
func NewSeverity(t MatchType, s captainslog.Severity) *Severity
```

**Ex. Usage**

```golang
s := NewSeverity(LessThan, captainslog.Warn)
```

### CLI

A convenience function is also supplied in the CLI form:

```
severity(lt, "warn")
```

### YAML

And in YAML:

```yaml
---
severity_matcher:
  match_type: lt
  severity: warn
```

## Key-Value Matcher

The **KV** matcher is where things start to get a little more interesting.
This structure provides a flexible key-value matcher for arbitrary JSON
data in a syslog message.

### Golang

To instantiate in Go, call:

```golang
func NewKV(key string, m MatchType, value interface{}) *KV
```

**Ex. Usage**

```golang
kv := NewKV("response.code", LessThan, 300)
```

This matcher would match a log whose JSON content contains, e.g.:
```json
{"response":{"code":204}}
```

But not, e.g.:
```json
{"response":{"code":401}}
```

Periods in the supplied `key` are used to dereference multi-level JSON data.
The supplied `value` can be anything, but the expectation from the library
is that the user supplies proper data types and their corresponding match
types.

### CLI

A convenience function is also supplied in the CLI form:

```
kv("response.code", lt, 300)
```

### YAML

Unfortunately, the ability to keep things totally generic breaks down when
representing this data in YAML. When supplying numeric rules, you must use the
`num_value` field in the definition, e.g.:

```yaml
---
kv_matcher:
  key: 'response.code'
  match_type: lt
  num_value: 300
```

The following schema holds for the YAML:
```yaml
---
kv_matcher:
  key: <string>
  match_type: <string>
  # One and only one of the following is required
  num_value: <float>
  str_value: <string>
  bool_value: <true or false>
```

## Unary Operator

The **UnaryOp** matcher is a generic matcher that allows for a unary
operator to apply to another matcher. In reality, the only supported
operator is the `not` operation.

### Golang

To instantiate in Go, simply call:

```golang
func NewUnaryOp(t UnaryOpType, m Matcher) *UnaryOp
```

**Ex. Usage**

```golang
o := NewUnaryOp(Not, someMatcher)
```

### CLI

A convenience function is also supplied in the CLI form:

**Warning:** This is a contrived example and should never be added in real
life!
```
not(host(prefix_match, "prod-somehost"))
```

### YAML

And in YAML:

```yaml
---
unary_op:
  type: not
  matcher:
    value_matcher:
      type: host
      match_type: prefix_match
      value: 'prod-somehost'
```

## N-Ary Operator

The **NAryOp** matcher is a generic matcher that allows for an n-ary
operator to apply to another matcher. In reality, the only supported
operators are the `and` and `or` operations.

### Golang

To instantiate in Go, simply call:

```golang
func NewNAryOp(t NAryOpType, m ...Matcher) *NAryOp
```

**Ex. Usage**

```golang
o := NewNAryOp(And, someMatcher, someOtherMatcher)
```

### CLI

The n-ary operators are more readable in the CLI via the natural forms with
parenthetical precedence rules, e.g.:
```
X and Y and Z
X and Y or Z
X and (Y or Z)
```

So we may see, e.g.

```
'program(prefix_match, "logCatcher_") and \
>  not(program(exact_match, "logCatcher_staging"))'
```

### YAML

And in YAML:

```yaml
---
n_ary_op:
  type: and
  matchers:
  - value_matcher:
      type: program
      match_type: prefix_match
      value: 'logCatcher_'
  - unary_op:
      type: not
      matcher:
        value_matcher:
          type: program
          match_type: exact_match
          value: 'logCatcher_staging'
```

## License

The project is licensed under the Apache License, Version 2.0.

You can find a copy of the license in the [LICENSE](LICENSE) file or visit the [Apache website](http://www.apache.org/licenses/LICENSE-2.0) for more details.

