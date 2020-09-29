package models

// TODO All transactions should have rollbacks in the case of errors
// TODO Fix the multiple types of returning errors to the client (should be using the correct RFC response)

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"os"
	"time"

	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	validator2 "gopkg.in/go-playground/validator.v9"
)

// ClassifierTypeSelectorSummary represents a summary of a classifier type selector.
type ClassifierTypeSelectorSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ClassifierTypeSelector represents the detail of a classifier type selector.
type ClassifierTypeSelector struct {
	ID              string   `json:"id"`
	Name            string   `json:"name" validate:"required,min=1,max=50,no-spaces"`
	ClassifierTypes []string `json:"classifierTypes" validate:"required,min=1,dive,min=1,max=50,no-spaces"`
}

// Survey represents the details of a survey.
type Survey struct {
	ID            string                   `json:"id"`
	ShortName     string                   `json:"shortName" validate:"required,no-spaces,max=20"`
	LongName      string                   `json:"longName" validate:"required,max=100"`
	Reference     string                   `json:"surveyRef" validate:"required,max=20"`
	LegalBasis    string                   `json:"legalBasis"`
	SurveyType    string                   `json:"surveyType"`
	SurveyMode    string                   `json:"surveyMode"`
	LegalBasisRef string                   `json:"legalBasisRef"`
	Classifiers   []ClassifierTypeSelector `json:"classifiers,omitempty"`
}

// LegalBasis - the legal basis for a survey consisting of a short reference and a long name
type LegalBasis struct {
	Reference string `json:"ref"`
	LongName  string `json:"longName"`
}

//API contains all the pre-prepared sql statements
type API struct {
	AllSurveysStmt                         *sql.Stmt
	GetSurveysBySurveyTypeStmt             *sql.Stmt
	GetSurveyStmt                          *sql.Stmt
	GetSurveyByShortNameStmt               *sql.Stmt
	GetSurveyByReferenceStmt               *sql.Stmt
	GetSurveyIDStmt                        *sql.Stmt
	GetClassifierTypeSelectorStmt          *sql.Stmt
	GetClassifierTypeSelectorByIDStmt      *sql.Stmt
	GetSurveyRefStmt                       *sql.Stmt
	PutSurveyDetailsBySurveyRefStmt        *sql.Stmt
	CreateSurveyStmt                       *sql.Stmt
	CreateSurveyClassifierTypeSelectorStmt *sql.Stmt
	CreateSurveyClassifierTypeStmt         *sql.Stmt
	GetLegalBasesStmt                      *sql.Stmt
	GetLegalBasisFromLongNameStmt          *sql.Stmt
	GetLegalBasisFromRefStmt               *sql.Stmt
	GetSurveyByShortnameStmt               *sql.Stmt
	GetSurveyPKByID                        *sql.Stmt
	CountMatchingClassifierTypeSelectors   *sql.Stmt
	Validator                              *validator2.Validate
	DB                                     *sql.DB
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

// SetUpRoutes balh
func SetUpRoutes(r *mux.Router, api *API) {
	r.HandleFunc("/info", api.Info).Methods("GET")
	r.HandleFunc("/surveys", use(api.AllSurveys, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/surveytype/{surveyType}", use(api.SurveysByType, basicAuth)).Methods("GET")
	r.HandleFunc("/legal-bases", use(api.AllLegalBases, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/{surveyId}", use(api.GetSurvey, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/shortname/{shortName}", use(api.GetSurveyByShortName, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/ref/{ref}", use(api.PutSurveyDetails, basicAuth)).Methods("PUT")
	r.HandleFunc("/surveys", use(api.PostSurveyDetails, basicAuth)).Methods("POST")
	r.HandleFunc("/surveys/ref/{ref}", use(api.GetSurveyByReference, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/{surveyId}/classifiertypeselectors", use(api.AllClassifierTypeSelectors, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/{surveyId}/classifiertypeselectors/{classifierTypeSelectorId}", use(api.GetClassifierTypeSelectorByID, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/{surveyId}/classifiers", use(api.PostSurveyClassifiers, basicAuth)).Methods("POST")
}

//NewAPI returns an API struct populated with all the created SQL statements
func NewAPI(db *sql.DB) (*API, error) {
	allSurveyStmt, err := createStmt("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, s.surveytype, s.surveymode, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref ORDER BY shortname ASC", db)
	if err != nil {
		return nil, err
	}

	getSurveysBySurveyTypeStmt, err := createStmt("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, s.surveytype, s.surveymode, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref WHERE s.surveyType = $1 ORDER BY shortname ASC", db)
	if err != nil {
		return nil, err
	}

	getSurveyStmt, err := createStmt("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, s.surveytype, s.surveymode, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref WHERE id = $1", db)
	if err != nil {
		return nil, err
	}

	getSurveyByShortNameStmt, err := createStmt("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, s.surveytype, s.surveymode, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref  WHERE LOWER(shortName) = LOWER($1)", db)
	if err != nil {
		return nil, err
	}

	getSurveyByReferenceStmt, err := createStmt("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, s.surveytype, s.surveymode, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref  WHERE LOWER(surveyref) = LOWER($1)", db)
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

	createSurvey, err := createStmt("INSERT INTO survey.survey ( surveypk, id, surveyref, shortname, longname, legalbasis, surveytype, surveymode ) VALUES ( nextval('survey.survey_surveypk_seq'), $1, $2, $3, $4, $5, $6, $7) RETURNING surveypk", db)
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

	getSurveyByShortname, err := createStmt("SELECT surveyref FROM survey.survey WHERE shortname = $1", db)
	if err != nil {
		return nil, err
	}

	createSurveyClassifierTypeSelectorStmt, err := createStmt("INSERT INTO survey.classifiertypeselector ( classifiertypeselectorpk, id, surveyfk, classifiertypeselector ) VALUES ( nextval('survey.classifiertypeselector_classifiertypeselectorpk_seq'), $1, $2, $3 ) RETURNING classifiertypeselectorpk as id", db)
	if err != nil {
		return nil, err
	}

	createSurveyClassifierTypeStmt, err := createStmt("INSERT INTO survey.classifiertype ( classifiertypepk, classifiertypeselectorfk, classifiertype ) VALUES ( nextval('survey.classifiertype_classifiertypepk_seq'), $1, $2 )", db)
	if err != nil {
		return nil, err
	}

	getSurveyPKByID, err := createStmt("SELECT surveypk FROM survey.survey WHERE id = $1", db)
	if err != nil {
		return nil, err
	}

	countMatchingClassifierTypeSelectorStmt, err := createStmt("SELECT COUNT(classifiertypeselector.id) FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.surveyfk = survey.surveypk WHERE survey.id = $1 AND classifiertypeselector.classifiertypeselector = $2", db)
	if err != nil {
		return nil, err
	}

	validator := createValidator()

	return &API{
			AllSurveysStmt:                         allSurveyStmt,
			GetSurveysBySurveyTypeStmt:             getSurveysBySurveyTypeStmt,
			GetSurveyStmt:                          getSurveyStmt,
			GetSurveyByShortNameStmt:               getSurveyByShortNameStmt,
			GetSurveyByReferenceStmt:               getSurveyByReferenceStmt,
			GetSurveyIDStmt:                        getSurveyIDStmt,
			GetClassifierTypeSelectorStmt:          getClassifierTypeSelectorStmt,
			GetClassifierTypeSelectorByIDStmt:      getClassifierTypeSelectorByIDStmt,
			GetSurveyRefStmt:                       getSurveyRefStmt,
			PutSurveyDetailsBySurveyRefStmt:        putSurveyDetailsBySurveyRefStmt,
			CreateSurveyStmt:                       createSurvey,
			CreateSurveyClassifierTypeSelectorStmt: createSurveyClassifierTypeSelectorStmt,
			CreateSurveyClassifierTypeStmt:         createSurveyClassifierTypeStmt,
			GetLegalBasesStmt:                      getLegalBases,
			GetLegalBasisFromLongNameStmt:          getLegalBasisFromLongName,
			GetLegalBasisFromRefStmt:               getLegalBasisFromRef,
			GetSurveyByShortnameStmt:               getSurveyByShortname,
			GetSurveyPKByID:                        getSurveyPKByID,
			CountMatchingClassifierTypeSelectors:   countMatchingClassifierTypeSelectorStmt,
			Validator:                              validator,
			DB:                                     db},
		nil
}

func stripChars(str string, fn runevalidator) string {
	return strings.Map(func(r rune) rune {
		if fn(r) {
			return -1
		}

		return r
	}, str)
}

type runevalidator func(rune) bool

// Could use validator:excludeall but would need to enumerate all space values.  This way we can leverage
// unicode.IsSpace
func validateNoSpaces(fl validator2.FieldLevel) bool {
	str := fl.Field().String()
	stripped := stripChars(str, unicode.IsSpace)

	return str == stripped
}

func createValidator() *validator2.Validate {
	validator := validator2.New()

	validator.RegisterValidation("no-spaces", validateNoSpaces)

	return validator
}

// PostSurveyDetails endpoint handler - creates a new survey based on JSON in request
func (api *API) PostSurveyDetails(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	var survey Survey
	err = json.Unmarshal(body, &survey)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	var legalBasis LegalBasis
	var errorMessage string

	if survey.LegalBasisRef != "" {
		legalBasis, err = api.getLegalBasisFromRef(survey.LegalBasisRef)
		errorMessage = fmt.Sprintf("Legal basis with reference %v does not exist", survey.LegalBasisRef)
	} else if survey.LegalBasis != "" {
		legalBasis, err = api.getLegalBasisFromLongName(survey.LegalBasis)
		errorMessage = fmt.Sprintf("Legal basis %v does not exist", survey.LegalBasis)
	} else {
		http.Error(w, "No legal basis specified for survey", http.StatusBadRequest)
		return
	}

	validSurveyTypes := map[string]bool{"Census": true, "Business": true, "Social": true}
	if _, ok := validSurveyTypes[survey.SurveyType]; !ok {
		http.Error(w, "Survey type must be one of [Census, Business, Social]", http.StatusBadRequest)
		return
	}

	if err == sql.ErrNoRows {
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Error getting legal basis - %v", err), http.StatusInternalServerError)
		return
	}

	// Generate a UUID to uniquely identify the new survey
	// TODO replace with Google uuid generator?
	surveyID, err := uuid.NewV4()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating random uuid"), http.StatusInternalServerError)
		return
	}
	// Check whether the survey shortname already exists in the database, returning
	// an error if it does
	_, err = api.getSurveyByShortname(survey.ShortName)
	if (err != nil && err != sql.ErrNoRows) || err == nil {
		http.Error(w, fmt.Sprintf("The survey with Abbreviation %v already exists", survey.ShortName), http.StatusConflict)
		return
	}

	// If the reference is unique then we can go ahead and create the survey and
	// it's classifiers
	if err = api.getSurveyRef(survey.Reference); err == sql.ErrNoRows {

		if err := api.Validator.Struct(survey); err != nil {
			http.Error(w, fmt.Sprintf("Survey failed to validate - %v", err), http.StatusBadRequest)
			return
		}

		surveyPK := 0
		err := api.CreateSurveyStmt.QueryRow(
			surveyID,
			survey.Reference,
			survey.ShortName,
			survey.LongName,
			legalBasis.Reference,
			survey.SurveyType,
			survey.SurveyMode,
		).Scan(&surveyPK)
		if err != nil {
			http.Error(w, fmt.Sprintf("Create survey details failed - %v", err), http.StatusInternalServerError)
			return
		}

		// If the main survey record has been correctly created then, if a set of
		// classifiers have been supplied, we want to create them.
		if survey.Classifiers != nil {
			for _, c := range survey.Classifiers {
				_, err := api.createClassifiers(int(surveyPK), surveyID.String(), c.Name, c.ClassifierTypes)
				if err != nil {
					logErrorAndRespond(w, "Failed to insert classifier '"+c.Name+"'", http.StatusInternalServerError, err)
					return
				}
			}
		}

		// Update the data passed in with the generated values so we can return them
		// to the caller
		survey.ID = surveyID.String()
		survey.LegalBasisRef = legalBasis.Reference
		survey.LegalBasis = legalBasis.LongName

		var js []byte
		js, err = json.Marshal(&survey)

		logger.Info("New survey created",
			zap.String("service", serviceName),
			zap.String("event", "created survey"),
			zap.String("survey_id", survey.ID),
			zap.String("survey_name", survey.LongName),
			zap.String("survey_type", survey.SurveyType),
			zap.String("survey_mode", survey.SurveyMode),
			zap.String("created", time.Now().UTC().Format(timeFormat)))

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write(js)
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Failed to validate survey ref - %v", err), http.StatusInternalServerError)
		return
	} else {
		http.Error(w, fmt.Sprintf("Survey with reference %v already exists", survey.Reference), http.StatusConflict)
		return
	}
}

// Insert a list of classifier types into the database given type selector primary key using transaction tx
func (api *API) insertClassifierTypes(classifierTypes []string, typeSelectorPK int, tx *sql.Tx) error {
	txCreateSurveyClassifierTypeStmt := tx.Stmt(api.CreateSurveyClassifierTypeStmt)
	for _, classifierType := range classifierTypes {
		_, err := txCreateSurveyClassifierTypeStmt.Exec(typeSelectorPK, classifierType)
		if err != nil {
			rollBack(tx)
			return err
		}
	}
	return nil
}

// Insert a classifier type selector into the database given survey primary key using transaction tx, return the PK and UUID
func (api *API) insertClassifierTypeSelector(name string, surveyPK int, tx *sql.Tx) (int, uuid.UUID, error) {
	txCreateSurveyClassifierTypeSelectorStmt := tx.Stmt(api.CreateSurveyClassifierTypeSelectorStmt)
	var typeSelectorPK int
	classifierTypeSelectorID, err := uuid.NewV4()
	if err != nil {
		tx.Rollback()
		return typeSelectorPK, uuid.Nil, errors.New("Error generating random uuid")
	}
	err = txCreateSurveyClassifierTypeSelectorStmt.
		QueryRow(classifierTypeSelectorID, surveyPK, name).
		Scan(&typeSelectorPK)
	if err != nil {
		rollBack(tx)
	}
	return typeSelectorPK, classifierTypeSelectorID, err
}

// Return a boolean true if a classifier type selector exists for the given survey ID and name
func (api *API) classifierTypeSelectorExists(name string, surveyID string) (bool, error) {
	var classifierMatchCount int
	err := api.CountMatchingClassifierTypeSelectors.
		QueryRow(surveyID, name).
		Scan(&classifierMatchCount)
	return classifierMatchCount > 0, err
}

// PostSurveyClassifiers endpoint handler - creates a new survey classifier
func (api *API) PostSurveyClassifiers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	for _, u := range []string{"surveyId"} {
		val, ok := vars[u]
		if !ok {
			http.Error(w, "Missing value for "+u, http.StatusBadRequest)
			return
		}

		if _, err := uuid.FromString(val); err != nil {
			http.Error(w, "The value ("+val+") used for "+u+" is not a valid UUID", http.StatusBadRequest)
			return
		}
	}

	surveyID := vars["surveyId"]

	// Check survey exists and get it's PK
	surveyPK, err := api.getSurveyPKByID(surveyID)

	if err == sql.ErrNoRows {
		writeRestErrorResponse(w, "Survey not found for ID '"+surveyID+"'", http.StatusNotFound)
		return
	}

	if err != nil {
		logErrorAndRespond(w, "Error retrieving survey by survey ID", http.StatusInternalServerError, err)
		return
	}

	// Read and unmarshal request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logErrorAndRespond(w, "Error creating classifier type selector for survey", http.StatusInternalServerError, err)
		return
	}
	var postData ClassifierTypeSelector
	err = json.Unmarshal(body, &postData)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}
	err = api.Validator.Struct(postData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	classifierID, err := api.createClassifiers(surveyPK, surveyID, postData.Name, postData.ClassifierTypes)
	if err != nil {
		logErrorAndRespond(w, "Failed to create classifiers", http.StatusInternalServerError, err)
		return
	}

	// Add inserted classifier to response object
	createdClassifier := postData
	createdClassifier.ID = classifierID

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdClassifier); err != nil {
		logError("Error encoding response to 'post survey classifiers'", err)
	}
}

func (api *API) createClassifiers(surveyPK int, surveyID, name string, types []string) (string, error) {
	logger.Info("Creating classifiers", zap.String("surveyID", surveyID))
	// Check if classifier type selector already exists
	classifierTypeSelectorAlreadyExists, err := api.classifierTypeSelectorExists(name, surveyID)
	if err != nil {
		return "", errors.Wrap(err, "Error counting existing classifier type selectors")
	}
	if classifierTypeSelectorAlreadyExists {
		return "", errors.New(fmt.Sprintf("Type selector with name '%s' already exists for this survey with ID '%s'", name, surveyID))
	}

	// Start database transaction
	tx, err := api.DB.Begin()
	if err != nil {
		return "", errors.Wrap(err, "Error creating database transaction")
	}

	// Insert classifier type selector and retrieve its primary key so that we can
	// use that to associate the classifier types
	typeSelectorPK, classifierTypeSelectorID, err := api.insertClassifierTypeSelector(name, surveyPK, tx)
	if err != nil {
		tx.Rollback()
		return "", errors.Wrap(err, "Error fetching type selector primary key")
	}

	// Insert classifier types
	err = api.insertClassifierTypes(types, typeSelectorPK, tx)
	if err != nil {
		tx.Rollback()
		return "", errors.Wrap(err, "Error inserting classifier types")
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return "", errors.Wrap(err, "Error committing transaction for posting survey classifier")
	}
	logger.Info("Finished creating classifiers", zap.String("surveyID", surveyID))
	return classifierTypeSelectorID.String(), nil
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
	logger.Info("Getting info", zap.String("url", r.URL.Path))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(NewVersion()); err != nil {
		http.Error(w, "info encoding failed", http.StatusInternalServerError)
	}
}

// AllSurveys returns a list of all known surveys
func (api *API) AllSurveys(w http.ResponseWriter, r *http.Request) {
	logger.Info("Getting AllSurveys", zap.String("url", r.URL.Path))
	var rows *sql.Rows
	var err error
	rows, err = api.AllSurveysStmt.Query()
	if err != nil {
		logError("Get all surveys returned error", err)
		http.Error(w, "Failed to retrieve surveys", http.StatusInternalServerError)
		return
	}
	parseSurveys(rows, w)
}

// SurveysByType returns surveys of a particular type
func (api *API) SurveysByType(w http.ResponseWriter, r *http.Request) {
	logger.Info("Getting SurveysByType", zap.String("url", r.URL.Path))
	var rows *sql.Rows
	var err error
	var surveyMap = map[string]string{
		"business": "Business",
		"social":   "Social",
		"census":   "Census",
	}
	vars := mux.Vars(r)
	surveyType := strings.ToLower(vars["surveyType"])

	if mappedSurveyType, ok := surveyMap[surveyType]; ok {

		rows, err = api.GetSurveysBySurveyTypeStmt.Query(mappedSurveyType)
		if err != nil {
			logError("Get surveys by type returned error", err)
			http.Error(w, "Failed to retrieve surveys", http.StatusInternalServerError)
			return
		}
		parseSurveys(rows, w)
		return
	}
	logError("Invalid surveyType in SurveysByType", fmt.Errorf("surveyType:%s", surveyType))
	http.Error(w, "Failed to retrieve surveys", http.StatusBadRequest)
}

func parseSurveys(rows *sql.Rows, w http.ResponseWriter) {
	var err error
	defer rows.Close()
	surveys := make([]*Survey, 0)

	for rows.Next() {
		survey := new(Survey)
		err = rows.Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference, &survey.LegalBasisRef, &survey.SurveyType, &survey.SurveyMode, &survey.LegalBasis)

		if err != nil {
			logError("Failed to get surveys from database", err)
			http.Error(w, "Failed to get surveys from database", http.StatusInternalServerError)
			return
		}

		surveys = append(surveys, survey)
	}

	if len(surveys) == 0 {
		logError("No surveys found", errors.New("no content"))
		http.Error(w, "No surveys found", http.StatusNoContent)
		return
	}

	data, err := json.Marshal(surveys)
	if err != nil {
		logError("Failed to marshal survey summary JSON", err)
		http.Error(w, "Failed to marshal survey summary JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// AllLegalBases returns details of all legal bases
func (api *API) AllLegalBases(w http.ResponseWriter, r *http.Request) {
	logger.Info("Getting AllLegalBases", zap.String("url", r.URL.Path))
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
	logger.Info("Getting Survey", zap.String("url", r.URL.Path))
	vars := mux.Vars(r)
	id := vars["surveyId"]
	survey := new(Survey)
	surveyRow := api.GetSurveyStmt.QueryRow(id)
	err := surveyRow.Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference, &survey.LegalBasisRef, &survey.SurveyType, &survey.SurveyMode, &survey.LegalBasis)

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
	logger.Info("Getting SurveyByShortName", zap.String("url", r.URL.Path))
	vars := mux.Vars(r)
	id := vars["shortName"]

	surveyRow := api.GetSurveyByShortNameStmt.QueryRow(id)

	survey := new(Survey)
	err := surveyRow.Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference, &survey.LegalBasisRef, &survey.SurveyType, &survey.SurveyMode, &survey.LegalBasis)

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
	logger.Info("Getting SurveyByReference", zap.String("url", r.URL.Path))
	vars := mux.Vars(r)
	id := vars["ref"]

	surveyRow := api.GetSurveyByReferenceStmt.QueryRow(id)
	survey := new(Survey)
	err := surveyRow.Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference, &survey.LegalBasisRef, &survey.SurveyType, &survey.SurveyMode, &survey.LegalBasis)

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
	logger.Info("Getting AllClassifierTypeSelectors", zap.String("url", r.URL.Path))
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
// the classifier type selector identified by the string classifierTypeSelectorID.
func (api *API) GetClassifierTypeSelectorByID(w http.ResponseWriter, r *http.Request) {
	logger.Info("Getting ClassifierTypeSelectorByID", zap.String("url", r.URL.Path))
	vars := mux.Vars(r)
	for _, u := range []string{"classifierTypeSelectorId", "surveyId"} {
		val, ok := vars[u]
		if !ok {
			http.Error(w, "Missing value for "+u, http.StatusBadRequest)
		}

		if _, err := uuid.FromString(val); err != nil {
			http.Error(w, "The value ("+val+") used for "+u+" is not a valid UUID", http.StatusBadRequest)
			return
		}
	}
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

// This function returns the survey ref for the survey with the given shortname
func (api *API) getSurveyByShortname(shortname string) (string, error) {
	var lShortname string
	err := api.GetSurveyByShortnameStmt.QueryRow(shortname).Scan(&lShortname)

	return lShortname, err
}

// This function returns the legal basis for a given legal basis longname
func (api *API) getLegalBasisFromLongName(longName string) (LegalBasis, error) {
	var legalBasis LegalBasis
	err := api.GetLegalBasisFromLongNameStmt.QueryRow(longName).Scan(&legalBasis.Reference, &legalBasis.LongName)

	return legalBasis, err
}

// This function returns the legal basis for a given legal basis ref
func (api *API) getLegalBasisFromRef(ref string) (LegalBasis, error) {
	var legalBasis LegalBasis
	err := api.GetLegalBasisFromRefStmt.QueryRow(ref).Scan(&legalBasis.Reference, &legalBasis.LongName)

	return legalBasis, err
}

// Get a survey primary key by UUID string
func (api *API) getSurveyPKByID(surveyID string) (int, error) {
	var surveyPK int
	err := api.GetSurveyPKByID.QueryRow(surveyID).Scan(&surveyPK)
	return surveyPK, err

}

func createStmt(sqlStatement string, db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare(sqlStatement)
}

//Close closes all db connections on the api struct
func (api *API) Close() {
	api.AllSurveysStmt.Close()
}

// Roll back a given transaction and log any errors which occur
func rollBack(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		logError("Error rolling back database transaction", err)
	}
}

// Log and message and an error and send an HTTP response
func logErrorAndRespond(w http.ResponseWriter, logMessage string, status int, err error) {
	logError(logMessage, err)
	http.Error(w, "Internal Server Error", status)
}

// Writes a NewRESTError and sends an HTTP response
func writeRestErrorResponse(w http.ResponseWriter, message string, status int) {
	response := NewRESTError(strconv.Itoa(status), message)
	data, err := json.Marshal(response)
	if err != nil {
		logErrorAndRespond(w, "Error marshalling NewRestError JSON", http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	w.Write(data)
}
