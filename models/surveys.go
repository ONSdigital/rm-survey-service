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
	AllSurveys           *sql.Stmt
	GetSurvey            *sql.Stmt
	GetSurveyByShortName *sql.Stmt
	GetSurveyByReference *sql.Stmt
	GetSurveyID          *sql.Stmt
}

func NewAPI(db *sql.DB, id string) (*API, error) {
	allSurveysSql := "SELECT id, shortname FROM survey.survey ORDER BY shortname ASC"
	allSurveyStmt, err := createStmt(allSurveysSql, db)
	if err != nil {
		return nil, err
	}

	getSurveySql := "SELECT id, shortname, longname, surveyref from survey.survey WHERE id = %s"
	getSurveyStmt, err := createStmt(getSurveySql, db)
	if err != nil {
		return nil, err
	}

	getSurveyByShortNameSql := "SELECT id, shortname, longname, surveyref from survey.survey WHERE LOWER(shortName) = LOWER(%s)"
	getSurveyByShortNameStmt, err := createStmt(getSurveySql, db)
	if err != nil {
		return nil, err
	}

	getSurveyByReferenceSql := "SELECT id, shortname, longname, surveyref from survey.survey WHERE LOWER(surveyref) = LOWER(%s)"
	getSurveyByReferenceStmt, err := createStmt(getSurveyByReferenceSql, db)
	if err != nil {
		return nil, err
	}

	getSurveyIdSql := "SELECT id FROM survey.survey WHERE id = %s"
	getSurveyIdStmt, err := createStmt(getSurveyByReferenceSql, db)
	if err != nil {
		return nil, err
	}

	return &API{
		AllSurveys:           allSurveyStmt,
		GetSurvey:            getSurveyStmt,
		GetSurveyByShortName: getSurveyByShortNameStmt,
		GetSurveyByReference: getSurveyByReferenceStmt,
		GetSurveyID:          getSurveyIdStmt,
	}, nil
}

func Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(NewVersion()); err != nil {
		panic(err)
	}
}

// AllSurveys returns summaries of all known surveys. The surveys are returned in ascending short name order.
func (api *API) AllSurveys(w http.ResponseWriter, r *http.Request) error {
	rows, err := db.Query(api.AllSurveys)

	if err != nil {
		LogError("Error getting list of surveys", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	defer rows.Close()
	surveySummaries := make([]*SurveySummary, 0)

	for rows.Next() {
		surveySummary := new(SurveySummary)
		err := rows.Scan(&surveySummary.ID, &surveySummary.ShortName)

		if err != nil {
			return context.JSON(http.StatusInternalServerError, err.Error())
		}

		surveySummaries = append(surveySummaries, surveySummary)
	}

	if err = rows.Err(); err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, surveys)
}

// GetSurvey returns the details of the survey identified by the string surveyID.
func (api *API) GetSurvey(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["surveyId"]

	survey := new(Survey)
	err := db.QueryRow(fmt.Sprintf(api.GetSurveyStmnt, id)).Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference)
	if err != nil {
		if err == sql.ErrNoRows {
			re := NewRESTError("404", "Survey not found")
			return context.JSON(http.StatusNotFound, re)
		}

		logError("Error getting survey '"+surveyID+"'", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, survey)
}

// GetSurveyByShortName returns the details of the survey identified by the string shortName.
func (api *API) GetSurveyByShortName(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["shortName"]

	survey := new(Survey)
	err := db.QueryRow(fmt.Sprintf(api.GetSurveyByShortNameStmt, id))
	Scan(&survey.ID,
		&survey.ShortName,
		&survey.LongName,
		&survey.Reference,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			re := NewRESTError("404", "Survey not found")
			return context.JSON(http.StatusNotFound, re)
		}

		LogError("Error getting survey '"+shortName+"'", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, survey)
}

// GetSurveyByReference returns the details of the survey identified by the string reference.
func (api *API) GetSurveyByReference(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["ref"]

	survey := new(Survey)
	err := db.QueryRow(fmt.Sprintf(api.GetSurveyByReferenceStmt, id)).Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference)
	if err != nil {
		if err == sql.ErrNoRows {
			re := NewRESTError("404", "Survey not found")
			return context.JSON(http.StatusNotFound, re)
		}

		LogError("Error getting survey '"+reference+"'", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, survey)
}

// AllClassifierTypeSelectors returns all the classifier type selectors for the survey identified by the string surveyID. The classifier type selectors are returned in ascending order.
func AllClassifierTypeSelectors(w http.ResponseWriter, r *http.Request) error {
	// We need to run a query first to check if the survey exists so an HTTP 404 can be correctly
	// returned if it doesn't exist. Without this check an HTTP 204 is incorrectly returned for an
	// invalid survey ID.
	vars := mux.Vars(r)
	colllectionId := vars["colllectionId"]
	surveyId := vars["surveyId"]

	err := getSurveyID(surveyID)
	if err != nil {
		LogError("Error getting list of classifier type selectors for survey '"+surveyID+"'", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	// Now we can get the classifier type selector records.
	rows, err := db.Query("SELECT classifiertypeselector.id, classifiertypeselector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.surveyfk = survey.surveypk WHERE survey.id = $1 ORDER BY classifiertypeselector ASC", surveyID)
	if err != nil {
		re := NewRESTError("404", "Survey not found")
		return context.JSON(http.StatusNotFound, re)
	}

	defer rows.Close()
	classifierTypeSelectorSummaries := make([]*ClassifierTypeSelectorSummary, 0)

	for rows.Next() {
		classifierTypeSelectorSummary := new(ClassifierTypeSelectorSummary)
		err := rows.Scan(&classifierTypeSelectorSummary.ID, &classifierTypeSelectorSummary.Name)

		if err != nil {
			LogError("Error getting list of classifier type selectors for survey '"+surveyID+"'", err)
			return context.JSON(http.StatusInternalServerError, err.Error())
		}

		classifierTypeSelectorSummaries = append(classifierTypeSelectorSummaries, classifierTypeSelectorSummary)
	}

	if err = rows.Err(); err != nil {
		LogError("Error getting list of classifier type selectors for survey '"+surveyID+"'", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	if err != nil {
		if err == sql.ErrNoRows {
			re := NewRESTError("404", "Survey not found")
			return context.JSON(http.StatusNotFound, re)
		}

		LogError("Error getting list of classifier type selectors for survey '"+surveyID+"'", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(classifierTypeSelectors) == 0 {
		return context.String(http.StatusNoContent, "")
	}

	return context.JSON(http.StatusOK, classifierTypeSelectors)
}

// GetClassifierTypeSelector returns the details of the classifier type selector for the survey identified by the string surveyID and
// the classifier type selector identified by the string classifierTypeSelectorID.
func GetClassifierTypeSelector(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["surveyId"]
	classifierType := vars["classifierTypeSelectorId"]

	// We need to run two queries first to check if the survey and classifier type selector both exist
	// so an HTTP 404 can be correctly returned if they don't exist. Without this check and HTTP 204 is
	// incorrectly return for an invalid survey ID or classifier type selector ID.
	err := getSurveyID(surveyID)
	if err != nil {
		LogError("Error getting classifier type selector '"+classifierTypeSelectorID+"' for survey '"+surveyID+"'", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	err := db.QueryRow("SELECT id FROM survey.classifiertypeselector WHERE id = $1", classifierTypeSelectorID).Scan(&id)
	if err != nil {
		LogError("Error getting classifier type selector '"+classifierTypeSelectorID+"' for survey '"+surveyID+"'", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	// Now we can get the classifier type selector and classifier type records.
	rows, err := db.Query("SELECT id, classifiertypeselector, classifiertype FROM survey.classifiertype INNER JOIN survey.classifiertypeselector ON classifiertype.classifiertypeselectorfk = classifiertypeselector.classifiertypeselectorpk WHERE classifiertypeselector.id = $1 ORDER BY classifiertype ASC", classifierTypeSelectorID)
	if err != nil {
		LogError("Error getting classifier type selector '"+classifierTypeSelectorID+"' for survey '"+surveyID+"'", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	defer rows.Close()
	classifierTypeSelector := new(ClassifierTypeSelector)

	// Using make here ensures the JSON contains an empty array if there are no classifier
	// types, rather than null.
	classifierTypes := make([]string, 0)
	var classifierType string

	for rows.Next() {
		err := rows.Scan(&classifierTypeSelector.ID, &classifierTypeSelector.Name, &classifierType)

		if err != nil {
			LogError("Error getting classifier type selector '"+classifierTypeSelectorID+"' for survey '"+surveyID+"'", err)
			return context.JSON(http.StatusInternalServerError, err.Error())
		}

		classifierTypes = append(classifierTypes, classifierType)
	}

	classifierTypeSelector.ClassifierTypes = classifierTypes

	if err = rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			re := NewRESTError("404", "Survey or classifier type selector not found")
			return context.JSON(http.StatusNotFound, re)
		}

		LogError("Error getting classifier type selector '"+classifierTypeSelectorID+"' for survey '"+surveyID+"'", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, classifierTypeSelector)
}

func (api *API) getSurveyID(surveyId string) error {
	var id string
	return db.QueryRow(api.GetSurveyIDStmt).Scan(&id)
}

func createStmt(sqlStatement string, db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare(sqlStatement)
}
