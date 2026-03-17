package personas

import (
	"encoding/csv"
	"io"
	"os"
	"strings"

	"reqx/internal/errs"
)

// Persona is one CSV row keyed by column name.
// Values are always treated as strings (no type coercion).
type Persona map[string]string

// LoadCSV reads a header-based CSV file into personas.
// Empty lines are ignored; missing cells become empty strings.
func LoadCSV(path string) ([]Persona, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errs.Wrap(err, errs.KindInvalidInput, "could not open personas CSV")
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = true

	headers, err := r.Read()
	if err != nil {
		if err == io.EOF {
			return nil, errs.InvalidInput("personas CSV is empty")
		}
		return nil, errs.Wrap(err, errs.KindInvalidInput, "could not read personas CSV header")
	}

	for i := range headers {
		headers[i] = strings.TrimSpace(strings.TrimPrefix(headers[i], "\uFEFF"))
	}

	seen := map[string]bool{}
	cols := make([]string, 0, len(headers))
	for _, h := range headers {
		h = strings.TrimSpace(h)
		if h == "" {
			continue
		}
		if seen[h] {
			continue
		}
		seen[h] = true
		cols = append(cols, h)
	}
	if len(cols) == 0 {
		return nil, errs.InvalidInput("personas CSV header has no valid columns")
	}

	out := make([]Persona, 0, 16)
	for {
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errs.Wrap(err, errs.KindInvalidInput, "could not read personas CSV row")
		}
		if len(rec) == 0 {
			continue
		}

		p := make(Persona, len(cols))
		for colIdx, col := range cols {
			if colIdx < len(rec) {
				p[col] = strings.TrimSpace(rec[colIdx])
			} else {
				p[col] = ""
			}
		}
		out = append(out, p)
	}
	return out, nil
}

