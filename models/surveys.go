package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
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

// Survey represents the details of a survey.
type Survey struct {
	ID         string `json:"id"`
	ShortName  string `json:"shortName"`
	LongName   string `json:"longName"`
	Reference  string `json:"surveyRef"`
	LegalBasis string `json:"legalBasis"`
	LegalBasisRef string `json:"legalBasisRef"`
}

// LegalBasis - the legal basis for a survey consisting of a short reference and a long name
type LegalBasis struct {
	Reference  string `json:"ref"`
	LongName   string `json:"longName"`
}

//API contains all the pre-prepared sql statments
type API struct {
	AllSurveysStmt                    *sql.Stmt
	GetSurveyStmt                     *sql.Stmt
	GetSurveyByShortNameStmt          *sql.Stmt
	GetSurveyByReferenceStmt          *sql.Stmt
	GetSurveyIDStmt                   *sql.Stmt
	GetClassifierTypeSelectorStmt     *sql.Stmt
	GetClassifierTypeSelectorByIDStmt *sql.Stmt
	GetSurveyRefStmt                  *sql.Stmt
	PutSurveyDetailsBySurveyRefStmt   *sql.Stmt
	CreateSurveyStmt                  *sql.Stmt
	GetLegalBasesStmt                 *sql.Stmt
	GetLegalBasisFromLongNameStmt     *sql.Stmt
	GetLegalBasisFromRefStmt          *sql.Stmt
}

//NewAPI returns an API struct populated with all the created SQL statements
func NewAPI(db *sql.DB) (*API, error) {
	allSurveyStmt, err := createStmt("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref ORDER BY shortname ASC", db)
	if err != nil {
		return nil, err
	}

	getSurveyStmt, err := createStmt("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref WHERE id = $1", db)
	if err != nil {
		return nil, err
	}

 	getSurveyByShortNameStmt, err := createStmt("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref  WHERE LOWER(shortName) = LOWER($1)", db)
	if err != nil {
		return nil, err
	}

	getSurveyByReferenceStmt, err := createStmt("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref  WHERE LOWER(surveyref) = LOWER($1)", db)
	if err != nil {
		return nil, err
	}

	getSurveyIDStmt, err := createStmt("SELECT id FROM survey.survey WHERE id = $1", db)
	if err != nil {
		return nil, err
	}

	getClassifierTypeSelectorStmt, err := createStmt("SELECT classifiertypeselector.id, classifiertypeselector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.surveyfk = survey.surveypk WHERE survey.id = $1 ORDER BY classifiertypeselector ASC", db)
	if err != nil {
		return nil, err
	}

	getClassifierTypeSelectorByIDStmt, err := createStmt("SELECT id, classifiertypeselector, classifiertype FROM survey.classifiertype INNER JOIN survey.classifiertypeselector ON classifiertype.classifiertypeselectorfk = classifiertypeselector.classifiertypeselectorpk WHERE classifiertypeselector.id = $1 ORDER BY classifiertype ASC", db)
	if err != nil {
		return nil, err
	}

	getSurveyRefStmt, err := createStmt("SELECT surveyref FROM survey.survey WHERE LOWER(surveyref) = LOWER($1)", db)
	if err != nil {
		return nil, err
	}

	putSurveyDetailsBySurveyRefStmt, err := createStmt("UPDATE survey.survey SET shortname = $2, longname = $3 WHERE LOWER(surveyref) = LOWER($1)", db)
	if err != nil {
		return nil, err
	}

	createSurvey, err := createStmt("INSERT INTO survey.survey ( surveypk, id, surveyref, shortname, longname, legalbasis ) VALUES ( nextval('survey.survey_surveypk_seq'), $1, $2, $3, $4, $5)", db)
	if err != nil {
		return nil, err
	}

	getLegalBases, err := createStmt("SELECT ref, longname FROM survey.legalbasis", db)
	if err != nil {
		return nil, err
	}

	getLegalBasisFromLongName, err := createStmt("SELECT ref, longname FROM survey.legalbasis WHERE longname = $1", db)
	if err != nil {
		return nil, err
	}

	getLegalBasisFromRef, err := createStmt("SELECT ref, longname FROM survey.legalbasis WHERE ref = $1", db)
	if err != nil {
		return nil, err
	}

	return &API{
		AllSurveysStmt:                    allSurveyStmt,
		GetSurveyStmt:                     getSurveyStmt,
		GetSurveyByShortNameStmt:          getSurveyByShortNameStmt,
		GetSurveyByReferenceStmt:          getSurveyByReferenceStmt,
		GetSurveyIDStmt:                   getSurveyIDStmt,
		GetClassifierTypeSelectorStmt:     getClassifierTypeSelectorStmt,
		GetClassifierTypeSelectorByIDStmt: getClassifierTypeSelectorByIDStmt,
		GetSurveyRefStmt:                  getSurveyRefStmt,
		PutSurveyDetailsBySurveyRefStmt:   putSurveyDetailsBySurveyRefStmt,
		CreateSurveyStmt:                  createSurvey,
		GetLegalBasesStmt:                 getLegalBases,
		GetLegalBasisFromLongNameStmt:     getLegalBasisFromLongName,
		GetLegalBasisFromRefStmt:          getLegalBasisFromRef }, nil
}

// PostSurveyDetails endpoint handler - creates a new survey based on JSON in request
func (api *API) PostSurveyDetails(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	var postData Survey
	err = json.Unmarshal(body, &postData)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
	}

	shortName := postData.ShortName
	longName := postData.LongName
	surveyRef := postData.Reference
	legalBasisLongName := postData.LegalBasis
	legalBasisRef := postData.LegalBasisRef
	var legalBasis LegalBasis;
	var errorMessage string;

	if legalBasisRef != "" {
		legalBasis, err = api.getLegalBasisFromRef(legalBasisRef)
		errorMessage = fmt.Sprintf("Legal basis with reference %v does not exist", legalBasisRef)
	} else if legalBasisLongName != "" {
		legalBasis, err = api.getLegalBasisFromLongName(legalBasisLongName)
		errorMessage = fmt.Sprintf("Legal basis %v does not exist", legalBasisLongName)
	} else {
		http.Error(w, "No legal basis specified for survey", http.StatusBadRequest)
		return
	}

	if err == sql.ErrNoRows {
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Error getting legal basis - %v", err), http.StatusInternalServerError)
		return
	}

	surveyID := uuid.NewV4()

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate a UUID for new survey - %v", err), http.StatusInternalServerError)
		return
	}

	err = api.getSurveyRef(surveyRef)

	if err == sql.ErrNoRows {
		// Reference is unique - this is good
		_, err = api.CreateSurveyStmt.Exec(surveyID, surveyRef, shortName, longName, legalBasis.Reference)

		if err != nil {
			http.Error(w, fmt.Sprintf("Create survey details failed - %v", err), http.StatusInternalServerError)
			return
		}

		var result Survey
		var js []byte

		result.Reference = surveyRef
		result.LongName = longName
		result.ShortName = shortName
		result.ID = surveyID.String()
		result.LegalBasisRef = legalBasis.Reference
		result.LegalBasis = legalBasis.LongName

		js, err = json.Marshal(&result)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write(js)
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Failed to validate survey ref - %v", err), http.StatusInternalServerError)
		return
	} else {
		http.Error(w, fmt.Sprintf("Survey with reference %v already exists", surveyRef), http.StatusConflict)
		return
	}
}

// PutSurveyDetails endpoint handler changes a survey short name using the survey reference
func (api *API) PutSurveyDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	surveyRef := vars["ref"]

	body, err := ioutil.ReadAll(r.Body)

	var putData Survey
	err = json.Unmarshal(body, &putData)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
	}

	shortName := putData.ShortName
	longName := putData.LongName

	err = api.getSurveyRef(surveyRef)

	if err == sql.ErrNoRows {
		re := NewRESTError("404", "Survey not found")
		data, err := json.Marshal(re)
		if err != nil {
			http.Error(w, "Error marshaling NewRestError JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write(data)

		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get survey ref - %v", err), http.StatusInternalServerError)
		return
	}

	_, err = api.PutSurveyDetailsBySurveyRefStmt.Exec(surveyRef, shortName, longName)

	if err != nil {
		http.Error(w, fmt.Sprintf("Update survey details query failed - %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

//Info endpoint handler returns info like name, version, origin, commit, branch
//and built
func (api *API) Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(NewVersion()); err != nil {
		http.Error(w, "info encoding failed", http.StatusInternalServerError)
	}
}

// AllSurveys returns summaries of all known surveys. The surveys are returned in ascending short name order.
func (api *API) AllSurveys(w http.ResponseWriter, r *http.Request) {
	rows, err := api.AllSurveysStmt.Query()

	if err != nil {
		http.Error(w, "AllSurveys query failed", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	surveys := make([]*Survey, 0)

	for rows.Next() {
		survey := new(Survey)
		err = rows.Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference, &survey.LegalBasisRef, &survey.LegalBasis)

		if err != nil {
			http.Error(w, "Failed to get surveys from database", http.StatusInternalServerError)
			return
		}

		surveys = append(surveys, survey)
	}

	if len(surveys) == 0 {
		http.Error(w, "No surveys found", http.StatusNoContent)
		return
	}

	data, err := json.Marshal(surveys)
	if err != nil {
		http.Error(w, "Failed to marshal survey summary JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// AllLegalBases returns details of all legal bases
func (api *API) AllLegalBases(w http.ResponseWriter, r *http.Request) {
	rows, err := api.GetLegalBasesStmt.Query()

	if err != nil {
		http.Error(w, fmt.Sprintf("AllLegalBases query failed - %v", err), http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	legalBases := make([]*LegalBasis, 0)

	for rows.Next() {
		legalBasis := new(LegalBasis)
		err = rows.Scan(&legalBasis.Reference, &legalBasis.LongName)

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get legal bases from database - %v", err), http.StatusInternalServerError)
			return
		}

		legalBases = append(legalBases, legalBasis)
	}

	if len(legalBases) == 0 {
		http.Error(w, "No legal bases found", http.StatusNoContent)
		return
	}

	data, err := json.Marshal(legalBases)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal survey summary JSON - %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetSurvey returns the details of the survey identified by the string surveyID.
func (api *API) GetSurvey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["surveyId"]
	survey := new(Survey)
	surveyRow := api.GetSurveyStmt.QueryRow(id)
	err := surveyRow.Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference, &survey.LegalBasisRef, &survey.LegalBasis)

	if err == sql.ErrNoRows {
		re := NewRESTError("404", "Survey not found")
		data, err := json.Marshal(re)
		if err != nil {
			http.Error(w, "Error marshaling NewRestError JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write(data)

		return
	}

	if err != nil {
		http.Error(w, "get survey query failed", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(survey)
	if err != nil {
		http.Error(w, "Failed to marshal survey JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetSurveyByShortName returns the details of the survey identified by the string shortName.
func (api *API) GetSurveyByShortName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["shortName"]
	surveyRow := api.GetSurveyByShortNameStmt.QueryRow(id)
	survey := new(Survey)
	err := surveyRow.Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference, &survey.LegalBasisRef, &survey.LegalBasis)

	if err == sql.ErrNoRows {
		re := NewRESTError("404", "Survey not found")
		data, err := json.Marshal(re)
		if err != nil {
			http.Error(w, "Error marshaling NewRestError JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write(data)

		return
	}

	if err != nil {
		http.Error(w, "get survey by shortname query failed", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(survey)
	if err != nil {
		http.Error(w, "Failed to marshal survey JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

// GetSurveyByReference returns the details of the survey identified by the string ref.
func (api *API) GetSurveyByReference(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ref"]
	surveyRow := api.GetSurveyByReferenceStmt.QueryRow(id)
	survey := new(Survey)
	err := surveyRow.Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference, &survey.LegalBasisRef, &survey.LegalBasis)

	if err == sql.ErrNoRows {
		re := NewRESTError("404", "Survey not found")
		data, err := json.Marshal(re)
		if err != nil {
			http.Error(w, "Error marshaling NewRestError JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write(data)

		return
	}

	if err != nil {
		http.Error(w, "get survey by reference query failed", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(survey)
	if err != nil {
		http.Error(w, "Failed to marshal survey JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

// AllClassifierTypeSelectors returns all the classifier type selectors for the survey identified by the string surveyID. The classifier type selectors are returned in ascending order.
func (api *API) AllClassifierTypeSelectors(w http.ResponseWriter, r *http.Request) {
	// We need to run a query first to check if the survey exists so an HTTP 404 can be correctly
	// returned if it doesn't exist. Without this check an HTTP 204 is incorrectly returned for an
	// invalid survey ID.
	vars := mux.Vars(r)
	surveyID := vars["surveyId"]

	err := api.getSurveyID(surveyID)

	if err == sql.ErrNoRows {
		re := NewRESTError("404", "Survey not found")
		data, err := json.Marshal(re)
		if err != nil {
			http.Error(w, "Error marshaling NewRestError JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write(data)

		return
	}

	if err != nil {
		http.Error(w, "Error getting list of classifier type selectors for survey '"+surveyID+"' - "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Now we can get the classifier type selector records.
	rows, err := api.GetClassifierTypeSelectorStmt.Query(surveyID)

	if err != nil {
		http.Error(w, "Error getting list of classifier type selectors for survey '"+surveyID+"' - "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	classifierTypeSelectorSummaries := make([]*ClassifierTypeSelectorSummary, 0)

	for rows.Next() {
		classifierTypeSelectorSummary := new(ClassifierTypeSelectorSummary)
		err = rows.Scan(&classifierTypeSelectorSummary.ID, &classifierTypeSelectorSummary.Name)

		if err != nil {
			//LogError("Error getting list of classifier type selectors for survey '"+surveyID+"'", err)
			http.Error(w, "Error getting list of classifier type selectors for survey '"+surveyID+"' - "+err.Error(), http.StatusInternalServerError)
			return
		}

		classifierTypeSelectorSummaries = append(classifierTypeSelectorSummaries, classifierTypeSelectorSummary)
	}

	data, err := json.Marshal(classifierTypeSelectorSummaries)
	if err != nil {
		http.Error(w, "Failed to marshal classifier type selector summary JSON", http.StatusInternalServerError)
		return
	}

	if len(classifierTypeSelectorSummaries) == 0 {
		http.Error(w, "No classifier type selectors found", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetClassifierTypeSelectorByID returns the details of the classifier type selector for the survey identified by the string surveyID and
// the classifier type selector identified by the str	ing classifierTypeSelectorID.
func (api *API) GetClassifierTypeSelectorByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	surveyID := vars["surveyId"]
	classifierTypeSelectorID := vars["classifierTypeSelectorId"]

	err := api.getSurveyID(surveyID)

	if err == sql.ErrNoRows {
		re := NewRESTError("404", "Classifier Type Selector not found")
		data, err := json.Marshal(re)
		if err != nil {
			http.Error(w, "Error marshaling NewRestError JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write(data)

		return
	}

	if err != nil {
		http.Error(w, "Error getting classifier type selector '"+classifierTypeSelectorID+"' for survey '"+surveyID+"' - "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Now we can get the classifier type selector and classifier type records.
	classifierRows, err := api.GetClassifierTypeSelectorByIDStmt.Query(classifierTypeSelectorID)

	if err != nil {
		http.Error(w, "Get classifiers query failed", http.StatusInternalServerError)

	}
	classifierTypeSelector := new(ClassifierTypeSelector)

	// Using make here ensures the JSON contains an empty array if there are no classifier
	// types, rather than null.
	classifierTypes := make([]string, 0)
	var classifierType string

	for classifierRows.Next() {
		err = classifierRows.Scan(&classifierTypeSelector.ID, &classifierTypeSelector.Name, &classifierType)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Get classifier type by id query failed", http.StatusInternalServerError)
			return
		}

		classifierTypes = append(classifierTypes, classifierType)
	}

	if len(classifierTypes) == 0 {
		re := NewRESTError("404", "Classifier Type Selector not found")
		data, err := json.Marshal(re)
		if err != nil {
			http.Error(w, "Error marshaling NewRestError JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write(data)

		return
	}

	classifierTypeSelector.ClassifierTypes = classifierTypes

	data, err := json.Marshal(classifierTypeSelector)
	if err != nil {
		http.Error(w, "Failed to marshal classifier type JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func (api *API) getSurveyID(surveyID string) error {
	var id string
	return api.GetSurveyIDStmt.QueryRow(surveyID).Scan(&id)
}

func (api *API) getSurveyRef(surveyRef string) error {
	var surveyref string
	return api.GetSurveyRefStmt.QueryRow(surveyRef).Scan(&surveyref)
}

// This function returns the legal basis for a given legal basis longname
func (api *API) getLegalBasisFromLongName(longName string) (LegalBasis, error) {
	var legalBasis LegalBasis;
	err := api.GetLegalBasisFromLongNameStmt.QueryRow(longName).Scan(&legalBasis.Reference, &legalBasis.LongName)

	return legalBasis, err
}

// This function returns the legal basis for a given legal basis ref
func (api *API) getLegalBasisFromRef(ref string) (LegalBasis, error) {
	var legalBasis LegalBasis;
	err := api.GetLegalBasisFromRefStmt.QueryRow(ref).Scan(&legalBasis.Reference, &legalBasis.LongName)

	return legalBasis, err
}

func createStmt(sqlStatement string, db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare(sqlStatement)
}

//Close closes all db connections on the api struct
func (api *API) Close() {
	api.AllSurveysStmt.Close()
}
