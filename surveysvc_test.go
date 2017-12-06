package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {

	os.Setenv("DATABASE_URL", "postgres://postgres:password@database/postgres?sslmode=disable")
	os.Setenv("security_user_name", "admin")
	os.Setenv("security_user_password", "secret")

	r := gofight.New()

	r.GET("/info").
		SetDebug(true).
		Run(EchoEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)

			data := []byte(r.Body.String())
			name, _ := jsonparser.GetString(data, "name")

			assert.Equal(t, "surveysvc", name)
		})
}
