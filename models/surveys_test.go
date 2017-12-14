package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInfoEndpoint(t *testing.T) {

	Convey("Info enpoint returns a 200 response", t, func() {
		db, mock, err := sqlmock.New()
		prepareMockStmts(mock)
		So(err, ShouldBeNil)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
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
		mock.ExpectPrepare("SELECT id, shortname FROM survey.survey").ExpectQuery().WillReturnRows(makeCollectionRow())
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/surveys", nil)
		So(err, ShouldBeNil)
		api.AllSurveys(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)
		expected := []SurveySummary{{ID: "testid", ShortName: "test-shortname"}}
		res := []SurveySummary{}
		json.Unmarshal(w.Body.Bytes(), &res)
		So(res[0].ID, ShouldEqual, expected[0].ID)
		So(res[0].ShortName, ShouldEqual, expected[0].ShortName)
	})
}

func TestSurveyListInternalServerError(t *testing.T) {
	Convey("Surveys list returns a 500", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT id, shortname FROM survey.survey").ExpectQuery().WillReturnError(fmt.Errorf("Testing internal server error"))
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/surveys", nil)
		So(err, ShouldBeNil)
		api.AllSurveys(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})
}

func TestSurveyListNotFound(t *testing.T) {
	Convey("Surveys list returns an 500 not found", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "shortname"})
		mock.ExpectPrepare("SELECT id, shortname FROM survey.survey").ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/surveys", nil)
		So(err, ShouldBeNil)
		api.AllSurveys(w, r)
		So(w.Code, ShouldEqual, http.StatusNoContent)
	})
}

func TestSurveyGetReturnsJson(t *testing.T) {
	Convey("Survey GET returns a survey resource", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "shortname", "longname", "surveyref"}).AddRow("testid", "test-shortname", "test-longname", "test-reference")
		mock.ExpectPrepare("SELECT id, shortname, longname, surveyref from survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/surveys/testid", nil)
		So(err, ShouldBeNil)
		api.GetSurvey(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)
		expected := Survey{ID: "testid", ShortName: "test-shortname", LongName: "test-longname", Reference: "test-reference"}
		res := Survey{}
		json.Unmarshal(w.Body.Bytes(), &res)
		So(res.ID, ShouldEqual, expected.ID)
		So(res.ShortName, ShouldEqual, expected.ShortName)
		So(res.LongName, ShouldEqual, expected.LongName)
		So(res.Reference, ShouldEqual, expected.Reference)
	})
}

func TestSurveyGetNotFound(t *testing.T) {
	Convey("Survey Get returns an 404 not found", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "shortname", "longname", "reference"})
		mock.ExpectPrepare("SELECT id, shortname, longname, surveyref from survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/survey/testid", nil)
		So(err, ShouldBeNil)
		api.GetSurvey(w, r)
		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

func TestSurveyGetInternalServerError(t *testing.T) {
	Convey("Survey GET returns a 500", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT id, shortname, longname, surveyref from survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnError(fmt.Errorf("Testing internal server error"))
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/survey/testid", nil)
		So(err, ShouldBeNil)
		api.GetSurvey(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})
}

func makeCollectionRow() *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "shortname"}).
		AddRow("testid", "test-shortname")
	return rows
}

func prepareMockStmts(m sqlmock.Sqlmock) {
	m.ExpectBegin()
	m.MatchExpectationsInOrder(false)
	m.ExpectPrepare("SELECT id, shortname FROM survey.survey ORDER BY shortname ASC")
	m.ExpectPrepare("SELECT id, shortname, longname, surveyref from survey.survey WHERE id = ?")
	m.ExpectPrepare("SELECT id, shortname, longname, surveyref from survey.survey WHERE LOWER(shortName) = LOWER(?)")
}
