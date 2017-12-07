package main

import (
	"encoding/base64"
	"net/http"
	"os"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/buger/jsonparser"
	"github.com/onsdigital/rm-survey-service/models"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	// Configure db
	dataSource, _ := configureEnvironment()
	models.InitDB(dataSource)

	// Set auth params
	os.Setenv("security_user_name", "admin")
	os.Setenv("security_user_password", "secret")

	// Run the tests and exit
	os.Exit(m.Run())
}

func TestInfoEndpoint(t *testing.T) {
	r := gofight.New()

	r.GET("/info").
		SetDebug(true).
		Run(EchoEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			name, _ := jsonparser.GetString(data, "name")
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "surveysvc", name)
		})
}

func TestGetSurveysEndpoint(t *testing.T) {
	r := gofight.New()

	auth_str := "admin:secret"
	encoded_auth_str := base64.StdEncoding.EncodeToString([]byte(auth_str))

	r.GET("/surveys").
		SetDebug(true).
		SetHeader(gofight.H{
			"Authorization": "Basic " + string(encoded_auth_str),
		}).
		Run(EchoEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			expected := `[{"id":"cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87","shortName":"BRES"}]`
			assert.Equal(t, expected, r.Body.String())
		})
}

func TestGetSurveyEndpoint(t *testing.T) {
	r := gofight.New()

	auth_str := "admin:secret"
	encoded_auth_str := base64.StdEncoding.EncodeToString([]byte(auth_str))
	id := "cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87"
	r.GET("/surveys/"+id).
		SetDebug(true).
		SetHeader(gofight.H{
			"Authorization": "Basic " + string(encoded_auth_str),
		}).
		Run(EchoEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			expected := `{"id":"cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87",` +
				`"shortName":"BRES",` +
				`"longName":"Business Register and Employment Survey",` +
				`"surveyRef":"221"}`
			assert.Equal(t, expected, r.Body.String())
		})
}

func TestGetSurveyByShortNameEndpoint(t *testing.T) {
	r := gofight.New()

	auth_str := "admin:secret"
	encoded_auth_str := base64.StdEncoding.EncodeToString([]byte(auth_str))
	r.GET("/surveys/shortname/BRES").
		SetDebug(true).
		SetHeader(gofight.H{
			"Authorization": "Basic " + string(encoded_auth_str),
		}).
		Run(EchoEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			expected := `{"id":"cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87",` +
				`"shortName":"BRES",` +
				`"longName":"Business Register and Employment Survey",` +
				`"surveyRef":"221"}`
			assert.Equal(t, expected, r.Body.String())
		})
}

func TestGetSurveyByReferenceEndpoint(t *testing.T) {
	r := gofight.New()

	auth_str := "admin:secret"
	encoded_auth_str := base64.StdEncoding.EncodeToString([]byte(auth_str))
	r.GET("/surveys/ref/221").
		SetDebug(true).
		SetHeader(gofight.H{
			"Authorization": "Basic " + string(encoded_auth_str),
		}).
		Run(EchoEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			expected := `{"id":"cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87",` +
				`"shortName":"BRES",` +
				`"longName":"Business Register and Employment Survey",` +
				`"surveyRef":"221"}`
			assert.Equal(t, expected, r.Body.String())
		})
}

func TestAllClassifierTypeSelectorsEndpoint(t *testing.T) {
	r := gofight.New()

	auth_str := "admin:secret"
	encoded_auth_str := base64.StdEncoding.EncodeToString([]byte(auth_str))
	r.GET("/surveys/cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87/classifiertypeselectors").
		SetDebug(true).
		SetHeader(gofight.H{
			"Authorization": "Basic " + string(encoded_auth_str),
		}).
		Run(EchoEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			expected := `[{"id":"efa868fb-fb80-44c7-9f33-d6800a17c4da","name":"COLLECTION_INSTRUMENT"},{"id":"e119ffd6-6fc1-426c-ae81-67a96f9a71ba","name":"COMMUNICATION_TEMPLATE"}]`
			assert.Equal(t, expected, r.Body.String())
		})
}
