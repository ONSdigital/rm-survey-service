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
	dataSource, _ := configureEnvironment()
	models.InitDB(dataSource)
	os.Exit(m.Run())
}

func TestEndpoints(t *testing.T) {

	os.Setenv("DATABASE_URL", "postgres://postgres:password@database/postgres?sslmode=disable")
	os.Setenv("security_user_name", "admin")
	os.Setenv("security_user_password", "secret")

	auth_str := "admin:secret"
	encoded_auth_str := base64.StdEncoding.EncodeToString([]byte(auth_str))

	r := gofight.New()

	r.GET("/info").
		SetDebug(true).
		Run(EchoEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			name, _ := jsonparser.GetString(data, "name")
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "surveysvc", name)
		})

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
