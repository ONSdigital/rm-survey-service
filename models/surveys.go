package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// ClassifierTypeSelectorSummary represents a summary of a classifier type selector.
type ClassifierTypeSelectorSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ClassifierTypeSelector represents the detail of a classifier type selector.
type ClassifierTypeSelector struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	ClassifierTypes []string `json:"classifierTypes"`
}

// SurveySummary represents a summary of a survey.
type SurveySummary struct {
	ID        string `json:"id"`
	ShortName string `json:"shortName"`
}

// Survey represents the details of a survey.
type Survey struct {
	ID        string `json:"id"`
	ShortName string `json:"shortName"`
	LongName  string `json:"longName"`
	Reference string `json:"surveyRef"`
}

type API struct {
	AllSurveysStmt           *sql.Stmt
	GetSurveyStmt            *sql.Stmt
	GetSurveyByShortNameStmt *sql.Stmt
	GetSurveyByReferenceStmt *sql.Stmt
	GetSurveyIDStmt          *sql.Stmt
}

func NewAPI(db *sql.DB) (*API, error) {
	allSurveysSQL := "SELECT id, shortname FROM survey.survey ORDER BY shortname ASC"
	allSurveyStmt, err := createStmt(allSurveysSQL, db)
	if err != nil {
		return nil, err
	}

	getSurveySQL := "SELECT id, shortname, longname, surveyref from survey.survey WHERE id = $1"
	getSurveyStmt, err := createStmt(getSurveySQL, db)
	if err != nil {
		return nil, err
	}

	getSurveyByShortNameSQL := "SELECT id, shortname, longname, surveyref from survey.survey WHERE LOWER(shortName) = LOWER($1)"
	getSurveyByShortNameStmt, err := createStmt(getSurveyByShortNameSQL, db)
	if err != nil {
		return nil, err
	}

	getSurveyByReferenceSQL := "SELECT id, shortname, longname, surveyref from survey.survey WHERE LOWER(surveyref) = LOWER($1)"
	getSurveyByReferenceStmt, err := createStmt(getSurveyByReferenceSQL, db)
	if err != nil {
		return nil, err
	}

	getSurveyIdSQL := "SELECT id FROM survey.survey WHERE id = $1"
	getSurveyIdStmt, err := createStmt(getSurveyIdSQL, db)
	if err != nil {
		return nil, err
	}

	return &API{
		AllSurveysStmt:           allSurveyStmt,
		GetSurveyStmt:            getSurveyStmt,
		GetSurveyByShortNameStmt: getSurveyByShortNameStmt,
		GetSurveyByReferenceStmt: getSurveyByReferenceStmt,
		GetSurveyIDStmt:          getSurveyIdStmt,
	}, nil
}

func (api *API) Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(NewVersion()); err != nil {
		panic(err)
	}
}

// AllSurveys returns summaries of all known surveys. The surveys are returned in ascending short name order.
func (api *API) AllSurveys(w http.ResponseWriter, r *http.Request) {
	rows, err := api.AllSurveysStmt.Query()

	if err != nil {
		//LogError("Error getting list of surveys", err)
		http.Error(w, "AllSurveys query failed", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	surveySummaries := make([]*SurveySummary, 0)

	for rows.Next() {
		surveySummary := new(SurveySummary)
		err := rows.Scan(&surveySummary.ID, &surveySummary.ShortName)

		if err != nil {
			http.Error(w, "No surveys found", http.StatusInternalServerError)
		}

		surveySummaries = append(surveySummaries, surveySummary)
	}

	if len(surveySummaries) == 0 {
		http.Error(w, "No surveys found", http.StatusNoContent)
		return
	}
	data, _ := json.Marshal(surveySummaries)
	w.Write(data)
}

// GetSurvey returns the details of the survey identified by the string surveyID.
func (api *API) GetSurvey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["surveyId"]
	survey := new(Survey)
	surveyRow := api.GetSurveyStmt.QueryRow(id)
	err := surveyRow.Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Survey not found", http.StatusNotFound)
		} else {
			http.Error(w, "get survey query failed", http.StatusInternalServerError)
		}
	}

	data, _ := json.Marshal(survey)
	w.Write(data)
}

// GetSurveyByShortName returns the details of the survey identified by the string shortName.
func (api *API) GetSurveyByShortName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["shortName"]
	surveyRow := api.GetSurveyByShortNameStmt.QueryRow(id)
	survey := new(Survey)
	err := surveyRow.Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Survey not found", http.StatusNotFound)
		} else {
			fmt.Println("ERROR: " + err.Error())
			http.Error(w, "get survey by shortname query failed", http.StatusInternalServerError)
		}
	}

	data, _ := json.Marshal(survey)
	w.Write(data)

}

// GetSurveyByReference returns the details of the survey identified by the string ref.
func (api *API) GetSurveyByReference(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ref"]
	surveyRow := api.GetSurveyByReferenceStmt.QueryRow(id)
	survey := new(Survey)
	err := surveyRow.Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Survey not found", http.StatusNotFound)
		} else {
			http.Error(w, "get survey by reference query failed", http.StatusInternalServerError)
		}
	}

	data, _ := json.Marshal(survey)
	w.Write(data)

}

// AllClassifierTypeSelectors returns all the classifier type selectors for the survey identified by the string surveyID. The classifier type selectors are returned in ascending order.
func (api *API) AllClassifierTypeSelectors(w http.ResponseWriter, r *http.Request) {
	// We need to run a query first to check if the survey exists so an HTTP 404 can be correctly
	// returned if it doesn't exist. Without this check an HTTP 204 is incorrectly returned for an
	// invalid survey ID.
	vars := mux.Vars(r)
	surveyId := vars["surveyId"]

	err := api.getSurveyID(surveyId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Survey not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error getting list of classifier type selectors for survey '"+surveyId+"' - "+err.Error(), http.StatusInternalServerError)
		}
	}

	// Now we can get the classifier type selector records.
	rows, err := db.Query("SELECT classifiertypeselector.id, classifiertypeselector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.surveyfk = survey.surveypk WHERE survey.id = $1 ORDER BY classifiertypeselector ASC", surveyId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Classifier type selector not found", http.StatusNotFound)
		} else {
			fmt.Println("HELLO")
			http.Error(w, "Error getting list of classifier type selectors for survey '"+surveyId+"' - "+err.Error(), http.StatusInternalServerError)
		}
	}

	defer rows.Close()
	classifierTypeSelectorSummaries := make([]*ClassifierTypeSelectorSummary, 0)

	for rows.Next() {
		classifierTypeSelectorSummary := new(ClassifierTypeSelectorSummary)
		err := rows.Scan(&classifierTypeSelectorSummary.ID, &classifierTypeSelectorSummary.Name)

		if err != nil {
			//LogError("Error getting list of classifier type selectors for survey '"+surveyID+"'", err)
			http.Error(w, "Error getting list of classifier type selectors for survey '"+surveyId+"' - "+err.Error(), http.StatusInternalServerError)
			return
		}

		classifierTypeSelectorSummaries = append(classifierTypeSelectorSummaries, classifierTypeSelectorSummary)
	}

	data, _ := json.Marshal(classifierTypeSelectorSummaries)
	w.Write(data)
}

func (api *API) getSurveyID(surveyId string) error {
	var id string
	return api.GetSurveyIDStmt.QueryRow(surveyId).Scan(&id)
}

func createStmt(sqlStatement string, db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare(sqlStatement)
}

func (api *API) Close() {
	api.AllSurveysStmt.Close()
}
