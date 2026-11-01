package validation

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"hash/crc32"
	"image"
	"image/jpeg"
	"image/png"
	"maps"
	"strings"
	"testing"
	"time"
)

var (
	submissionDate = time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC)
	genders        = []string{"MALE", "FEMALE"}
)

func jpegBase64(t *testing.T) string {
	t.Helper()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 4, 4)), nil); err != nil {
		t.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func pngBytes(t *testing.T, width, height int) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, image.NewRGBA(image.Rect(0, 0, width, height))); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func validCustomer(t *testing.T) map[string]any {
	t.Helper()
	photo := jpegBase64(t)
	return map[string]any{
		"government-id-photo":       photo,
		"government-id-expiry-date": "2030-05-14",
		"national-id-number":        "3201012345678901",
		"full-legal-name":           "Alex Morgan",
		"selfie-photo":              photo,
		"phone-number":              "+15551234567",
		"email":                     "alex@example.com",
		"gender":                    "FEMALE",
		"date-of-birth":             "1990-02-11",
		"place-of-birth":            "Springfield",
		"marital-status":            "SINGLE",
		"occupation":                "EMPLOYEE",
		"work-duration-years":       json.Number("6"),
		"monthly-income":            json.Number("4500.50"),
		"number-of-dependents":      json.Number("1"),
		"domicile-street":           "12 Example Street",
		"domicile-district":         "Central",
		"domicile-city":             "Springfield",
		"domicile-province":         "North",
	}
}

func mandatoryFields(t *testing.T) []string {
	t.Helper()
	fields := make([]string, 0, 18)
	for name := range maps.Keys(validCustomer(t)) {
		if name != "email" {
			fields = append(fields, name)
		}
	}
	return fields
}

func run(t *testing.T, customer map[string]any) (Customer, []Reason) {
	t.Helper()
	body, err := json.Marshal(map[string]any{"customer": customer})
	if err != nil {
		t.Fatal(err)
	}
	c, reasons, err := Validate(body, submissionDate, genders)
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	return c, reasons
}

func codes(reasons []Reason) map[string]string {
	out := make(map[string]string, len(reasons))
	for _, r := range reasons {
		out[strings.TrimPrefix(r.Field, "customer.")] = r.Code
	}
	return out
}

func TestValidCustomerPasses(t *testing.T) {
	c, reasons := run(t, validCustomer(t))
	if len(reasons) != 0 {
		t.Fatalf("unexpected reasons: %+v", reasons)
	}
	if c.NationalIDNumber != "3201012345678901" || c.MonthlyIncome != "4500.50" || c.WorkDurationYears != 6 {
		t.Errorf("fields not parsed: %+v", c)
	}
	if c.GovernmentIDPhoto.ContentType != "image/jpeg" || len(c.SelfiePhoto.Content) == 0 {
		t.Errorf("photos not decoded: %+v", c.GovernmentIDPhoto.ContentType)
	}
	if c.Email == nil || *c.Email != "alex@example.com" {
		t.Errorf("email not parsed: %v", c.Email)
	}
}

func TestMissingMandatoryField(t *testing.T) {
	blanks := map[string]func(map[string]any, string){
		"absent":     func(c map[string]any, f string) { delete(c, f) },
		"null":       func(c map[string]any, f string) { c[f] = nil },
		"empty":      func(c map[string]any, f string) { c[f] = "" },
		"whitespace": func(c map[string]any, f string) { c[f] = "  \t " },
	}
	for _, field := range mandatoryFields(t) {
		for name, blank := range blanks {
			t.Run(field+"/"+name, func(t *testing.T) {
				customer := validCustomer(t)
				blank(customer, field)
				_, reasons := run(t, customer)
				got := codes(reasons)
				if len(got) != 1 || got[field] != CodeMissing {
					t.Errorf("reasons = %+v, want only %s MISSING", reasons, field)
				}
			})
		}
	}
}

func TestSeveralProblemsReportedTogether(t *testing.T) {
	customer := validCustomer(t)
	delete(customer, "gender")
	delete(customer, "phone-number")
	customer["marital-status"] = "ENGAGED"
	customer["occupation"] = "STUDENT"

	got := codes(func() []Reason { _, r := run(t, customer); return r }())
	want := map[string]string{
		"gender":         CodeMissing,
		"phone-number":   CodeMissing,
		"marital-status": CodeInvalid,
		"occupation":     CodeInvalid,
	}
	if !maps.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCustomerObjectProblems(t *testing.T) {
	t.Run("absent", func(t *testing.T) {
		_, reasons, err := Validate([]byte(`{}`), submissionDate, genders)
		if err != nil {
			t.Fatal(err)
		}
		if len(reasons) != len(mandatoryFields(t)) {
			t.Errorf("got %d reasons", len(reasons))
		}
	})
	t.Run("wrong type", func(t *testing.T) {
		_, reasons, err := Validate([]byte(`{"customer":5}`), submissionDate, genders)
		if err != nil {
			t.Fatal(err)
		}
		if len(reasons) != 1 || reasons[0].Field != "customer" || reasons[0].Code != CodeInvalid {
			t.Errorf("reasons = %+v", reasons)
		}
	})
}

func TestInvalidValues(t *testing.T) {
	tests := []struct {
		field string
		value any
	}{
		{"marital-status", "single"},
		{"marital-status", "Single"},
		{"occupation", "employee"},
		{"national-id-number", "12-34"},
		{"national-id-number", strings.Repeat("1", 256)},
		{"phone-number", "1234567"},
		{"phone-number", "1234567890123456"},
		{"phone-number", "+1 555 123 4567"},
		{"gender", "female"},
		{"gender", "OTHER"},
		{"date-of-birth", "1990-02-30"},
		{"date-of-birth", "11/02/1990"},
		{"date-of-birth", "2026-11-01"},
		{"government-id-expiry-date", "2030-13-01"},
		{"work-duration-years", json.Number("-1")},
		{"work-duration-years", json.Number("1.5")},
		{"work-duration-years", json.Number("32768")},
		{"number-of-dependents", json.Number("1e1")},
		{"monthly-income", json.Number("0")},
		{"monthly-income", json.Number("-5")},
		{"monthly-income", json.Number("1.234")},
		{"monthly-income", json.Number("1e3")},
		{"monthly-income", json.Number("1234567890123456789")},
		{"government-id-photo", "not base64!"},
		{"selfie-photo", base64.StdEncoding.EncodeToString([]byte("not an image"))},
		// Wrong JSON types.
		{"monthly-income", "lots"},
		{"work-duration-years", "6"},
		{"full-legal-name", json.Number("7")},
		{"gender", true},
		{"phone-number", []string{"1"}},
		// NUL characters.
		{"full-legal-name", "Alex\x00Morgan"},
		{"domicile-city", "\x00"},
	}
	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			customer := validCustomer(t)
			customer[tt.field] = tt.value
			_, reasons := run(t, customer)
			got := codes(reasons)
			if len(got) != 1 || got[tt.field] != CodeInvalid {
				t.Errorf("reasons = %+v, want only %s INVALID", reasons, tt.field)
			}
		})
	}
}

func TestReasonsNeverContainSubmittedValues(t *testing.T) {
	const secret = "SECRETVALUE"
	customer := validCustomer(t)
	for _, field := range []string{"phone-number", "gender", "marital-status", "occupation",
		"national-id-number", "date-of-birth", "government-id-expiry-date", "government-id-photo",
		"selfie-photo", "work-duration-years", "monthly-income", "number-of-dependents"} {
		customer[field] = secret + "!"
	}
	_, reasons := run(t, customer)
	if len(reasons) == 0 {
		t.Fatal("expected reasons")
	}
	for _, r := range reasons {
		if strings.Contains(r.Message, secret) {
			t.Errorf("%s message leaks the value: %q", r.Field, r.Message)
		}
	}
}

func TestPhotos(t *testing.T) {
	t.Run("png accepted", func(t *testing.T) {
		customer := validCustomer(t)
		customer["selfie-photo"] = base64.StdEncoding.EncodeToString(pngBytes(t, 2, 2))
		c, reasons := run(t, customer)
		if len(reasons) != 0 || c.SelfiePhoto.ContentType != "image/png" {
			t.Errorf("reasons = %+v, type = %q", reasons, c.SelfiePhoto.ContentType)
		}
	})
	t.Run("truncated image", func(t *testing.T) {
		raw := pngBytes(t, 32, 32)
		customer := validCustomer(t)
		customer["selfie-photo"] = base64.StdEncoding.EncodeToString(raw[:len(raw)/2])
		_, reasons := run(t, customer)
		if got := codes(reasons); got["selfie-photo"] != CodeInvalid {
			t.Errorf("reasons = %+v", reasons)
		}
	})
	t.Run("line-wrapped base64 at exactly 5 MiB", func(t *testing.T) {
		raw := append(pngBytes(t, 2, 2), make([]byte, 5<<20)...)[:5<<20]
		encoded := base64.StdEncoding.EncodeToString(raw)
		var wrapped strings.Builder
		for len(encoded) > 76 {
			wrapped.WriteString(encoded[:76] + "\r\n")
			encoded = encoded[76:]
		}
		wrapped.WriteString(encoded)

		customer := validCustomer(t)
		customer["selfie-photo"] = wrapped.String()
		c, reasons := run(t, customer)
		if len(reasons) != 0 || len(c.SelfiePhoto.Content) != 5<<20 {
			t.Errorf("reasons = %+v, size = %d", reasons, len(c.SelfiePhoto.Content))
		}
	})
	t.Run("over 5 MiB", func(t *testing.T) {
		customer := validCustomer(t)
		customer["selfie-photo"] = base64.StdEncoding.EncodeToString(make([]byte, 5<<20+1))
		_, reasons := run(t, customer)
		if got := codes(reasons); got["selfie-photo"] != CodeInvalid {
			t.Errorf("reasons = %+v", reasons)
		}
	})
	t.Run("over 50 megapixels", func(t *testing.T) {
		customer := validCustomer(t)
		customer["selfie-photo"] = base64.StdEncoding.EncodeToString(oversizedPNGHeader(t))
		_, reasons := run(t, customer)
		if got := codes(reasons); got["selfie-photo"] != CodeInvalid {
			t.Errorf("reasons = %+v", reasons)
		}
	})
}

// oversizedPNGHeader returns a PNG that claims 8000x8000 pixels, patched in the
// IHDR chunk with its CRC recomputed.  DecodeConfig trusts the header, so the
// pixel limit must trip before any full decode is attempted.
func oversizedPNGHeader(t *testing.T) []byte {
	t.Helper()
	raw := pngBytes(t, 1, 1)
	const chunkType, widthOffset, heightOffset, crcOffset = 12, 16, 20, 29
	binary.BigEndian.PutUint32(raw[widthOffset:], 8000)
	binary.BigEndian.PutUint32(raw[heightOffset:], 8000)
	binary.BigEndian.PutUint32(raw[crcOffset:], crc32.ChecksumIEEE(raw[chunkType:crcOffset]))
	cfg, err := png.DecodeConfig(bytes.NewReader(raw))
	if err != nil || cfg.Width != 8000 {
		t.Fatalf("patched header not accepted by DecodeConfig: %v", err)
	}
	return raw
}

func TestEmailNeverRejects(t *testing.T) {
	values := map[string]any{"absent": nil, "empty": "", "number": json.Number("5"), "object": map[string]any{"a": 1}, "nul": "a\x00b"}
	for name, value := range values {
		t.Run(name, func(t *testing.T) {
			customer := validCustomer(t)
			customer["email"] = value
			if name == "absent" {
				delete(customer, "email")
			}
			c, reasons := run(t, customer)
			if len(reasons) != 0 {
				t.Fatalf("reasons = %+v", reasons)
			}
			if c.Email != nil {
				t.Errorf("email = %q, want treated as absent", *c.Email)
			}
		})
	}
}

func TestEligibility(t *testing.T) {
	tests := []struct {
		name   string
		date   time.Time
		field  string
		value  string
		reason string
	}{
		{"under 18 by one day", submissionDate, "date-of-birth", "2008-11-01", CodeCustomerUnderAge},
		{"18 today", submissionDate, "date-of-birth", "2008-10-31", ""},
		{"born today is under age, not invalid", submissionDate, "date-of-birth", "2026-10-31", CodeCustomerUnderAge},
		{"29 February birth on 28 February", day(2026, 2, 28), "date-of-birth", "2008-02-29", CodeCustomerUnderAge},
		{"29 February birth on 1 March", day(2026, 3, 1), "date-of-birth", "2008-02-29", ""},
		{"ID expired yesterday", submissionDate, "government-id-expiry-date", "2026-10-30", CodeGovernmentIDExpired},
		{"ID expires today", submissionDate, "government-id-expiry-date", "2026-10-31", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer := validCustomer(t)
			customer[tt.field] = tt.value
			body, err := json.Marshal(map[string]any{"customer": customer})
			if err != nil {
				t.Fatal(err)
			}
			_, reasons, err := Validate(body, tt.date, genders)
			if err != nil {
				t.Fatal(err)
			}
			got := codes(reasons)
			if tt.reason == "" {
				if len(got) != 0 {
					t.Errorf("reasons = %+v, want none", reasons)
				}
				return
			}
			if len(got) != 1 || got[tt.field] != tt.reason {
				t.Errorf("reasons = %+v, want only %s %s", reasons, tt.field, tt.reason)
			}
		})
	}
}

func TestSubmissionDateUsesItsOwnCalendarDay(t *testing.T) {
	// 23:30 on 31 October in UTC-5 is already 1 November in UTC.
	zone := time.FixedZone("UTC-5", -5*60*60)
	customer := validCustomer(t)
	customer["government-id-expiry-date"] = "2026-10-31"
	body, err := json.Marshal(map[string]any{"customer": customer})
	if err != nil {
		t.Fatal(err)
	}
	_, reasons, err := Validate(body, time.Date(2026, 10, 31, 23, 30, 0, 0, zone), genders)
	if err != nil {
		t.Fatal(err)
	}
	if len(reasons) != 0 {
		t.Errorf("reasons = %+v", reasons)
	}
}

func TestSkipsEligibilityWhenDateFailed(t *testing.T) {
	customer := validCustomer(t)
	customer["date-of-birth"] = "not-a-date"
	customer["government-id-expiry-date"] = ""
	_, reasons := run(t, customer)
	want := map[string]string{"date-of-birth": CodeInvalid, "government-id-expiry-date": CodeMissing}
	if got := codes(reasons); !maps.Equal(got, want) || len(reasons) != 2 {
		t.Errorf("reasons = %+v", reasons)
	}
}

func TestRetainedSections(t *testing.T) {
	got, err := RetainedSections([]byte(`{"customer":{"full-legal-name":"Alex"},"collateral":{"a":1},"loan":{"b":2}}`))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(got), "Alex") || !strings.Contains(string(got), `"collateral":{"a":1}`) ||
		!strings.Contains(string(got), `"loan":{"b":2}`) {
		t.Errorf("got %s", got)
	}
}

func day(year int, month time.Month, d int) time.Time {
	return time.Date(year, month, d, 0, 0, 0, 0, time.UTC)
}
