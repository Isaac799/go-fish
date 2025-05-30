package bridge

import "time"

const (
	// TimeFormatHTMLTime is the go format string for HTML time
	TimeFormatHTMLTime = "15:04"
	// TimeFormatHTMLDate is the go format string for HTML date
	TimeFormatHTMLDate = "2006-01-02"
	// TimeFormatHTMLDateTime is the go format string for HTML date time
	TimeFormatHTMLDateTime = "2006-01-02T15:04"
)

// HTMLPrintTime is used to print time into
type HTMLPrintTime func(t *time.Time) string

// PrintDate prints date in the way HTML expects
var PrintDate HTMLPrintTime = func(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(TimeFormatHTMLDate)
}

// PrintTime prints time in the way HTML expects
var PrintTime HTMLPrintTime = func(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(TimeFormatHTMLTime)
}

// PrintDateTime prints date time in the way HTML expects
var PrintDateTime HTMLPrintTime = func(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(TimeFormatHTMLDateTime)
}
