package validation

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Validate runs every customer data rule against a raw submission body, so all
// problems appear in one rejection.  submissionDate is read as a calendar day in
// its own location.  A non-nil error means the body could not be parsed at all;
// the caller retries, because the request path only records JSON objects.
func Validate(body []byte, submissionDate time.Time, genders []string) (Customer, []Reason, error) {
	var envelope struct {
		Customer json.RawMessage `json:"customer"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return Customer{}, nil, fmt.Errorf("validation: parse body: %w", err)
	}

	fields := map[string]json.RawMessage{}
	if len(envelope.Customer) > 0 && !isNull(envelope.Customer) {
		if err := json.Unmarshal(envelope.Customer, &fields); err != nil {
			return Customer{}, []Reason{{
				Field:   "customer",
				Code:    CodeInvalid,
				Message: "customer must be a JSON object",
			}}, nil
		}
	}

	year, month, day := submissionDate.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	c := rules{fields: fields}
	customer := Customer{
		GovernmentIDPhoto:      c.photo("government-id-photo"),
		GovernmentIDExpiryDate: c.expiryDate("government-id-expiry-date", today),
		NationalIDNumber:       c.nationalID("national-id-number"),
		FullLegalName:          c.text("full-legal-name"),
		SelfiePhoto:            c.photo("selfie-photo"),
		PhoneNumber:            c.phone("phone-number"),
		Email:                  c.email("email"),
		Gender:                 c.oneOf("gender", genders),
		DateOfBirth:            c.birthDate("date-of-birth", today),
		PlaceOfBirth:           c.text("place-of-birth"),
		MaritalStatus:          c.oneOf("marital-status", maritalStatuses),
		Occupation:             c.oneOf("occupation", occupations),
		WorkDurationYears:      c.smallInt("work-duration-years"),
		MonthlyIncome:          c.income("monthly-income"),
		NumberOfDependents:     c.smallInt("number-of-dependents"),
		DomicileStreet:         c.text("domicile-street"),
		DomicileDistrict:       c.text("domicile-district"),
		DomicileCity:           c.text("domicile-city"),
		DomicileProvince:       c.text("domicile-province"),
	}
	if len(c.reasons) > 0 {
		return Customer{}, c.reasons, nil
	}
	return customer, nil, nil
}

// RetainedSections returns the body reduced to the sections kept after a pass,
// for US-004 and US-005.  The customer section is dropped.
func RetainedSections(body []byte) ([]byte, error) {
	var sections struct {
		Collateral json.RawMessage `json:"collateral,omitempty"`
		Loan       json.RawMessage `json:"loan,omitempty"`
	}
	if err := json.Unmarshal(body, &sections); err != nil {
		return nil, fmt.Errorf("validation: parse body: %w", err)
	}
	retained, err := json.Marshal(sections)
	if err != nil {
		return nil, fmt.Errorf("validation: marshal sections: %w", err)
	}
	return retained, nil
}

// rules checks one field at a time and collects the reasons.  Each method
// returns the parsed value, or the zero value after recording a reason.
type rules struct {
	fields  map[string]json.RawMessage
	reasons []Reason
}

func (c *rules) fail(name, code, message string) {
	c.reasons = append(c.reasons, Reason{Field: "customer." + name, Code: code, Message: message})
}

func (c *rules) missing(name string) {
	c.fail(name, CodeMissing, name+" is required")
}

// present returns the raw value of a mandatory field, recording MISSING when it
// is absent, null, or blank text.
func (c *rules) present(name string) (json.RawMessage, bool) {
	raw, ok := c.fields[name]
	if !ok || isNull(raw) {
		c.missing(name)
		return nil, false
	}
	var s string
	if json.Unmarshal(raw, &s) == nil && strings.TrimSpace(s) == "" {
		c.missing(name)
		return nil, false
	}
	return raw, true
}

// text returns a mandatory text field.
func (c *rules) text(name string) string {
	raw, ok := c.present(name)
	if !ok {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		c.fail(name, CodeInvalid, name+" must be text")
		return ""
	}
	if strings.ContainsRune(s, 0) {
		c.fail(name, CodeInvalid, name+" must not contain NUL characters")
		return ""
	}
	return s
}

// checked returns a mandatory text field that must also satisfy ok.  The second
// result is false when a reason was recorded.
func (c *rules) checked(name, message string, ok func(string) bool) (string, bool) {
	before := len(c.reasons)
	s := c.text(name)
	if len(c.reasons) > before {
		return "", false
	}
	if !ok(s) {
		c.fail(name, CodeInvalid, message)
		return "", false
	}
	return s, true
}

func (c *rules) nationalID(name string) string {
	s, _ := c.checked(name, name+" must contain only letters and digits", func(s string) bool {
		return len(s) <= maxNationalIDLength && nationalIDPattern.MatchString(s)
	})
	return s
}

func (c *rules) phone(name string) string {
	s, _ := c.checked(name, name+" must be an optional + followed by 8 to 15 digits", phonePattern.MatchString)
	return s
}

func (c *rules) oneOf(name string, accepted []string) string {
	s, _ := c.checked(name, name+" is not an accepted value", func(s string) bool {
		return slices.Contains(accepted, s)
	})
	return s
}

// date returns a mandatory field in YYYY-MM-DD form.
func (c *rules) date(name string) (time.Time, bool) {
	s, ok := c.checked(name, name+" must be a real date in YYYY-MM-DD format", func(s string) bool {
		_, err := time.Parse(time.DateOnly, s)
		return err == nil
	})
	if !ok {
		return time.Time{}, false
	}
	parsed, _ := time.Parse(time.DateOnly, s)
	return parsed, true
}

func (c *rules) birthDate(name string, today time.Time) time.Time {
	birth, ok := c.date(name)
	if !ok {
		return time.Time{}
	}
	if birth.After(today) {
		c.fail(name, CodeInvalid, name+" must not be after the submission date")
		return time.Time{}
	}
	// AddDate normalises 29 February to 1 March in non-leap years.
	if birth.AddDate(minimumAgeYears, 0, 0).After(today) {
		c.fail(name, CodeCustomerUnderAge, "customer must be at least 18 years old on the submission date")
		return time.Time{}
	}
	return birth
}

func (c *rules) expiryDate(name string, today time.Time) time.Time {
	expiry, ok := c.date(name)
	if !ok {
		return time.Time{}
	}
	if expiry.Before(today) {
		c.fail(name, CodeGovernmentIDExpired, "government-issued ID is expired on the submission date")
		return time.Time{}
	}
	return expiry
}

// number returns the literal text of a mandatory JSON number.
func (c *rules) number(name, message string) (string, bool) {
	raw, ok := c.present(name)
	if !ok {
		return "", false
	}
	literal := string(bytes.TrimSpace(raw))
	if literal[0] != '-' && (literal[0] < '0' || literal[0] > '9') {
		c.fail(name, CodeInvalid, message)
		return "", false
	}
	return literal, true
}

func (c *rules) smallInt(name string) int16 {
	message := name + " must be a whole number from 0 to 32767"
	literal, ok := c.number(name, message)
	if !ok {
		return 0
	}
	if !wholePattern.MatchString(literal) {
		c.fail(name, CodeInvalid, message)
		return 0
	}
	value, err := strconv.Atoi(literal)
	if err != nil || value > maxSmallInt {
		c.fail(name, CodeInvalid, message)
		return 0
	}
	return int16(value)
}

func (c *rules) income(name string) string {
	message := name + " must be greater than 0, with up to 18 whole digits and 2 decimals"
	literal, ok := c.number(name, message)
	if !ok {
		return ""
	}
	if !incomePattern.MatchString(literal) || strings.Trim(literal, "0.") == "" {
		c.fail(name, CodeInvalid, message)
		return ""
	}
	return literal
}

func (c *rules) photo(name string) Photo {
	message := name + " must be a base64 JPEG or PNG, at most 5 MiB and 50 megapixels"
	s := c.text(name)
	if s == "" {
		return Photo{}
	}
	// The size is checked after decoding because line breaks in the base64 text
	// don't count towards it.  The request body limit caps the work done here.
	content, err := base64.StdEncoding.DecodeString(s)
	if err != nil || len(content) > maxPhotoBytes {
		c.fail(name, CodeInvalid, message)
		return Photo{}
	}
	// DecodeConfig reads only the header, which keeps a small file from expanding
	// into a huge decode; the full decode then catches truncated files.
	config, format, err := image.DecodeConfig(bytes.NewReader(content))
	if err != nil || int64(config.Width)*int64(config.Height) > maxPhotoPixels {
		c.fail(name, CodeInvalid, message)
		return Photo{}
	}
	contentType, accepted := photoTypes[format]
	if !accepted {
		c.fail(name, CodeInvalid, message)
		return Photo{}
	}
	if _, _, err := image.Decode(bytes.NewReader(content)); err != nil {
		c.fail(name, CodeInvalid, message)
		return Photo{}
	}
	return Photo{ContentType: contentType, Content: content}
}

// email returns the optional email.  It is never a rejection reason (FR-003), so
// a value that isn't usable text is treated as absent.
func (c *rules) email(name string) *string {
	raw, ok := c.fields[name]
	if !ok {
		return nil
	}
	var s string
	if json.Unmarshal(raw, &s) != nil || strings.TrimSpace(s) == "" || strings.ContainsRune(s, 0) {
		return nil
	}
	return &s
}

func isNull(raw json.RawMessage) bool {
	return string(bytes.TrimSpace(raw)) == "null"
}
