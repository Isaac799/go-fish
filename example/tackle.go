package main

import "time"

func printTime(t time.Time) string {
	return t.Format(time.RFC1123)
}
