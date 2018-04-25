package models

import (
	"bytes"
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
		rows := sqlmock.NewRows([]string{"id", "shortname", "longname", "surveyref", "legalbasis","longname"}).AddRow("testid", "test-shortname", "test-longname", "test-reference", "test-legalbasis-ref","test-legalbasis-longname")
		mock.ExpectPrepare("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref").ExpectQuery().WillReturnRows(rows)
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
		expected := []Survey{{ID: "testid", ShortName: "test-shortname"}}
		res := []Survey{}
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
		mock.ExpectPrepare("SELECT id, shortname, longname, surveyref, legalbasis FROM survey.survey").ExpectQuery().WillReturnError(fmt.Errorf("Testing internal server error"))
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
		rows := sqlmock.NewRows([]string{"id", "shortname", "longname", "surveyref", "legalbasis","longname"})
		mock.ExpectPrepare("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref").ExpectQuery().WillReturnRows(rows)
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
		rows := sqlmock.NewRows([]string{"id", "shortname", "longname", "surveyref", "legalbasis","longname"}).AddRow("testid", "test-shortname", "test-longname", "test-reference", "test-legalbasis-ref","test-legalbasis-longname")
		mock.ExpectPrepare("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
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
		rows := sqlmock.NewRows([]string{"id", "shortname", "longname", "surveyref", "legalbasis","longname"})
		mock.ExpectPrepare("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
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
		mock.ExpectPrepare("SELECT id, shortname, longname, surveyref, legalbasis from survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnError(fmt.Errorf("Testing internal server error"))
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

func TestGetSurveyByShortnameReturnsJSON(t *testing.T) {
	Convey("Survey GET by shortname returns a survey resource", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "shortname", "longname", "surveyref", "legalbasis","longname"}).AddRow("testid", "test-shortname", "test-longname", "test-reference", "test-legalbasis-ref","test-legalbasis-longname")
		mock.ExpectPrepare("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref").ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/surveys/shortname/test-shortname", nil)
		So(err, ShouldBeNil)
		api.GetSurveyByShortName(w, r)
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

func TestSurveyGetByShortNameNotFound(t *testing.T) {
	Convey("Survey Get by shortname returns an 404 not found", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "shortname", "longname", "surveyref", "legalbasis","longname"})
		mock.ExpectPrepare("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref").ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/survey/shortname/test-shortname", nil)
		So(err, ShouldBeNil)
		api.GetSurveyByShortName(w, r)
		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

func TestSurveyGetByShortNameInternalServerError(t *testing.T) {
	Convey("Survey GET by shortname returns a 500", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT id, shortname, longname, surveyref, legalbasis from survey.survey").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnError(fmt.Errorf("Testing internal server error"))
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/survey/testid", nil)
		So(err, ShouldBeNil)
		api.GetSurveyByShortName(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})
}

func TestGetSurveyByReferenceReturnsJSON(t *testing.T) {
	Convey("Survey GET by reference returns a survey resource", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "shortname", "longname", "surveyref", "legalbasis","longname"}).AddRow("testid", "test-shortname", "test-longname", "test-reference", "test-legalbasis-ref","test-legalbasis-longname")
		mock.ExpectPrepare("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref").ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/surveys/ref/test-reference", nil)
		So(err, ShouldBeNil)
		api.GetSurveyByReference(w, r)
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

func TestSurveyGetByReferenceNotFound(t *testing.T) {
	Convey("Survey Get by reference returns an 404 not found", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "shortname", "longname", "surveyref", "legalbasis","longname"})
		mock.ExpectPrepare("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref").ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/survey/ref/test-reference", nil)
		So(err, ShouldBeNil)
		api.GetSurveyByReference(w, r)
		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

func TestSurveyGetByReferenceInternalServerError(t *testing.T) {
	Convey("Survey GET by reference returns a 500", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT id, shortname, longname, surveyref, legalbasis from survey.survey").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnError(fmt.Errorf("Testing internal server error"))
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
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
		rows := sqlmock.NewRows([]string{"id", "classifiertypeselector"}).AddRow("test-id", "test-name")
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(idRow)
		mock.ExpectPrepare("SELECT classifiertypeselector.id, classifiertypeselector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.surveyfk = survey.surveypk WHERE survey.id = .* ORDER BY classifiertypeselector ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/surveys/test-id/classifiertypeselectors/", nil)
		So(err, ShouldBeNil)
		api.AllClassifierTypeSelectors(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)
		expected := ClassifierTypeSelectorSummary{ID: "test-id", Name: "test-name"}
		res := []ClassifierTypeSelectorSummary{}
		json.Unmarshal(w.Body.Bytes(), &res)
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
		rows := sqlmock.NewRows([]string{"id", "classifiertypeselector"}).AddRow("test-id", "test-name")
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(idRow)
		mock.ExpectPrepare("SELECT classifiertypeselector.id, classifiertypeselector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.surveyfk = survey.surveypk WHERE survey.id = .* ORDER BY classifiertypeselector ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/surveys/test-id/classifiertypeselectors/", nil)
		So(err, ShouldBeNil)
		api.AllClassifierTypeSelectors(w, r)
		So(w.Code, ShouldEqual, http.StatusNotFound)
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
		mock.ExpectPrepare("SELECT classifiertypeselector.id, classifiertypeselector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.surveyfk = survey.surveypk WHERE survey.id = .* ORDER BY classifiertypeselector ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/surveys/test-id/classifiertypeselectors/", nil)
		So(err, ShouldBeNil)
		api.AllClassifierTypeSelectors(w, r)
		So(w.Code, ShouldEqual, http.StatusNoContent)
	})
}

func TestAllClassifierTypeSelectorsSurveyReturnsInternalServerError(t *testing.T) {
	Convey("ClassifierType GET returns a 500", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "classifiertypeselector"}).AddRow("test-id", "test-name")
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnError(fmt.Errorf("Testing internal server error"))
		mock.ExpectPrepare("SELECT classifiertypeselector.id, classifiertypeselector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.surveyfk = survey.surveypk WHERE survey.id = .* ORDER BY classifiertypeselector ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/surveys/test-id/classifiertypeselectors/", nil)
		So(err, ShouldBeNil)
		api.AllClassifierTypeSelectors(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})
}

func TestAllClassifierTypeSelectorsReturnsInternalServerError(t *testing.T) {
	Convey("ClassifierType GET returns a 500", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT classifiertypeselector.id, classifiertypeselector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.surveyfk = survey.surveypk WHERE survey.id = .* ").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnError(fmt.Errorf("Testing internal server error"))
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/surveys/test-id/classifiertypeselectors", nil)
		So(err, ShouldBeNil)
		api.AllClassifierTypeSelectors(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})
}

func TestClassifierTypeSelectorByIdReturnsJSON(t *testing.T) {
	Convey("ClassifierType GET by ID returns a classifier resource", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		idRow := sqlmock.NewRows([]string{"id"}).AddRow("id").AddRow("id")
		rows := sqlmock.NewRows([]string{"id", "classifiertypeselector", "classifiertype"}).AddRow("test-id", "test-name", "test-type")
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(idRow)
		mock.ExpectPrepare("SELECT id, classifiertypeselector, classifiertype FROM survey.classifiertype INNER JOIN survey.classifiertypeselector ON classifiertype.classifiertypeselectorfk = classifiertypeselector.classifiertypeselectorpk WHERE classifiertypeselector.id = .* ORDER BY classifiertype ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/surveys/test-id/classifiertypeselectors/", nil)
		So(err, ShouldBeNil)
		api.GetClassifierTypeSelectorByID(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)
		var a = []string{"test"}
		expected := ClassifierTypeSelector{ID: "test-id", Name: "test-name", ClassifierTypes: a}
		res := ClassifierTypeSelector{}
		json.Unmarshal(w.Body.Bytes(), &res)
		So(res.ID, ShouldEqual, expected.ID)
		So(res.Name, ShouldEqual, expected.Name)
	})
}

func TestClassifierTypeSelectorByIdReturns404(t *testing.T) {
	Convey("ClassifierType GET by ID returns a classifier resource", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		idRow := sqlmock.NewRows([]string{"id"})
		rows := sqlmock.NewRows([]string{"id", "classifiertypeselector", "classifiertype"}).AddRow("test-id", "test-name", "test-type")
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(idRow)
		mock.ExpectPrepare("SELECT id, classifiertypeselector, classifiertype FROM survey.classifiertype INNER JOIN survey.classifiertypeselector ON classifiertype.classifiertypeselectorfk = classifiertypeselector.classifiertypeselectorpk WHERE classifiertypeselector.id = .* ORDER BY classifiertype ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/surveys/test-id/classifiertypeselectors/", nil)
		So(err, ShouldBeNil)
		api.GetClassifierTypeSelectorByID(w, r)
		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

func TestClassifierTypeSelectorByIdNoClassifierTypesReturns404(t *testing.T) {
	Convey("ClassifierType GET by ID returns 404 if no classifier types exist", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		idRow := sqlmock.NewRows([]string{"id"}).AddRow("test-id")
		rows := sqlmock.NewRows([]string{"id", "classifiertypeselector", "classifiertype"})
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(idRow)
		mock.ExpectPrepare("SELECT id, classifiertypeselector, classifiertype FROM survey.classifiertype INNER JOIN survey.classifiertypeselector ON classifiertype.classifiertypeselectorfk = classifiertypeselector.classifiertypeselectorpk WHERE classifiertypeselector.id = .* ORDER BY classifiertype ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/surveys/test-id/classifiertypeselectors/test-selector", nil)
		So(err, ShouldBeNil)
		api.GetClassifierTypeSelectorByID(w, r)
		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

func TestClassifierTypeSelectorByIdInternalServerError(t *testing.T) {
	Convey("ClassifierType GET by reference returns a classifier resource", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		rows := sqlmock.NewRows([]string{"id", "classifiertypeselector", "classifiertype"}).AddRow("test-id", "test-name", "test-type")
		mock.ExpectPrepare("SELECT id FROM survey.survey WHERE id = ?").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnError(fmt.Errorf("Testing internal server error"))
		mock.ExpectPrepare("SELECT id, classifiertypeselector, classifiertype FROM survey.classifiertype INNER JOIN survey.classifiertypeselector ON classifiertype.classifiertypeselectorfk = classifiertypeselector.classifiertypeselectorpk WHERE classifiertypeselector.id = .* ORDER BY classifiertype ASC").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/surveys/test-id/classifiertypeselectors/", nil)
		So(err, ShouldBeNil)
		api.GetClassifierTypeSelectorByID(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})
}

func TestPutSurveyDetailsBySurveyRefSuccess(t *testing.T) {
	Convey("Survey Details PUT by Survey Reference success", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		refRow := sqlmock.NewRows([]string{"surveyref"}).AddRow("456")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT surveyref FROM survey.survey WHERE LOWER\\(surveyref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(refRow)
		mock.ExpectPrepare("UPDATE survey.survey SET shortname = .+, longname = .+ WHERE LOWER\\(surveyref\\) = LOWER\\(.+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name"}`)
		r, err := http.NewRequest("PUT", "http://localhost:9090/surveys/ref/456", bytes.NewBuffer(jsonStr))
		So(err, ShouldBeNil)
		api.PutSurveyDetails(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

func TestPutSurveyDetailsBySurveyRefInternalServerError(t *testing.T) {
	Convey("Survey Details PUT by Survey Reference success", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT surveyref FROM survey.survey WHERE LOWER\\(surveyref\\) = LOWER\\(.+\\)").ExpectQuery().WillReturnError(fmt.Errorf("Testing internal server error"))
		mock.ExpectPrepare("UPDATE survey.survey SET shortname = .+, longname = .+ WHERE LOWER\\(surveyref\\) = LOWER\\(.+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name"}`)
		r, err := http.NewRequest("PUT", "http://localhost:9090/surveys/ref/456", bytes.NewBuffer(jsonStr))
		So(err, ShouldBeNil)
		api.PutSurveyDetails(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})
}

func TestCreateNewSurvey(t *testing.T) {
	Convey("Create new survey success", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref","longname"}).AddRow("STA1947", "Statistics of Trade Act 1947")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT surveyref FROM survey.survey WHERE LOWER\\(surveyref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("INSERT INTO survey.survey \\( surveypk, id, surveyref, shortname, longname, legalbasis \\) VALUES \\( .+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectPrepare("SELECT ref, longname FROM survey.legalbasis WHERE longname = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"99","LegalBasis":"Statistics of Trade Act 1947"}`)
		r, err := http.NewRequest("POST", "http://localhost:9090/surveys", bytes.NewBuffer(jsonStr))
		So(err, ShouldBeNil)
		api.PostSurveyDetails(w, r)
		So(w.Code, ShouldEqual, http.StatusCreated)
	})
}

func TestCreateNewSurveyNonExistentLegalBasisRef(t *testing.T) {
	Convey("Create new survey with non existent legal basis ref", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref","longname"})
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT surveyref FROM survey.survey WHERE LOWER\\(surveyref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("INSERT INTO survey.survey \\( surveypk, id, surveyref, shortname, longname, legalbasis \\) VALUES \\( .+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectPrepare("SELECT ref, longname FROM survey.legalbasis WHERE ref = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"99","LegalBasisRef":"Statistics of Trade Act 1947"}`)
		r, err := http.NewRequest("POST", "http://localhost:9090/surveys", bytes.NewBuffer(jsonStr))
		So(err, ShouldBeNil)
		api.PostSurveyDetails(w, r)
		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func TestCreateNewSurveyNonExistentLegalBasisLongName(t *testing.T) {
	Convey("Create new survey with non existent legal basis ref", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref","longname"})
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT surveyref FROM survey.survey WHERE LOWER\\(surveyref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("INSERT INTO survey.survey \\( surveypk, id, surveyref, shortname, longname, legalbasis \\) VALUES \\( .+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectPrepare("SELECT ref, longname FROM survey.legalbasis WHERE longname = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"99","LegalBasis":"Statistics of Trade Act 1947"}`)
		r, err := http.NewRequest("POST", "http://localhost:9090/surveys", bytes.NewBuffer(jsonStr))
		So(err, ShouldBeNil)
		api.PostSurveyDetails(w, r)
		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func TestCreateNewSurveyNonExistentLegalBasis(t *testing.T) {
	Convey("Create new survey with non existent legal basis ref", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref","longname"})
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT surveyref FROM survey.survey WHERE LOWER\\(surveyref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("INSERT INTO survey.survey \\( surveypk, id, surveyref, shortname, longname, legalbasis \\) VALUES \\( .+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectPrepare("SELECT ref, longname FROM survey.legalbasis WHERE longname = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"99"}`)
		r, err := http.NewRequest("POST", "http://localhost:9090/surveys", bytes.NewBuffer(jsonStr))
		So(err, ShouldBeNil)
		api.PostSurveyDetails(w, r)
		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func TestCreateNewSurveyNonNumericRef(t *testing.T) {
	Convey("Create new survey with non numeric refernce", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref","longname"}).AddRow("STA1947", "Statistics of Trade Act 1947")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT surveyref FROM survey.survey WHERE LOWER\\(surveyref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("INSERT INTO survey.survey \\( surveypk, id, surveyref, shortname, longname, legalbasis \\) VALUES \\( .+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectPrepare("SELECT ref, longname FROM survey.legalbasis WHERE longname = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"99A"}`)
		r, err := http.NewRequest("POST", "http://localhost:9090/surveys", bytes.NewBuffer(jsonStr))
		So(err, ShouldBeNil)
		api.PostSurveyDetails(w, r)
		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func TestCreateNewSurveyRefTooLong(t *testing.T) {
	Convey("Create new survey with non numeric refernce", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref","longname"}).AddRow("STA1947", "Statistics of Trade Act 1947")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT surveyref FROM survey.survey WHERE LOWER\\(surveyref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("INSERT INTO survey.survey \\( surveypk, id, surveyref, shortname, longname, legalbasis \\) VALUES \\( .+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectPrepare("SELECT ref, longname FROM survey.legalbasis WHERE longname = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name","SurveyRef":"012345678901234567890"}`)
		r, err := http.NewRequest("POST", "http://localhost:9090/surveys", bytes.NewBuffer(jsonStr))
		So(err, ShouldBeNil)
		api.PostSurveyDetails(w, r)
		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func TestCreateNewSurveyShortNameWithSpace(t *testing.T) {
	Convey("Create new survey with non numeric refernce", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref","longname"}).AddRow("STA1947", "Statistics of Trade Act 1947")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT surveyref FROM survey.survey WHERE LOWER\\(surveyref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("INSERT INTO survey.survey \\( surveypk, id, surveyref, shortname, longname, legalbasis \\) VALUES \\( .+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectPrepare("SELECT ref, longname FROM survey.legalbasis WHERE longname = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		var jsonStr = []byte(`{"ShortName": "test short name", "LongName":"test-long-name","SurveyRef":"0123"}`)
		r, err := http.NewRequest("POST", "http://localhost:9090/surveys", bytes.NewBuffer(jsonStr))
		So(err, ShouldBeNil)
		api.PostSurveyDetails(w, r)
		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func TestCreateNewSurveyShortNameTooLong(t *testing.T) {
	Convey("Create new survey with non numeric refernce", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref","longname"}).AddRow("STA1947", "Statistics of Trade Act 1947")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT surveyref FROM survey.survey WHERE LOWER\\(surveyref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("INSERT INTO survey.survey \\( surveypk, id, surveyref, shortname, longname, legalbasis \\) VALUES \\( .+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectPrepare("SELECT ref, longname FROM survey.legalbasis WHERE longname = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		var jsonStr = []byte(`{"ShortName": "test-short-name-0123456", "LongName":"test-long-name","SurveyRef":"0123"}`)
		r, err := http.NewRequest("POST", "http://localhost:9090/surveys", bytes.NewBuffer(jsonStr))
		So(err, ShouldBeNil)
		api.PostSurveyDetails(w, r)
		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func TestCreateNewSurveyLongNameTooLong(t *testing.T) {
	Convey("Create new survey with non numeric refernce", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		rows := sqlmock.NewRows([]string{"surveyref"})
		legalBasis := sqlmock.NewRows([]string{"ref","longname"}).AddRow("STA1947", "Statistics of Trade Act 1947")
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT surveyref FROM survey.survey WHERE LOWER\\(surveyref\\) = LOWER\\(.+\\)").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectPrepare("INSERT INTO survey.survey \\( surveypk, id, surveyref, shortname, longname, legalbasis \\) VALUES \\( .+\\)").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectPrepare("SELECT ref, longname FROM survey.legalbasis WHERE longname = .+").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(legalBasis)
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		var jsonStr = []byte(`{"ShortName": "test-short-name", "LongName":"test-long-name-012345678-012345678-012345678-012345678-012345678-012345678-012345678-01234567899999999-0123456789","SurveyRef":"0123"}`)
		r, err := http.NewRequest("POST", "http://localhost:9090/surveys", bytes.NewBuffer(jsonStr))
		So(err, ShouldBeNil)
		api.PostSurveyDetails(w, r)
		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func prepareMockStmts(m sqlmock.Sqlmock) {
	m.ExpectBegin()
	m.MatchExpectationsInOrder(false)
	m.ExpectPrepare("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref ORDER BY shortname ASC")
	m.ExpectPrepare("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref WHERE id = ?")
	m.ExpectPrepare("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref  WHERE LOWER\\(shortName\\) = LOWER\\(.+\\)")
	m.ExpectPrepare("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref  WHERE LOWER\\(surveyref\\) = LOWER\\(.+\\)")
	m.ExpectPrepare("SELECT id, s.shortname, s.longname, s.surveyref, s.legalbasis, lb.longname FROM survey.survey s INNER JOIN survey.legalbasis lb on s.legalbasis = lb.ref")

	m.ExpectPrepare("SELECT ref, longname FROM survey.legalbasis WHERE longname = .+")
	m.ExpectPrepare("SELECT ref, longname FROM survey.legalbasis WHERE ref = .+")

	m.ExpectPrepare("SELECT id, shortname, longname, surveyref, legalbasis from survey.survey")
	m.ExpectPrepare("SELECT surveyref FROM survey.survey WHERE LOWER\\(surveyref\\) = LOWER\\(.*\\)")
	m.ExpectPrepare("UPDATE survey.survey SET shortname = .*, longname = .* WHERE LOWER\\(surveyref\\) = LOWER\\(.*\\)")
	m.ExpectPrepare("SELECT id FROM survey.survey WHERE id = .*")
	m.ExpectPrepare("SELECT classifiertypeselector.id, classifiertypeselector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.surveyfk = survey.surveypk WHERE survey.id .*")
	m.ExpectPrepare("SELECT id, classifiertypeselector, classifiertype FROM survey.classifiertype INNER JOIN survey.classifiertypeselector ON classifiertype.classifiertypeselectorfk = classifiertypeselector.classifiertypeselectorpk .*")
	m.ExpectPrepare("INSERT INTO survey.survey \\( surveypk, id, surveyref, shortname, longname, legalbasis \\) VALUES \\( .+\\)")
	m.ExpectPrepare("SELECT ref, longname FROM survey.legalbasis")

}
