package models_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/ONSdigital/rm-survey-service/models"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

// Examples of valid data to use
const shortName = "test-shortname"
const longName = "test-longname"
const legalBasisLongName = "test-legalbasis-longname"
const reference = "test-reference"
const surveyType = "Business"
const surveyID = "67602ba2-8af6-4298-af66-4e46a62f32c8"
const classifierID = "c0482274-9e96-4001-8797-4b487454c187"
const surveyMode = "SEFT"
const eQVersion = "v2"

var httpClient = &http.Client{}

func TestInfoEndpoint(t *testing.T) {
	Convey("Info enpoint returns a 200 response", t, func() {
		db, mock, err := sqlmock.New()
		prepareMockStmts(mock)
		So(err, ShouldBeNil)
		db.Begin()
		defer db.Close()
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		r, err := http.NewRequest("GET", "http://localhost:9090/health", nil)
		So(err, ShouldBeNil)
		w := httptest.NewRecorder()
		So(err, ShouldBeNil)
		api.Info(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

func TestSurveyListReturnsJson(t *testing.T) {
	Convey("Surveys list returns an array of surveys", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "short_name", "long_name", "survey_ref", "legal_basis", "survey_type", "survey_mode", "long_name"}).AddRow(surveyID, shortName, longName, reference, "test-legalbasis-ref", "test-surveytype", surveyMode, legalBasisLongName)
		mock.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref").ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusOK)
		expected := []models.Survey{{ID: surveyID, ShortName: shortName}}
		res := []models.Survey{}
		body, err := io.ReadAll(resp.Body)
		json.Unmarshal(body, &res)
		So(res[0].ID, ShouldEqual, expected[0].ID)
		So(res[0].ShortName, ShouldEqual, expected[0].ShortName)
	})
}

func TestSurveyListInternalServerError(t *testing.T) {
	Convey("Surveys list returns a 500", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT id, short_name, long_name, survey_ref, survey_type, legal_basis, survey_mode FROM survey.survey").ExpectQuery().WillReturnError(fmt.Errorf("Testing internal server error"))
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusInternalServerError)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Failed to retrieve surveys")
	})
}

func TestSurveyListNotFound(t *testing.T) {
	Convey("Surveys list returns an 500 not found", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "short_name", "long_name", "survey_ref", "legal_basis", "survey_type", "survey_mode", "long_name"})
		mock.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref").ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)
		So(resp.StatusCode, ShouldEqual, http.StatusNoContent)
	})
}

func TestSurveyListBySurveyTypeReturnsJson(t *testing.T) {
	Convey("Surveys list restricted by survey type returns an array of surveys", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "short_name", "long_name", "survey_ref", "legal_basis", "survey_type", "survey_mode", "long_name"}).AddRow("testid", shortName, longName, reference, "test-legalbasis-ref", surveyType, surveyMode, legalBasisLongName)
		mock.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref WHERE s.survey_type =").ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/surveytype/" + surveyType
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusOK)
		expected := []models.Survey{{ID: "testid", SurveyType: "Business"}}
		res := []models.Survey{}
		body, err := io.ReadAll(resp.Body)
		json.Unmarshal(body, &res)
		So(res[0].ID, ShouldEqual, expected[0].ID)
		So(res[0].SurveyType, ShouldEqual, expected[0].SurveyType)
	})
}

func TestSurveyListBySurveyTypeIncorrectCaseReturnsJson(t *testing.T) {
	Convey("Surveys list restricted by survey type of wrong case returns an array of surveys", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "short_name", "long_name", "survey_ref", "legal_basis", "survey_type", "survey_mode", "long_name"}).AddRow(surveyID, shortName, longName, reference, "test-legalbasis-ref", surveyType, surveyMode, legalBasisLongName)
		mock.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref WHERE s.survey_type =").ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/surveytype/BuSiNeSS"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")
		resp, err := httpClient.Do(r)
		So(resp.StatusCode, ShouldEqual, http.StatusOK)
		expected := []models.Survey{{ID: surveyID, SurveyType: surveyType}}
		res := []models.Survey{}
		body, err := io.ReadAll(resp.Body)
		json.Unmarshal(body, &res)
		So(res[0].ID, ShouldEqual, expected[0].ID)
		So(res[0].SurveyType, ShouldEqual, expected[0].SurveyType)
	})
}

func TestSurveyListBySurveyModeEQ(t *testing.T) {
	Convey("Testing to see if survey mode returns what is 'eQ'", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "short_name", "long_name", "survey_ref", "legal_basis", "survey_type", "survey_mode", "long_name"}).AddRow(surveyID, shortName, longName, reference, "test-legalbasis-ref", surveyType, "eQ", legalBasisLongName)
		mock.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref WHERE s.survey_type =").ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/surveytype/Business"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")
		resp, err := httpClient.Do(r)
		So(resp.StatusCode, ShouldEqual, http.StatusOK)
		expected := []models.Survey{{SurveyMode: "eQ"}}
		res := []models.Survey{}
		body, err := io.ReadAll(resp.Body)
		json.Unmarshal(body, &res)
		So(res[0].SurveyMode, ShouldEqual, expected[0].SurveyMode)
	})
}

func TestSurveyListBySurveyModeSEFT(t *testing.T) {
	Convey("Testing to see if survey mode returns what is 'SEFT'", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "short_name", "long_name", "survey_ref", "legal_basis", "survey_type", "survey_mode", "long_name"}).AddRow(surveyID, shortName, longName, reference, "test-legalbasis-ref", surveyType, surveyMode, legalBasisLongName)
		mock.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref WHERE s.survey_type =").ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/surveytype/Business"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusOK)
		expected := []models.Survey{{SurveyMode: "SEFT"}}
		res := []models.Survey{}
		body, err := io.ReadAll(resp.Body)
		json.Unmarshal(body, &res)
		So(res[0].SurveyMode, ShouldEqual, expected[0].SurveyMode)
	})
}

func TestSurveyListBySurveyTypeReturnsErrorForUnknownType(t *testing.T) {
	Convey("Surveys list restricted by survey type returns an array of surveys", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "short_name", "long_name", "survey_ref", "legal_basis", "survey_type", "survey_mode", "long_name"}).AddRow("testid", shortName, longName, reference, "test-legalbasis-ref", surveyType, surveyMode, legalBasisLongName)
		mock.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref WHERE s.surveyType =").ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/surveytype/SomeUnknownType"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, _ := io.ReadAll(resp.Body)
		So(string(body), ShouldEqual, "Failed to retrieve surveys\n")
	})
}

func TestSurveyGetReturnsJson(t *testing.T) {
	Convey("Survey GET returns a survey resource", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "short_name", "long_name", "survey_ref", "legal_basis", "survey_type", "survey_mode", "long_name"}).AddRow(surveyID, shortName, longName, reference, "test-legalbasis-ref", surveyType, surveyMode, legalBasisLongName)
		mock.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusOK)
		expected := models.Survey{ID: surveyID, ShortName: shortName, LongName: longName, Reference: reference, SurveyType: surveyType, SurveyMode: surveyMode}
		res := models.Survey{}
		body, err := io.ReadAll(resp.Body)
		json.Unmarshal(body, &res)
		So(res.ID, ShouldEqual, expected.ID)
		So(res.ShortName, ShouldEqual, expected.ShortName)
		So(res.LongName, ShouldEqual, expected.LongName)
		So(res.Reference, ShouldEqual, expected.Reference)
		So(res.SurveyType, ShouldEqual, expected.SurveyType)
		So(res.SurveyMode, ShouldEqual, expected.SurveyMode)
	})
}

func TestSurveyGetNotFound(t *testing.T) {
	Convey("Survey Get returns an 404 not found", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "short_name", "long_name", "survey_ref", "legal_basis", "survey_type", "survey_mode", "long_name"})
		mock.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusNotFound)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, `{"code":"404","message":"Survey not found",`)
	})
}

func TestSurveyGetInternalServerError(t *testing.T) {
	Convey("Survey GET returns a 500", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT id, short_name, long_name, survey_ref, legal_basis, survey_mode, survey_type from survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnError(fmt.Errorf("Testing internal server error"))
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusInternalServerError)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "get survey query failed")
	})
}

func TestGetSurveyByShortnameReturnsJSON(t *testing.T) {
	Convey("Survey GET by shortname returns a survey resource", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "short_name", "long_name", "survey_ref", "legal_basis", "survey_type", "survey_mode", "eq_version", "long_name"}).AddRow(surveyID, shortName, longName, reference, "test-legalbasis-ref", "test-surveytype", surveyMode, eQVersion, legalBasisLongName)
		mock.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, s.eq_version, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref").ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/shortname/test-shortname"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusOK)
		expected := models.Survey{ID: surveyID, ShortName: shortName, LongName: longName, Reference: reference, SurveyType: "test-surveytype", SurveyMode: surveyMode}
		res := models.Survey{}
		body, err := io.ReadAll(resp.Body)
		json.Unmarshal(body, &res)
		So(res.ID, ShouldEqual, expected.ID)
		So(res.ShortName, ShouldEqual, expected.ShortName)
		So(res.LongName, ShouldEqual, expected.LongName)
		So(res.Reference, ShouldEqual, expected.Reference)
		So(res.SurveyType, ShouldEqual, expected.SurveyType)
		So(res.SurveyMode, ShouldEqual, expected.SurveyMode)
	})
}

func TestSurveyGetByShortNameNotFound(t *testing.T) {
	Convey("Survey Get by shortname returns an 404 not found", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "short_name", "long_name", "survey_ref", "legal_basis", "survey_type", "survey_mode", "long_name"})
		mock.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref").ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/survey/shortname/test-shortname"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusNotFound)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "404 page not found")
	})
}

func TestSurveyGetByShortNameInternalServerError(t *testing.T) {
	Convey("Survey GET by shortname returns a 500", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT id, short_name, long_name, survey_ref, legal_basis, survey_mode from survey.survey").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnError(fmt.Errorf("Testing internal server error"))
		db.Begin()
		defer db.Close()
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/survey/testid", nil)
		So(err, ShouldBeNil)
		api.GetSurveyByShortName(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})
}

// /////
func TestGetSurveyByReferenceReturnsJSON(t *testing.T) {
	Convey("Survey GET by reference returns a survey resource", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "short_name", "long_name", "survey_ref", "legal_basis", "survey_type", "survey_mode", "long_name"}).AddRow(surveyID, shortName, longName, reference, "test-legalbasis-ref", surveyType, surveyMode, legalBasisLongName)
		mock.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref").ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/ref/test-reference"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		c := &http.Client{}
		resp, err := c.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusOK)
		expected := models.Survey{ID: surveyID, ShortName: shortName, LongName: longName, Reference: reference, SurveyType: surveyType, SurveyMode: surveyMode}
		res := models.Survey{}
		body, err := io.ReadAll(resp.Body)
		json.Unmarshal(body, &res)
		So(res.ID, ShouldEqual, expected.ID)
		So(res.ShortName, ShouldEqual, expected.ShortName)
		So(res.LongName, ShouldEqual, expected.LongName)
		So(res.Reference, ShouldEqual, expected.Reference)
		So(res.SurveyType, ShouldEqual, expected.SurveyType)
		So(res.SurveyMode, ShouldEqual, expected.SurveyMode)
	})
}

func TestSurveyGetByReferenceNotFound(t *testing.T) {
	Convey("Survey Get by reference returns an 404 not found", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "short_name", "long_name", "survey_ref", "legal_basis", "survey_type", "survey_mode", "long_name"})
		mock.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref").ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/ref/test-reference"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)
		So(resp.StatusCode, ShouldEqual, http.StatusNotFound)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, `{"code":"404","message":"Survey not found",`)
	})
}

func TestSurveyGetByReferenceInternalServerError(t *testing.T) {
	Convey("Survey GET by reference returns a 500", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT id, short_name, long_name, survey_ref, legal_basis, survey_mode from survey.survey").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnError(fmt.Errorf("Testing internal server error"))
		db.Begin()
		defer db.Close()
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/survey/testid", nil)
		So(err, ShouldBeNil)
		api.GetSurveyByReference(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})
}

func TestAllClassifierTypeSelectorsReturnsJSON(t *testing.T) {
	Convey("ClassifierType GET by reference returns a classifier resource", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		idRow := sqlmock.NewRows([]string{"id"}).AddRow("id").AddRow("id")
		rows := sqlmock.NewRows([]string{"id", "classifiertypeselector"}).AddRow(surveyID, "test-name")
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(idRow)
		mock.ExpectPrepare("SELECT classifiertypeselector.id, classifier_type_selector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.survey_fk = survey.survey_pk WHERE survey.id = .* ORDER BY classifier_type_selector ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiertypeselectors"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusOK)
		expected := models.ClassifierTypeSelectorSummary{ID: surveyID, Name: "test-name"}
		res := []models.ClassifierTypeSelectorSummary{}
		body, err := io.ReadAll(resp.Body)
		json.Unmarshal(body, &res)
		So(res[0].ID, ShouldEqual, expected.ID)
		So(res[0].Name, ShouldEqual, expected.Name)
	})
}

func TestAllClassifierTypeSelectorsSurveyNotFound(t *testing.T) {
	Convey("ClassifierType GET returns a 404 not found", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		idRow := sqlmock.NewRows([]string{"id"})
		rows := sqlmock.NewRows([]string{"id", "classifiertypeselector"}).AddRow(surveyID, "test-name")
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(idRow)
		mock.ExpectPrepare("SELECT classifiertypeselector.id, classifier_type_selector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.survey_fk = survey.survey_pk WHERE survey.id = .* ORDER BY classifier_type_selector ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiertypeselectors"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		c := &http.Client{}
		resp, err := c.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusNotFound)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, `{"code":"404","message":"Survey not found",`)
	})
}

func TestAllClassifierTypeSelectorsNotFound(t *testing.T) {
	Convey("ClassifierType GET returns a 204 no content", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		idRow := sqlmock.NewRows([]string{"id"}).AddRow("test-id")
		rows := sqlmock.NewRows([]string{"id", "classifiertypeselector"})
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(idRow)
		mock.ExpectPrepare("SELECT classifiertypeselector.id, classifier_type_selector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.survey_fk = survey.survey_pk WHERE survey.id = .* ORDER BY classifier_type_selector ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiertypeselectors"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusNoContent)
	})
}

func TestAllClassifierTypeSelectorsSurveyReturnsInternalServerError(t *testing.T) {
	Convey("ClassifierType GET returns a 500 when search of survey fails", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "classifier_type_selector"}).AddRow(surveyID, "test-name")
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnError(fmt.Errorf("Testing internal server error"))
		mock.ExpectPrepare("SELECT classifiertypeselector.id, classifier_type_selector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.survey_fk = survey.surveypk WHERE survey.id = .* ORDER BY classifier_type_selector ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiertypeselectors"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusInternalServerError)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Error getting list of classifier type selectors for survey '"+surveyID+"' - Testing internal server error")
	})
}

func TestAllClassifierTypeSelectorsReturnsInternalServerError(t *testing.T) {
	Convey("ClassifierType GET returns a 500 when classifiertypeselector search fails", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		idRow := sqlmock.NewRows([]string{"id"}).AddRow("id").AddRow("id")
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(idRow)
		mock.ExpectPrepare("SELECT classifiertypeselector.id, classifier_type_selector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.survey_fk = survey.survey_pk WHERE survey.id = .* ").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnError(fmt.Errorf("Testing internal server error"))
		db.Begin()
		defer db.Close()
		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiertypeselectors"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)
		So(resp.StatusCode, ShouldEqual, http.StatusInternalServerError)
		body, err := io.ReadAll(resp.Body)

		So(string(body), ShouldStartWith, "Error getting list of classifier type selectors for survey '"+surveyID+"' - Testing internal server error")
	})
}

func TestClassifierTypeSelectorByIdReturnsJSON(t *testing.T) {
	Convey("ClassifierType GET by ID returns a classifier resource", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		idRow := sqlmock.NewRows([]string{"id"}).AddRow("id").AddRow("id")
		rows := sqlmock.NewRows([]string{"id", "classifier_type_selector", "classifier_type"}).AddRow(surveyID, "test-name", classifierID)
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(idRow)
		mock.ExpectPrepare("SELECT id, classifier_type_selector, classifier_type FROM survey.classifiertype INNER JOIN survey.classifiertypeselector ON classifiertype.classifier_type_selector_fk = classifiertypeselector.classifier_type_selector_pk WHERE classifiertypeselector.id = .* ORDER BY classifier_type ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiertypeselectors/" + classifierID
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusOK)
		var a = []string{"test"}
		expected := models.ClassifierTypeSelector{ID: surveyID, Name: "test-name", ClassifierTypes: a}
		res := models.ClassifierTypeSelector{}
		body, err := io.ReadAll(resp.Body)
		json.Unmarshal(body, &res)
		So(res.ID, ShouldEqual, expected.ID)
		So(res.Name, ShouldEqual, expected.Name)
	})
}

func TestClassifierTypeSelectorByIdSurveyIdIsInvalidUuid(t *testing.T) {
	Convey("ClassifierType GET by ID will return a 400 when supplied with an invalid uuid in the survey_id", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		idRow := sqlmock.NewRows([]string{"id"}).AddRow("id").AddRow("id")
		rows := sqlmock.NewRows([]string{"id", "classifier_type_selector", "classifier_type"}).AddRow(surveyID, "test-name", classifierID)
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(idRow)
		mock.ExpectPrepare("SELECT id, classifier_type_selector, classifier_type FROM survey.classifiertype INNER JOIN survey.classifiertypeselector ON classifiertype.classifier_type_selector_fk = classifiertypeselector.classifier_type_selector_pk WHERE classifiertypeselector.id = .* ORDER BY classifier_type ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/not-a-uuid/classifiertypeselectors/" + classifierID
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "The value (not-a-uuid) used for surveyId is not a valid UUID")
	})
}

func TestClassifierTypeSelectorByIdClassifierIdIsInvalidUuid(t *testing.T) {
	Convey("ClassifierType GET by ID will return a 400 when supplied with an invalid uuid in the survey_id", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		idRow := sqlmock.NewRows([]string{"id"}).AddRow("id").AddRow("id")
		rows := sqlmock.NewRows([]string{"id", "classifier_type_selector", "classifier_type"}).AddRow(surveyID, "test-name", classifierID)
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(idRow)
		mock.ExpectPrepare("SELECT id, classifier_type_selector, classifier_type FROM survey.classifiertype INNER JOIN survey.classifiertypeselector ON classifiertype.classifier_type_selector_fk = classifiertypeselector.classifier_type_selector_pk WHERE classifiertypeselector.id = .* ORDER BY classifier_type ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiertypeselectors/not-a-uuid"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "The value (not-a-uuid) used for classifierTypeSelectorId is not a valid UUID")
	})
}

func TestClassifierTypeSelectorByIdReturns404(t *testing.T) {
	Convey("ClassifierType GET by ID returns a classifier resource", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		idRow := sqlmock.NewRows([]string{"id"})
		rows := sqlmock.NewRows([]string{"id", "classifier_type_selector", "classifier_type"}).AddRow(surveyID, "test-name", "test-type")
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(idRow)
		mock.ExpectPrepare("SELECT id, classifier_type_selector, classifier_type FROM survey.classifiertype INNER JOIN survey.classifiertypeselector ON classifiertype.classifier_type_selector_fk = classifiertypeselector.classifier_type_selector_pk WHERE classifiertypeselector.id = .* ORDER BY classifier_type ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiertypeselectors/bed34d98-f546-40d7-83ba-9ed636f95ac2"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(err, ShouldBeNil)
		So(resp.StatusCode, ShouldEqual, http.StatusNotFound)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, `{"code":"404","message":"Classifier Type Selector not found",`)
	})
}

func TestClassifierTypeSelectorByIdNoClassifierTypesReturns404(t *testing.T) {
	Convey("ClassifierType GET by ID returns 404 if no classifier types exist", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		idRow := sqlmock.NewRows([]string{"id"}).AddRow(surveyID)
		rows := sqlmock.NewRows([]string{"id", "classifier_type_selector", "classifier_type"})
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(idRow)
		mock.ExpectPrepare("SELECT id, classifier_type_selector, classifier_type FROM survey.classifiertype INNER JOIN survey.classifiertypeselector ON classifiertype.classifier_type_selector_fk = classifiertypeselector.classifier_type_selector_pk WHERE classifiertypeselector.id = .* ORDER BY classifier_type ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiertypeselectors/" + classifierID
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusNotFound)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, `{"code":"404","message":"Classifier Type Selector not found",`)
	})
}

func TestClassifierTypeSelectorByIdInternalServerError(t *testing.T) {
	Convey("ClassifierType GET by reference returns a classifier resource", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "classifier_type_selector", "classifier_type"}).AddRow(surveyID, "test-name", "test-type")
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnError(fmt.Errorf("Testing internal server error"))
		mock.ExpectPrepare("SELECT id, classifier_type_selector, classifier_type FROM survey.classifiertype INNER JOIN survey.classifiertypeselector ON classifiertype.classifier_type_selector_fk = classifiertypeselector.classifier_type_selector_pk WHERE classifiertypeselector.id = .* ORDER BY classifier_type ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiertypeselectors/" + classifierID
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("GET", url, nil)
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusInternalServerError)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Error getting classifier type selector '"+classifierID+"' for survey '"+surveyID+"' - Testing internal server error")
	})
}

func TestPutSurveyDetailsBySurveyRefSuccess(t *testing.T) {
	Convey("Survey Details PUT by Survey Reference success", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		refRow := sqlmock.NewRows([]string{"survey_ref"}).AddRow("456")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(refRow)
		mock.ExpectPrepare("UPDATE survey.survey SET short_name = .+, long_name = .+ WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/ref/456"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name"}`)
		r, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusOK)
	})
}

func TestPutSurveyDetailsBySurveyRefInternalServerError(t *testing.T) {
	Convey("Survey Details PUT by Survey Reference success", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectQuery().WillReturnError(fmt.Errorf("Testing internal server error"))
		mock.ExpectPrepare("UPDATE survey.survey SET short_name = .+, long_name = .+ WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		db.Begin()
		defer db.Close()
		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/ref/456"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name"}`)
		r, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)
		So(resp.StatusCode, ShouldEqual, http.StatusInternalServerError)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Failed to get survey ref - Testing internal server error")
	})
}

func TestCreateNewSurvey(t *testing.T) {
	Convey("Create new survey success", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"survey_ref"})
		legalBasis := sqlmock.NewRows([]string{"ref", "long_name"}).AddRow("STA1947", "Statistics of Trade Act 1947")
		newSurveyPK := sqlmock.NewRows([]string{"survey_pk"}).AddRow("1000")

		prepareMockStmts(mock)

		mock.ExpectRollback()
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs("99").WillReturnRows(rows)
		mock.ExpectPrepare("INSERT INTO survey.survey \\( survey_pk, id, survey_ref, short_name, long_name, legal_basis, survey_type, survey_mode, eq_version \\) VALUES \\( .+\\) RETURNING survey_pk").ExpectQuery().WithArgs(sqlmock.AnyArg(), "99", "test-short-name", "test-long-name", "STA1947", "Social", "SEFT", "v2").WillReturnRows(newSurveyPK)
		mock.ExpectPrepare("SELECT ref, long_name FROM survey.legalbasis WHERE long_name = .+").ExpectQuery().WithArgs("Statistics of Trade Act 1947").WillReturnRows(legalBasis)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE short_name = .+").ExpectQuery().WithArgs("test-short-name").WillReturnRows(rows)

		// Insert first classifier with one type
		mock.ExpectBegin()
		mock.ExpectPrepare("SELECT COUNT\\(classifiertypeselector.id\\) FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.survey_fk = survey.survey_pk WHERE survey.id = .+ AND classifiertypeselector.classifier_type_selector = .+").ExpectQuery().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"Count"}).AddRow(0))
		mock.ExpectPrepare("INSERT INTO survey.classifiertypeselector .+").ExpectQuery().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1000"))
		mock.ExpectPrepare("INSERT INTO survey.classifiertype .+").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// Insert second classifier with two types
		mock.ExpectBegin()
		mock.ExpectPrepare("SELECT COUNT\\(classifiertypeselector.id\\) FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.survey_fk = survey.survey_pk WHERE survey.id = .+ AND classifiertypeselector.classifier_type_selector = .+").ExpectQuery().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"Count"}).AddRow(0))
		mock.ExpectPrepare("INSERT INTO survey.classifiertypeselector .+").ExpectQuery().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1000"))
		mock.ExpectPrepare("INSERT INTO survey.classifiertype .+").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectPrepare("INSERT INTO survey.classifiertype .+").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"99","LegalBasis":"Statistics of Trade Act 1947","SurveyType":"Social", "SurveyMode":"SEFT", "EQVersion": "v2"}`)

		r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusCreated)
	})
}

func TestCreateNewSurveyInvalidSurveyType(t *testing.T) {
	Convey("Create new survey success", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref", "longname"}).AddRow("STA1947", "Statistics of Trade Act 1947")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("SELECT ref, long_name FROM survey.legal_basis WHERE long_name = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE short_name = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"99","LegalBasis":"Statistics of Trade Act 1947","SurveyType":"Invalid", "SurveyMode":"SEFT"}`)

		r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldEqual, "Survey type must be one of [Census, Business, Social]\n")
	})
}

func TestCreateNewSurveySurveyTypeDoesNotExist(t *testing.T) {
	Convey("Create new survey success", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref", "longname"}).AddRow("STA1947", "Statistics of Trade Act 1947")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("SELECT ref, long_name FROM survey.legalbasis WHERE long_name = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE short_name = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"99","LegalBasis":"Statistics of Trade Act 1947", "SurveyMode":"SEFT"}`)

		r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldEqual, "Survey type must be one of [Census, Business, Social]\n")
	})
}

func TestCreateNewSurveyNonExistentLegalBasisRef(t *testing.T) {
	Convey("Create new survey with non existent legal basis ref", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref", "longname"})
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("SELECT ref, long_name FROM survey.legalbasis WHERE ref = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		db.Begin()
		defer db.Close()
		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"99","LegalBasisRef":"Statistics of Trade Act 1947", "SurveyType":"Social", "SurveyMode":"SEFT"}`)

		r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Legal basis with reference Statistics of Trade Act 1947 does not exist")
	})
}

func TestCreateNewSurveyNonExistentLegalBasisLongName(t *testing.T) {
	Convey("Create new survey with non existent legal basis ref", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref", "longname"})
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("INSERT INTO survey.survey \\( survey_pk, id, survey_ref, short_name, long_name, legal_basis \\) VALUES \\( .+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectPrepare("SELECT ref, long_name FROM survey.legalbasis WHERE long_name = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"99","LegalBasis":"Statistics of Trade Act 1947", "SurveyType":"Business", "SurveyMode":"SEFT"}`)

		r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Legal basis Statistics of Trade Act 1947 does not exist")
	})
}

func TestCreateNewSurveyNonExistentSurveyModeName(t *testing.T) {
	Convey("Create new survey with non existent legal basis ref", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref", "longname"})
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("INSERT INTO survey.survey \\( survey_pk, id, survey_ref, short_name, long_name, survey_mode, legal_basis \\) VALUES \\( .+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectPrepare("SELECT ref, long_name FROM survey.legalbasis WHERE long_name = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"99","LegalBasis":"Statistics of Trade Act 1947", "SurveyType":"Business", "SurveyMode":"SEFT"}`)

		r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Legal basis Statistics of Trade Act 1947 does not exist")
	})
}

func TestCreateNewSurveyNonExistentLegalBasis(t *testing.T) {
	Convey("Create new survey with non existent legal basis ref", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref", "longname"})
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("SELECT ref, long_name FROM survey.legalbasis WHERE ref = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"99","LegalBasisRef":"STA1947", "SurveyType":"Business", "SurveyMode":"SEFT"}`)

		r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Legal basis with reference STA1947 does not exist")
	})
}

// func TestCreateNewSurveyNonNumericRef(t *testing.T) {
// 	Convey("Create new survey with non numeric refernce", t, func() {
// 		db, mock, err := sqlmock.New()
// 		So(err, ShouldBeNil)
// 		rows := sqlmock.NewRows([]string{"surveyref"})
// 		legalBasis := sqlmock.NewRows([]string{"ref", "longname"}).AddRow("STA1947", "Statistics of Trade Act 1947")
// 		newSurveyPK := sqlmock.NewRows([]string{"surveypk"}).AddRow("1000")
// 		prepareMockStmts(mock)
// 		mock.ExpectPrepare("SELECT surveyref FROM survey.survey WHERE LOWER\\(surveyref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
// 		mock.ExpectPrepare("INSERT INTO survey.survey \\( surveypk, id, surveyref, shortname, longname, legalbasis, surveytype \\) VALUES \\( .+\\) RETURNING surveypk").ExpectQuery().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(newSurveyPK)
// 		mock.ExpectPrepare("SELECT ref, longname FROM survey.legalbasis WHERE ref = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
// 		mock.ExpectPrepare("SELECT surveyref FROM survey.survey WHERE shortname = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
// 		db.Begin()
// 		defer db.Close()

// 		// When
// 		api, err := models.NewAPI(db)
// 		So(err, ShouldBeNil)
// 		defer api.Close()

// 		// Create a new router and plug in the defined routes
// 		router := mux.NewRouter()
// 		models.SetUpRoutes(router, api)

// 		ts := httptest.NewServer(router)
// 		defer ts.Close()
// 		url := ts.URL + "/surveys"
// 		// User and password not set so base64encode the dividing character
// 		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
// 		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"99A","LegalBasisRef":"STA1947", "SurveyType":"Business"}`)

// 		r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
// 		r.Header.Set("Authorization", "Basic: "+basicAuth)
// 		r.Header.Set("Content-Type", "application/json")

// 		resp, err := httpClient.Do(r)

// 		// FIXME This error should throw a 400 status code and different text.  Previously it was giving a 400 because of a missing survey type.
// 		// Since fixing the mock sql statements, it's now giving a 201 as there isn't actually any validation in the Survey struct to stop the
// 		// Reference field having numbers in it.  The validation on this would need to be fixed and then this test amended.
// 		//So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
// 		body, err := io.ReadAll(resp.Body)
// 		So(string(body), ShouldEndWith, `,"shortName":"test-short-name","longName":"test-long-name","surveyRef":"99A","legalBasis":"Statistics of Trade Act 1947","surveyType":"Business","legalBasisRef":"STA1947"}`)
// 		So(resp.StatusCode, ShouldEqual, http.StatusCreated)
// 	})
// }

func TestCreateNewSurveyRefTooLong(t *testing.T) {
	Convey("Create new survey with non numeric refernce", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref", "longname"}).AddRow("STA1947", "Statistics of Trade Act 1947")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("SELECT ref, long_name FROM survey.legalbasis WHERE ref = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE short_name = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"012345678901234567890","LegalBasisRef":"STA1947","SurveyType":"Social", "SurveyMode":"SEFT"}`)

		r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Survey failed to validate - Key: 'Survey.Reference'")
	})
}

func TestCreateNewSurveyShortNameWithSpace(t *testing.T) {
	Convey("Create new survey with non numeric refernce", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref", "longname"}).AddRow("STA1947", "Statistics of Trade Act 1947")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("SELECT ref, long_name FROM survey.legalbasis WHERE ref = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE short_name = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		var jsonStr = []byte(`{"ShortName": "test short name", "LongName":"test-long-name","SurveyRef":"0123","LegalBasisRef":"STA1947","SurveyType":"Social","SurveyMode":"SEFT"}`)

		r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Survey failed to validate - Key: 'Survey.ShortName' Error:Field validation for 'ShortName' failed on the 'no-spaces' tag")
	})
}

func TestCreateNewSurveyShortNameTooLong(t *testing.T) {
	Convey("Create new survey with non numeric refernce", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref", "longname"}).AddRow("STA1947", "Statistics of Trade Act 1947")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("SELECT ref, long_name FROM survey.legalbasis WHERE ref = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE short_name = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		var jsonStr = []byte(`{"ShortName": "test-short-name-0123456", "LongName":"test-long-name","SurveyRef":"0123","LegalBasisRef":"STA1947", "SurveyType":"Business", "SurveyMode":"SEFT"}`)
		r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Survey failed to validate - Key: 'Survey.ShortName' Error:Field validation for 'ShortName' failed on the 'max' tag")
	})
}

func TestCreateNewSurveyLongNameTooLong(t *testing.T) {
	Convey("Create new survey with non numeric refernce", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref", "longname"}).AddRow("STA1947", "Statistics of Trade Act 1947")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("SELECT ref, long_name FROM survey.legalbasis WHERE ref = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE short_name = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name-012345678-012345678-012345678-012345678-012345678-012345678-012345678-01234567899999999-0123456789","SurveyRef":"123","LegalBasisRef":"STA1947", "SurveyType":"Business", "SurveyMode":"SEFT"}`)
		r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Survey failed to validate - Key: 'Survey.LongName' Error:Field validation for 'LongName' failed on the 'max' tag")
	})
}

func TestCreateNewSurveyDupilcateSurveyRef(t *testing.T) {
	Convey("Create new survey with duplicate survey ref", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		surveyRefRows := sqlmock.NewRows([]string{"survey_ref"}).AddRow("0123")
		shortNameRows := sqlmock.NewRows([]string{"short_name"})
		legalBasis := sqlmock.NewRows([]string{"ref", "long_name"}).AddRow("STA1947", "Statistics of Trade Act 1947")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(surveyRefRows)
		mock.ExpectPrepare("INSERT INTO survey.survey \\( survey_pk, id, survey_ref, short_name, long_name, survey_type, legal_basis \\) VALUES \\( .+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectPrepare("SELECT ref, long_name FROM survey.legalbasis WHERE ref = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE short_name = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(shortNameRows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"0123","LegalBasisRef":"STA1947","SurveyType":"Social", "SurveyMode":"SEFT"}`)
		r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusConflict)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Survey with reference 0123 already exists")
	})
}

func TestCreateNewSurveyDupilcateShortName(t *testing.T) {
	Convey("Create new survey with duplicate survey ref", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"survey_ref"}).AddRow("0123")
		noRows := sqlmock.NewRows([]string{"survey_ref"}).AddRow("0123")
		legalBasis := sqlmock.NewRows([]string{"ref", "long_name"}).AddRow("STA1947", "Statistics of Trade Act 1947")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(noRows)
		mock.ExpectPrepare("INSERT INTO survey.survey \\( survey_pk, id, survey_ref, short_name, long_name, survey_type, legal_basis \\) VALUES \\( .+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectPrepare("SELECT ref, long_name FROM survey.legalbasis WHERE ref = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		mock.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE short_name = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"0123","LegalBasisRef":"STA1947","SurveyType":"Social", "SurveyMode":"SEFT"}`)
		r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)

		So(resp.StatusCode, ShouldEqual, http.StatusConflict)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "The survey with Abbreviation test-short-name already exists")
	})
}

func TestCreateNewSurveyClassifiers(t *testing.T) {
	Convey("Create new survey classifiers", t, func() {

		// Given
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		surveyPKRows := sqlmock.NewRows([]string{"surveypk"}).AddRow("1000")
		classifierTypeSelectorMatchesRow := sqlmock.NewRows([]string{"Count"}).AddRow(0)
		classifierTypeSelectorPKRows := sqlmock.NewRows([]string{"id"}).AddRow("1000")
		prepareMockStmts(mock)
		mock.ExpectBegin()
		mock.ExpectPrepare("INSERT INTO survey.classifiertypeselector \\( classifier_type_selector_pk, id, survey_fk, classifier_type_selector \\) VALUES \\( .+, .+, .+, .+ \\) RETURNING classifier_type_selector_pk as id").ExpectQuery().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(classifierTypeSelectorPKRows)
		mock.ExpectPrepare("INSERT INTO survey.classifiertype \\( classifier_type_pk, classifier_type_selector_fk, classifier_type \\) VALUES \\( .+, .+, .+ \\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectPrepare("SELECT survey_pk FROM survey.survey WHERE id = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(surveyPKRows)
		mock.ExpectPrepare("SELECT COUNT\\(classifiertypeselector.id\\) FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.survey_fk = survey.survey_pk WHERE survey.id = .+ AND classifiertypeselector.classifier_type_selector = .+").ExpectQuery().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(classifierTypeSelectorMatchesRow)
		mock.ExpectCommit()
		var postData = []byte(`{"name": "test", "classifierTypes": ["TEST1"]}`)

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiers"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)
		So(err, ShouldBeNil)

		// Then
		So(resp.StatusCode, ShouldEqual, http.StatusCreated)
	})
}

func TestCreateNewSurveyClassifiersInvalidUuid(t *testing.T) {
	Convey("will return a 400 when supplied with an invalid uuid in the survey_id", t, func() {

		// Given
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		surveyPKRows := sqlmock.NewRows([]string{"surveypk"}).AddRow("1000")
		classifierTypeSelectorMatchesRow := sqlmock.NewRows([]string{"Count"}).AddRow(0)
		classifierTypeSelectorPKRows := sqlmock.NewRows([]string{"id"}).AddRow("1000")
		prepareMockStmts(mock)
		mock.ExpectBegin()
		mock.ExpectPrepare("INSERT INTO survey.classifiertypeselector \\( classifier_type_selector_pk, id, survey_fk, classifier_type_selector \\) VALUES \\( .+, .+, .+, .+ \\) RETURNING classifier_type_selector_pk as id").ExpectQuery().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(classifierTypeSelectorPKRows)
		mock.ExpectPrepare("INSERT INTO survey.classifiertype \\( classifier_type_pk, classifier_type_selector_fk, classifier_type \\) VALUES \\( .+, .+, .+ \\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectPrepare("SELECT surveypk FROM survey.survey WHERE id = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(surveyPKRows)
		mock.ExpectPrepare("SELECT COUNT\\(classifiertypeselector.id\\) FROM survey.classifier_type_selector INNER JOIN survey.survey ON classifiertypeselector.survey_fk = survey.surveypk WHERE survey.id = .+ AND classifiertypeselector.classifier_type_selector = .+").ExpectQuery().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(classifierTypeSelectorMatchesRow)
		mock.ExpectCommit()
		var postData = []byte(`{"name": "test", "classifierTypes": ["TEST1"]}`)

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/not-a-uuid/classifiers"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)
		So(err, ShouldBeNil)

		// Then
		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "The value (not-a-uuid) used for surveyId is not a valid UUID")
	})
}

func TestCreateNewSurveyClassifiersSurveyNotFound(t *testing.T) {
	Convey("Create new classifiers survey not found", t, func() {

		// Given
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		surveyPKRows := sqlmock.NewRows([]string{"survey_pk"})
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_pk FROM survey.survey WHERE id = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(surveyPKRows)
		var postData = []byte(`{"name": "test", "classifierTypes": ["TEST1"]}`)

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiers"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)
		So(err, ShouldBeNil)

		// Then
		So(resp.StatusCode, ShouldEqual, http.StatusNotFound)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, `{"code":"404","message":"Survey not found for ID '67602ba2-8af6-4298-af66-4e46a62f32c8'",`)
	})
}

func TestCreateNewSurveyClassifiersAlreadyExistsConflict(t *testing.T) {
	Convey("Create new survey classifier already exists", t, func() {

		// Given
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		surveyPKRows := sqlmock.NewRows([]string{"survey_pk"}).AddRow("1000")
		classifierTypeSelectorMatchesRow := sqlmock.NewRows([]string{"Count"}).AddRow(1)
		prepareMockStmts(mock)
		mock.ExpectBegin()
		mock.ExpectPrepare("SELECT survey_pk FROM survey.survey WHERE id = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(surveyPKRows)
		mock.ExpectPrepare("SELECT COUNT\\(classifiertypeselector.id\\) FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.survey_fk = survey.survey_pk WHERE survey.id = .+ AND classifiertypeselector.classifier_type_selector = .+").ExpectQuery().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(classifierTypeSelectorMatchesRow)
		mock.ExpectRollback()
		var postData = []byte(`{"name": "test", "classifierTypes": ["TEST1"]}`)

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiers"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)
		So(err, ShouldBeNil)

		// Then
		So(resp.StatusCode, ShouldEqual, http.StatusInternalServerError)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Internal Server Error")
	})
}

func TestCreateNewSurveyClassifiersNoName(t *testing.T) {
	Convey("Create new survey classifier already exists", t, func() {

		// Given
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		surveyPKRows := sqlmock.NewRows([]string{"survey_pk"}).AddRow("1000")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_pk FROM survey.survey WHERE id = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(surveyPKRows)
		var postData = []byte(`{"classifierTypes": ["TEST"]}`)

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiers"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)
		So(err, ShouldBeNil)

		// Then
		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Invalid request body")
	})
}

func TestCreateNewSurveyClassifiersEmptyName(t *testing.T) {
	Convey("Create new survey classifier already exists", t, func() {

		// Given
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		surveyPKRows := sqlmock.NewRows([]string{"survey_pk"}).AddRow("1000")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_pk FROM survey.survey WHERE id = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(surveyPKRows)
		var postData = []byte(`{"name": "", "classifierTypes": ["TEST1"]}`)

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiers"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)
		So(err, ShouldBeNil)

		// Then
		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Invalid request body")
	})
}

func TestCreateNewSurveyClassifiersWhitespaceName(t *testing.T) {
	Convey("Create new survey classifier already exists", t, func() {

		// Given
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		surveyPKRows := sqlmock.NewRows([]string{"survey_pk"}).AddRow("1000")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_pk FROM survey.survey WHERE id = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(surveyPKRows)
		var postData = []byte(`{"name": " ", "classifierTypes": ["TEST1"]}`)

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiers"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)
		So(err, ShouldBeNil)

		// Then
		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Invalid request body")
	})
}

func TestCreateNewSurveyClassifiersNoClassifierTypes(t *testing.T) {
	Convey("Create new survey classifier already exists", t, func() {

		// Given
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		surveyPKRows := sqlmock.NewRows([]string{"survey_pk"}).AddRow("1000")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_pk FROM survey.survey WHERE id = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(surveyPKRows)
		var postData = []byte(`{"name": "test"}`)

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiers"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)
		So(err, ShouldBeNil)

		// Then
		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Invalid request body")
	})
}

func TestCreateNewSurveyClassifiersEmptyClassifierTypes(t *testing.T) {
	Convey("Create new survey classifier already exists", t, func() {

		// Given
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		surveyPKRows := sqlmock.NewRows([]string{"survey_pk"}).AddRow("1000")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_pk FROM survey.survey WHERE id = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(surveyPKRows)
		var postData = []byte(`{"name": "test",  "classifierTypes": []}`)

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiers"

		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))

		r, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)
		So(err, ShouldBeNil)

		// Then
		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Invalid request body")
	})
}

func TestCreateNewSurveyClassifiersWhitespaceClassifierTypes(t *testing.T) {
	Convey("Create new survey classifier already exists", t, func() {

		// Given
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		surveyPKRows := sqlmock.NewRows([]string{"survey_pk"}).AddRow("1000")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_pk FROM survey.survey WHERE id = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(surveyPKRows)
		var postData = []byte(`{"name": "test",  "classifierTypes": [" "]}`)

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiers"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)
		So(err, ShouldBeNil)

		// Then
		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Invalid request body")
	})
}

func TestCreateNewSurveyClassifiersEmptyStringClassifierTypes(t *testing.T) {
	Convey("Create new survey classifier already exists", t, func() {

		// Given
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		surveyPKRows := sqlmock.NewRows([]string{"survey_pk"}).AddRow("1000")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_pk FROM survey.survey WHERE id = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(surveyPKRows)
		var postData = []byte(`{"name": "test",  "classifierTypes": [""]}`)

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiers"
		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))
		r, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)
		So(err, ShouldBeNil)

		// Then
		So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Invalid request body")
	})
}

func TestCreateNewSurveyClassifiers500Error(t *testing.T) {
	Convey("Create new survey classifier returns 500", t, func() {

		// Given
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT survey_pk FROM survey.survey WHERE id = .+").ExpectQuery().WillReturnError(fmt.Errorf("Testing internal server error"))
		var postData = []byte(`{"name": "test", "classifierTypes": ["TEST1"]}`)

		// When
		api, err := models.NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()

		// Create a new router and plug in the defined routes
		router := mux.NewRouter()
		models.SetUpRoutes(router, api)

		ts := httptest.NewServer(router)
		defer ts.Close()
		url := ts.URL + "/surveys/" + surveyID + "/classifiers"

		// User and password not set so base64encode the dividing character
		basicAuth := base64.StdEncoding.EncodeToString([]byte(":"))

		r, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
		r.Header.Set("Authorization", "Basic: "+basicAuth)
		r.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(r)
		So(err, ShouldBeNil)

		// Then
		So(resp.StatusCode, ShouldEqual, http.StatusInternalServerError)
		body, err := io.ReadAll(resp.Body)
		So(string(body), ShouldStartWith, "Internal Server Error")
	})
}

func prepareMockStmts(m sqlmock.Sqlmock) {
	m.ExpectBegin()
	m.MatchExpectationsInOrder(false)
	m.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref ORDER BY short_name ASC")
	m.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref WHERE id = ?")
	m.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, s.eq_version, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref  WHERE LOWER\\(short_name\\) = LOWER\\(.+\\)")
	m.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref  WHERE LOWER\\(survey_ref\\) = LOWER\\(.+\\)")
	m.ExpectPrepare("SELECT id, s.short_name, s.long_name, s.survey_ref, s.legal_basis, s.survey_type, s.survey_mode, lb.long_name FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legal_basis = lb.ref")

	m.ExpectPrepare("SELECT ref, long_name FROM survey.legalbasis WHERE long_name = .+")
	m.ExpectPrepare("SELECT ref, long_name FROM survey.legalbasis WHERE ref = .+")

	m.ExpectPrepare("SELECT id, short_name, long_name, survey_ref, legal_basis, survey_type, survey_mode from survey.survey")
	m.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE LOWER\\(survey_ref\\) = LOWER\\(.*\\)")
	m.ExpectPrepare("UPDATE survey.survey SET short_name = .*, long_name = .* WHERE LOWER\\(survey_ref\\) = LOWER\\(.*\\)")
	m.ExpectPrepare("SELECT id FROM survey.survey WHERE id = .*")
	m.ExpectPrepare("DELETE FROM survey.survey WHERE id = .*")
	m.ExpectPrepare("SELECT classifiertypeselector.id, classifier_type_selector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.survey_fk = survey.survey_pk WHERE survey.id .*")
	m.ExpectPrepare("SELECT id, classifier_type_selector, classifier_type FROM survey.classifiertype INNER JOIN survey.classifiertypeselector ON classifiertype.classifier_type_selector_fk = classifiertypeselector.classifier_type_selector_pk .*")
	m.ExpectPrepare("INSERT INTO survey.survey \\( survey_pk, id, survey_ref, short_name, long_name, legal_basis, survey_type, survey_mode, eq_version \\) VALUES \\( .+\\)")
	m.ExpectPrepare("SELECT ref, long_name FROM survey.legalbasis")
	m.ExpectPrepare("SELECT survey_ref FROM survey.survey WHERE short_name = .+")
	m.ExpectPrepare("INSERT INTO survey.classifiertypeselector \\( classifier_type_selector_pk, id, survey_fk, classifier_type_selector \\) VALUES \\( .+\\) RETURNING classifier_type_selector_pk as id")
	m.ExpectPrepare("INSERT INTO survey.classifiertype \\( classifier_type_pk, classifier_type_selector_fk, classifier_type \\) VALUES \\( .+\\)")
	m.ExpectPrepare("SELECT survey_pk FROM survey.survey WHERE id = .+")
	m.ExpectPrepare("SELECT COUNT\\(classifiertypeselector.id\\) FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.survey_fk = survey.survey_pk WHERE survey.id = .+ AND classifiertypeselector.classifier_type_selector = .+")
}
