package models

import (
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func Router(api *API) *mux.Router {
	// Webserver - strictslash set to true to match trailing slashes to routes
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/info", api.Info).Methods("GET")
	r.HandleFunc("/surveys", use(api.AllSurveys, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/{surveyId}", use(api.GetSurvey, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/shortname/{shortName}", use(api.GetSurveyByShortName, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/ref/{ref}", use(api.GetSurveyByReference, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/{surveyId}/classifiertypeselectors", use(api.AllClassifierTypeSelectors, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/{surveyId}/classifiertypeselectors/{classifierTypeSelectorId}", use(api.GetClassifierTypeSelectorByID, basicAuth)).Methods("GET")
	return r
}

func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	// Taken from https://gist.github.com/elithrar/9146306
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		if pair[0] != os.Getenv("security_user_name") || pair[1] != os.Getenv("security_user_password") {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	}
}
