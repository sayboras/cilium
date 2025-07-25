//go:build !disable_pgv
// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/type/http/v3/cookie.proto

package httpv3

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on Cookie with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Cookie) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Cookie with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in CookieMultiError, or nil if none found.
func (m *Cookie) ValidateAll() error {
	return m.validate(true)
}

func (m *Cookie) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetName()) < 1 {
		err := CookieValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if d := m.GetTtl(); d != nil {
		dur, err := d.AsDuration(), d.CheckValid()
		if err != nil {
			err = CookieValidationError{
				field:  "Ttl",
				reason: "value is not a valid duration",
				cause:  err,
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		} else {

			gte := time.Duration(0*time.Second + 0*time.Nanosecond)

			if dur < gte {
				err := CookieValidationError{
					field:  "Ttl",
					reason: "value must be greater than or equal to 0s",
				}
				if !all {
					return err
				}
				errors = append(errors, err)
			}

		}
	}

	// no validation rules for Path

	for idx, item := range m.GetAttributes() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CookieValidationError{
						field:  fmt.Sprintf("Attributes[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CookieValidationError{
						field:  fmt.Sprintf("Attributes[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CookieValidationError{
					field:  fmt.Sprintf("Attributes[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return CookieMultiError(errors)
	}

	return nil
}

// CookieMultiError is an error wrapping multiple validation errors returned by
// Cookie.ValidateAll() if the designated constraints aren't met.
type CookieMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CookieMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CookieMultiError) AllErrors() []error { return m }

// CookieValidationError is the validation error returned by Cookie.Validate if
// the designated constraints aren't met.
type CookieValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CookieValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CookieValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CookieValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CookieValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CookieValidationError) ErrorName() string { return "CookieValidationError" }

// Error satisfies the builtin error interface
func (e CookieValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCookie.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CookieValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CookieValidationError{}

// Validate checks the field values on CookieAttribute with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *CookieAttribute) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CookieAttribute with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CookieAttributeMultiError, or nil if none found.
func (m *CookieAttribute) ValidateAll() error {
	return m.validate(true)
}

func (m *CookieAttribute) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetName()) < 1 {
		err := CookieAttributeValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetName()) > 16384 {
		err := CookieAttributeValidationError{
			field:  "Name",
			reason: "value length must be at most 16384 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_CookieAttribute_Name_Pattern.MatchString(m.GetName()) {
		err := CookieAttributeValidationError{
			field:  "Name",
			reason: "value does not match regex pattern \"^[^\\x00\\n\\r]*$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetValue()) > 16384 {
		err := CookieAttributeValidationError{
			field:  "Value",
			reason: "value length must be at most 16384 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_CookieAttribute_Value_Pattern.MatchString(m.GetValue()) {
		err := CookieAttributeValidationError{
			field:  "Value",
			reason: "value does not match regex pattern \"^[^\\x00\\n\\r]*$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return CookieAttributeMultiError(errors)
	}

	return nil
}

// CookieAttributeMultiError is an error wrapping multiple validation errors
// returned by CookieAttribute.ValidateAll() if the designated constraints
// aren't met.
type CookieAttributeMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CookieAttributeMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CookieAttributeMultiError) AllErrors() []error { return m }

// CookieAttributeValidationError is the validation error returned by
// CookieAttribute.Validate if the designated constraints aren't met.
type CookieAttributeValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CookieAttributeValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CookieAttributeValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CookieAttributeValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CookieAttributeValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CookieAttributeValidationError) ErrorName() string { return "CookieAttributeValidationError" }

// Error satisfies the builtin error interface
func (e CookieAttributeValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCookieAttribute.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CookieAttributeValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CookieAttributeValidationError{}

var _CookieAttribute_Name_Pattern = regexp.MustCompile("^[^\x00\n\r]*$")

var _CookieAttribute_Value_Pattern = regexp.MustCompile("^[^\x00\n\r]*$")
