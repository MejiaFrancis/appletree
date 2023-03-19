// Filename: cmd/api/schools.go

package main

import (
	"fmt"
	"net/http"
	"time"

	"appletree.francismejia.net/internal/data"
)

// createSchoolHandler for the "POST /v1/schools" endpoint
func (app *application) createSchoolHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new school..")
}

// showSchoolHandler for the "GET /v1/schools/:id" endpoint
func (app *application) showSchoolHandler(w http.ResponseWriter, r *http.Request) {
	//params := httprouter.ParamsFromContext(r.Context())
	//id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w,r)
		//http.NotFound(w, r)
		return
	}
	// Display the school id
	//fmt.Fprintf(w, "show the details for school %d\n", id)
	school := data.School{
		ID:       id,
		CreateAt: time.Now(),
		Name:     "University of Belmopan",
		Level:    "University",
		Contact:  "Abel Blanco",
		Phone:    "323-4545",
		Email:    "university@ub.bz",
		Website:  "https://uob.edu.bz",
		Address:  "17 Apple Avenue",
		Mode:     []string{"blended", "online", "face-to-face"},
		Version:  1,
	}
	err = app.WriteJSON(w, http.StatusOK, envelope{"school" : school}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		
		return
	}
}
