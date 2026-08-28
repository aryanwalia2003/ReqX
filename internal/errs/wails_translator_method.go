package errs

import (
	"encoding/json"
	"errors"
	"log"
)

// WailsErrorResponse ek bound method ke error ka JSON shape hai jo frontend
// tak pahunchta — dekho app/frontend/src/lib/errors.ts, jo isi shape ko parse
// karta hai.
type WailsErrorResponse struct {
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

// FormatForWails app/main.go me options.App.ErrorFormatter ke taur pe wire
// hota — har bound method ka returned error frontend tak pahunchne se pehle
// yahan se guzarta, taaki hamesha {kind, message} JSON ki tarah mile, bare Go
// error string ki tarah nahi. Same convention jo WriteHTTPError (CLI/UI side)
// already follow karta hai.
func FormatForWails(err error) any {
	var appErr AppError
	if errors.As(err, &appErr) {
		log.Printf("[ERR] %s | Cause: %v | Meta: %v\nStack: %s\n",
			appErr.Kind(), appErr.Unwrap(), appErr.Metadata(), appErr.StackTrace())

		payload, _ := json.Marshal(WailsErrorResponse{
			Kind:    string(appErr.Kind()),
			Message: appErr.Message(),
		})
		return string(payload)
	}

	// errs.Wrap/New se na guzra hua error — client ko internal detail kabhi
	// mat dikhao, bas server-side log karo taaki missing wrap pakda jaye.
	log.Printf("[ERR] Unhandled Standard Error: %v\n", err)
	payload, _ := json.Marshal(WailsErrorResponse{
		Kind:    string(KindInternal),
		Message: "An unexpected error occurred.",
	})
	return string(payload)
}
