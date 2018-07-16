// +build integration

package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func createBasicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// Assumes the service is running on localhost, port 9090
func TestAPI_Info(t *testing.T) {
	Convey("Can post new survey classifiers", t, func() {

		// Given
		request, err := http.NewRequest("GET", "http://localhost:9090/info", nil)
		So(err, ShouldBeNil)
		client := http.Client{}

		// When
		response, err := client.Do(request)
		So(err, ShouldBeNil)

		// Then
		So(response.StatusCode, ShouldEqual, http.StatusOK)
	})
}

// Assumes the service is running on localhost, port 9090
func TestAPI_PostSurveyClassifiers(t *testing.T) {
	Convey("Can post new survey classifiers", t, func() {

		// Create a classifier made of a classifier type selector containing it's classifier types
		classifier := ClassifierTypeSelector{
			Name:            "TEST_SELECTOR_TYPE1",
			ClassifierTypes: []string{"TEST1", "TEST2"},
		}

		// Create HTTP request to post the classifier as JSON
		postData, err := json.Marshal(classifier)
		So(err, ShouldBeNil)
		request, err := http.NewRequest("POST", "http://localhost:9090/surveys/cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87/classifiers", bytes.NewReader(postData))
		So(err, ShouldBeNil)
		apiAuth := createBasicAuth("admin", "secret")
		request.Header.Add("Authorization", apiAuth)
		request.Header.Add("Content-Type", "application/json")

		// Post the classifier and assert success 201 is the response status
		client := http.Client{}
		response, err := client.Do(request)
		So(err, ShouldBeNil)
		So(response.StatusCode, ShouldEqual, http.StatusCreated)

		// Decode and parse the response body to get the ID
		var setupResponseClassifier ClassifierTypeSelector
		defer response.Body.Close()
		json.NewDecoder(response.Body).Decode(&setupResponseClassifier)

		// Use the ID to get the classifier we posted by a GET request
		getClassifier, err := http.NewRequest("GET", "http://localhost:9090/surveys/cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87/classifiertypeselectors/"+setupResponseClassifier.ID, nil)
		So(err, ShouldBeNil)
		getClassifier.Header.Add("Authorization", apiAuth)
		getResponseClassifiers, err := client.Do(getClassifier)
		So(err, ShouldBeNil)
		So(getResponseClassifiers.StatusCode, ShouldEqual, http.StatusOK)

		// Assert that the classifier returned by the GET request is the same as we originally posted
		var retrievedClassifier ClassifierTypeSelector
		defer getResponseClassifiers.Body.Close()
		err = json.NewDecoder(getResponseClassifiers.Body).Decode(&retrievedClassifier)
		So(err, ShouldBeNil)
		So(retrievedClassifier.Name, ShouldEqual, classifier.Name)
		So(retrievedClassifier.ClassifierTypes, ShouldContain, classifier.ClassifierTypes[0])
		So(retrievedClassifier.ClassifierTypes, ShouldContain, classifier.ClassifierTypes[1])
		So(retrievedClassifier.ClassifierTypes, ShouldHaveLength, 2)
	})
}
