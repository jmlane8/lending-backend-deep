// Package validation holds the customer data rules that run after a submission
// is acknowledged.  It is pure: all input is passed in and it performs no I/O.
//
// Reason messages must never contain submitted values (FR-011).
package validation

import (
	_ "image/jpeg" // register the JPEG decoder
	_ "image/png"  // register the PNG decoder
	"regexp"
	"time"
)

const (
	maxPhotoBytes  = 5 << 20
	maxPhotoPixels = 50_000_000
	maxSmallInt    = 32767

	// maxNationalIDLength keeps the value inside the size limit of the unique
	// index on customer.national_id_number; a longer value could never be stored.
	maxNationalIDLength = 255

	minimumAgeYears = 18
)

var (
	nationalIDPattern = regexp.MustCompile(`^[A-Za-z0-9]+$`)
	phonePattern      = regexp.MustCompile(`^\+?[0-9]{8,15}$`)
	wholePattern      = regexp.MustCompile(`^[0-9]+$`)
	incomePattern     = regexp.MustCompile(`^[0-9]{1,18}(\.[0-9]{1,2})?$`)

	maritalStatuses = []string{"SINGLE", "MARRIED", "DIVORCED", "WIDOWED"}
	occupations     = []string{"EMPLOYEE", "ENTREPRENEUR", "CIVIL_SERVANT"}
	photoTypes      = map[string]string{"jpeg": "image/jpeg", "png": "image/png"}
)

// Reason codes reported to the partner.
const (
	CodeMissing             = "MISSING"
	CodeInvalid             = "INVALID"
	CodeCustomerUnderAge    = "CUSTOMER_UNDER_AGE"
	CodeGovernmentIDExpired = "GOVERNMENT_ID_EXPIRED"
)

// Reason describes one problem with a submitted field.
type Reason struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Photo is a decoded image.
type Photo struct {
	ContentType string // "image/jpeg" | "image/png"
	Content     []byte
}

// Customer holds the parsed customer data of a submission that passed every rule.
type Customer struct {
	GovernmentIDPhoto      Photo
	GovernmentIDExpiryDate time.Time
	NationalIDNumber       string
	FullLegalName          string
	SelfiePhoto            Photo
	PhoneNumber            string
	Email                  *string
	Gender                 string
	DateOfBirth            time.Time
	PlaceOfBirth           string
	MaritalStatus          string
	Occupation             string
	WorkDurationYears      int16
	MonthlyIncome          string // decimal text; kept as text to avoid float rounding
	NumberOfDependents     int16
	DomicileStreet         string
	DomicileDistrict       string
	DomicileCity           string
	DomicileProvince       string
}
