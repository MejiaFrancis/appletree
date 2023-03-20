// Filename: cmd/api/healthcheck.go

package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "available",
		"sysem_info": map[string]string{
			"envvironment": app.config.env,
			"version":      version,
		},
	}
	//js := `{"status" : "avialable", "environment":%q, "version": %q}`
	//js = fmt.Sprintf(js, app.config.env, version)
	// add header information
	//w.Header().Set("Content-Type", "application/json") //this is how you set a header
	//w.Write([]byte(js))
	//fmt.Fprintln(w, "status: available")  // places in the js on previous line to create a JSON file
	//fmt.Fprintf(w, "environment: %s\n", app.config.env) // places in the js on previous line
	//fmt.Fprintf(w, "version: %s\n", version)  // places in the js on previous line

	err := app.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
