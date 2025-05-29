package bridge

import "time"

// HTMLPrintTime is used to print time into
type HTMLPrintTime func(t *time.Time) string

// PrintDate prints date in the way HTML expects
var PrintDate HTMLPrintTime = func(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

// PrintTime prints time in the way HTML expects
var PrintTime HTMLPrintTime = func(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("15:04")
}

// PrintDateTime prints date time in the way HTML expects
var PrintDateTime HTMLPrintTime = func(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02T15:04")
}
