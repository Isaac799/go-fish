package bridge

import (
	"testing"
	"time"
)

type mockPrimitivesOne struct {
	Name           string    `gofish:"name"`
	Alive          bool      `gofish:"alive"`
	Age            uint      `gofish:"age"`
	FavoriteNumber float32   `gofish:"favorite number"`
	Birthday       time.Time `gofish:"birthday"`
	ClockIn        time.Time `gofish:"clock in"`
	VacationStart  time.Time `gofish:"vacation start"`
}

type mockPrimitivesMany struct {
	Name           []string    `gofish:"name"`
	Alive          []bool      `gofish:"alive"`
	Age            []uint      `gofish:"age"`
	FavoriteNumber []float32   `gofish:"favorite number"`
	Birthday       []time.Time `gofish:"birthday"`
	ClockIn        []time.Time `gofish:"clock in"`
	VacationStart  []time.Time `gofish:"vacation start"`
}

func TestReflection(t *testing.T) {
	m := map[string]string{
		"name":            "Jane Doe",
		"alive":           "true",
		"age":             "27",
		"favorite number": "87.3",
		"birthday":        "1970-01-01",
		"clock in":        "10:00",
		"vacation start":  "1998-01-01T10:15",
		"favoriteColors":  "0,2",
	}
	person := mockPrimitivesOne{}

	birthday, _ := time.Parse(TimeFormatHTMLDate, "1970-01-01")
	clockIn, _ := time.Parse(TimeFormatHTMLTime, "10:00")
	vacationStart, _ := time.Parse(TimeFormatHTMLDateTime, "1998-01-01T10:15")

	AttributesToStruct(m, &person)

	assert(t, person.Name, "Jane Doe")
	assert(t, person.Alive, true)
	assert(t, person.Age, 27)
	assert(t, person.FavoriteNumber, 87.3)
	assert(t, person.Birthday, birthday)
	assert(t, person.ClockIn, clockIn)
	assert(t, person.VacationStart, vacationStart)
}

func TestReflectionSlice(t *testing.T) {
	m := map[string]string{
		"name":            "Jane Doe, John Doe",
		"alive":           "true, false",
		"age":             "27, 34",
		"favorite number": "87.3, 19.4",
		"birthday":        "1970-01-01, 1988-01-01",
		"clock in":        "10:00, 3:00",
		"vacation start":  "1998-01-01T10:15, 1988-01-01T3:15",
		"favoriteColors":  "0,2",
	}
	person := mockPrimitivesMany{}

	birthday, _ := time.Parse(TimeFormatHTMLDate, "1970-01-01")
	clockIn, _ := time.Parse(TimeFormatHTMLTime, "10:00")
	vacationStart, _ := time.Parse(TimeFormatHTMLDateTime, "1998-01-01T10:15")

	birthday2, _ := time.Parse(TimeFormatHTMLDate, "1988-01-01")
	clockIn2, _ := time.Parse(TimeFormatHTMLTime, "3:00")
	vacationStart2, _ := time.Parse(TimeFormatHTMLDateTime, "1988-01-01T3:15")

	AttributesToStruct(m, &person)

	assertIndexes(t, person.Name, []string{"Jane Doe", "John Doe"})
	assertIndexes(t, person.Alive, []bool{true, false})
	assertIndexes(t, person.Age, []uint{27, 34})
	assertIndexes(t, person.FavoriteNumber, []float32{87.3, 19.4})
	assertIndexes(t, person.Birthday, []time.Time{birthday, birthday2})
	assertIndexes(t, person.ClockIn, []time.Time{clockIn, clockIn2})
	assertIndexes(t, person.VacationStart, []time.Time{vacationStart, vacationStart2})
}

func TestFormValue(t *testing.T) {
	el := mockDateEl()
	expectedBirthday, _ := time.Parse(TimeFormatHTMLDate, "1970-01-01")
	el.Attributes["value"] = "1970-01-01"
	form := el.Form()
	birthday, err := ValueOf[time.Time](form, "birthday")
	assertNoError(t, err)
	assert(t, birthday, expectedBirthday)
}

func TestAttributesToStruct(t *testing.T) {
	el := mockDateEl()
	expectedBirthday, _ := time.Parse(TimeFormatHTMLDate, "1970-01-01")
	el.Attributes["value"] = "1970-01-01"
	validation, err := ValidationForElement[time.Time](el)
	birthday := validation.Value
	assertNoError(t, err)
	assert(t, birthday, expectedBirthday)
}
