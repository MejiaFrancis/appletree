// Filename: cmd/api/healthcheck.go

package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	js := `{"status" : "avialable", "environment":%q, "version": %q}`
	js = fmt.Sprintf(js, app.config.env, version)
	// add header information
	w.Header().Set("Content-Type", "application/json") //this is how you set a header
	w.Write([]byte(js))
	//fmt.Fprintln(w, "status: available")  // places in the js on previous line to create a JSON file
	//fmt.Fprintf(w, "environment: %s\n", app.config.env) // places in the js on previous line
	//fmt.Fprintf(w, "version: %s\n", version)  // places in the js on previous line
}
