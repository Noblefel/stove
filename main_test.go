package main

import (
	"strings"
	"testing"
)

func TestReadCSV(t *testing.T) {
	t.Run("get rows", func(t *testing.T) {
		s := strings.NewReader("Name,Address")
		want := "<tr><td>Name</td><td>Address</td></tr>"

		rows, err := readCSV(s, false)
		if err != nil {
			t.Errorf("returned error: %v", err)
		}

		if !strings.EqualFold(want, rows) {
			t.Errorf("want %q, got %q", want, rows)
		}
	})

	t.Run("get rows with num", func(t *testing.T) {
		s := strings.NewReader("Name\nJohn Doe")
		want := "<tr><td>No</td><td>Name</td></tr><tr><td>1</td><td>John Doe</td></tr>"

		rows, err := readCSV(s, true)
		if err != nil {
			t.Errorf("returned error: %v", err)
		}

		if !strings.EqualFold(want, rows) {
			t.Errorf("want %q, got %q", want, rows)
		}
	})

	t.Run("fail when reading line", func(t *testing.T) {
		s := strings.NewReader(`"`)

		rows, err := readCSV(s, false)
		if err == nil {
			t.Errorf("should return error, but got %q", rows)
		}
	})
}

func TestSetupHTML(t *testing.T) {
	t.Run("get html string", func(t *testing.T) {
		htmlString := setupHTML([]byte("[%title%]"), "", "sample title")
		want := "sample title"

		if !strings.EqualFold(want, htmlString) {
			t.Errorf("want %q, got %q", want, htmlString)
		}
	})
}
