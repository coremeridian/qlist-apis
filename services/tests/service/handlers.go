// Package service handles and routes access endpoints
package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"coremeridian.xyz/app/qlist/cmd/api/services/tests/models"
	"coremeridian.xyz/app/qlist/cmd/api/services/tests/utils"
)

func info(w http.ResponseWriter, r *http.Request) {
	infoMessage := struct {
		Message string
	}{
		Message: "This is info about the API",
	}
	json.NewEncoder(w).Encode(infoMessage)
}

func (app *HTTPApplication) showTests(w http.ResponseWriter, r *http.Request) {
	testOptions, err := models.NewTestOptionsFromQuery(r.URL.Query())
	if err != nil {
		app.HTTP.ServerError(w, err)
		return
	}
	testsArray, err := app.Tests.Latest(10, testOptions)
	if err != nil {
		app.HTTP.ServerError(w, err)
		return
	}

	tests, err := json.MarshalIndent(testsArray, "", " ")
	if err != nil {
		app.HTTP.ServerError(w, err)
		return
	}
	fmt.Println(string(tests))
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(tests))
}

func (app *HTTPApplication) createTest(w http.ResponseWriter, r *http.Request) {
	var test models.Test

	mr := utils.DecodeJSONBody(w, r, &test)
	if mr != nil {
		var err *utils.JsonMarshalError
		if errors.As(mr, &err) {
			app.HTTP.ErrorLog().Println(err)
			app.HTTP.ClientError(w, err.StatusCode)
		} else {
			app.HTTP.ServerError(w, err)
		}
		return
	}

	accessToken := r.Context().Value(accessTokenContext).(string)
	_, err := app.HTTP.Authorize(r.Context(), "resource:create:Test", test.Title, accessToken)
	if err != nil {
		app.HTTP.ErrorLog().Println(err)
		app.HTTP.ClientError(w, http.StatusForbidden)
		return
	}

	_, err = app.Tests.Insert(&test)
	if err != nil {
		app.HTTP.ServerError(w, err)
		return
	}

	testJSON, err := json.MarshalIndent(test, "", " ")
	if err != nil {
		app.HTTP.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(testJSON))
}
